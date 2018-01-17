package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/daticahealth/cli/config"
	"github.com/daticahealth/cli/models"
	u2f "github.com/marshallbrekka/go-u2fhost"
)

// Signin signs in a user and returns the representative user model. If an
// error occurs, nil is returned for the user and the error field is populated.
func (a *SAuth) Signin() (*models.User, error) {
	// if we're already signed in with a valid session, don't sign in again
	if user, err := a.Verify(); err == nil {
		return user, nil
	}
	f := a.signInWithKey
	if a.Settings.PrivateKeyPath == "" {
		f = a.signInWithCredentials
	}
	signinResp, err := f("")
	if err != nil {
		return nil, err
	}

	var user *models.User

	if signinResp.MFAID != "" {
		user, tryDifferent, err = a.mfaSignin(signinResp.MFAID, signinResp.MFAPreferredMode, signinResp.Challenge)
		if err != nil {
			if tryDifferent {
				err = a.Prompts.YesNo("%s\nWould you like to try a different mfa method (y/n)?")
				if err != nil {
					return nil, err
				}
				fmt.Println("1.) Authenticator App")
				fmt.Println("2.) Email")
				fmt.Println("3.) U2F Key")
				i, err := strconv.Atoi(a.Prompts.CaptureInput("Which method?"))
				if err != nil {
					return nil, err
				}
				var string mfaType
				switch i {
				case 1:
					mfaType = "authenticator"
				case 2:
					mfaType = "email"
				case 3:
					mfaType = "u2f"
				}
				signinResp, err = f(mfaType)
				user, _, err = a.mfaSignin(signinResp.MFAID, signinResp.MFAPreferredMode, signinResp.Challenge)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	} else {
		user = signinResp.toUser()
	}

	a.Settings.UsersID = user.UsersID
	a.Settings.Email = user.Email
	a.Settings.SessionToken = user.SessionToken

	config.SaveSettings(a.Settings)

	return user, nil
}

type signinResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	SessionToken     string `json:"sessionToken"`
	MFAID            string `json:"mfaID"`
	MFAPreferredMode string `json:"mfaPreferredType"`
	MFAChallenge     string `json:"mfaChallenge,omitempty"`
}

type u2fSignRequest struct {
	AppID          string             `json:"appId"`
	Challenge      string             `json:"challenge"`
	RegisteredKeys []u2fRegisteredKey `json:"registeredKeys"`
}

type u2fRegisteredKey struct {
	Version   string `json:"version"`
	KeyHandle string `json:"keyHandle"`
	AppID     string `json:"appId"`
}

func (sr *signinResponse) toUser() *models.User {
	return &models.User{
		UsersID:      sr.ID,
		Email:        sr.Email,
		SessionToken: sr.SessionToken,
	}
}

func (a *SAuth) signInWithCredentials(mfaType string) (*signinResponse, error) {
	login := models.Login{
		Identifier: a.Settings.Email,
		Password:   a.Settings.Password,
	}
	if a.Settings.Email == "" || a.Settings.Password == "" {
		email, password, err := a.Prompts.EmailPassword(a.Settings.Email, a.Settings.Password)
		if err != nil {
			return nil, err
		}
		login = models.Login{
			Identifier: email,
			Password:   password,
		}
	}

	b, err := json.Marshal(login)
	if err != nil {
		return nil, err
	}
	headers := a.Settings.HTTPManager.GetHeaders(a.Settings.SessionToken, a.Settings.Version, a.Settings.Pod, a.Settings.UsersID)
	var mfaQuery string
	if len(mfaType) > 0 {
		mfaQuery = fmt.Sprintf("?mfaType=%s", mfaType)
	}
	resp, statusCode, err := a.Settings.HTTPManager.Post(b, fmt.Sprintf("%s%s/auth/signin%s", a.Settings.AuthHost, a.Settings.AuthHostVersion, mfaQuery), headers)
	if err != nil {
		return nil, err
	}
	signinResp := &signinResponse{}
	return signinResp, a.Settings.HTTPManager.ConvertResp(resp, statusCode, signinResp)
}

func (a *SAuth) signInWithKey(mfaType string) (*signinResponse, error) {
	body := struct {
		PublicKey string `json:"publicKey"`
		Signature string `json:"signature"`
	}{}

	bytes, err := ioutil.ReadFile(a.Settings.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("Private key is not PEM-encoded")
	}
	bytes = block.Bytes
	if x509.IsEncryptedPEMBlock(block) {
		passphrase := a.Prompts.KeyPassphrase(a.Settings.PrivateKeyPath)
		bytes, err = x509.DecryptPEMBlock(block, []byte(passphrase))
		if err != nil {
			return nil, err
		}
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(bytes)
	if err != nil {
		return nil, err
	}
	publicKey, err := ioutil.ReadFile(a.Settings.PrivateKeyPath + ".pub")
	if err != nil {
		return nil, err
	}
	body.PublicKey = string(publicKey)

	headers := a.Settings.HTTPManager.GetHeaders(a.Settings.SessionToken, a.Settings.Version, a.Settings.Pod, a.Settings.UsersID)
	message := fmt.Sprintf("%s&%s", headers["X-Request-Nonce"][0], headers["X-Request-Timestamp"][0])
	hashedMessage := sha256.Sum256([]byte(message))
	signature, err := privateKey.Sign(rand.Reader, hashedMessage[:], crypto.SHA256)
	if err != nil {
		return nil, err
	}
	body.Signature = base64.StdEncoding.EncodeToString(signature)

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	var mfaQuery string
	if len(mfaType) > 0 {
		mfaQuery = fmt.Sprintf("?mfaType=%s", mfaType)
	}
	resp, statusCode, err := a.Settings.HTTPManager.Post(b, fmt.Sprintf("%s%s/auth/signin/key%s", a.Settings.AuthHost, a.Settings.AuthHostVersion, mfaQuery), headers)
	if err != nil {
		return nil, err
	}
	signinResp := &signinResponse{}
	return signinResp, a.Settings.HTTPManager.ConvertResp(resp, statusCode, signinResp)
}

func (a *SAuth) mfaSignin(mfaID, preferredMode, challenge string) (*models.User, bool, error) {
	var token string
	fmt.Println("This account has two-factor authentication enabled.")
	if preferredMode == "u2f" {
		data, err := base64.URLEncoding.DecodeString(challenge)
		if err != nil {
			return nil, false, err
		}
		sr := &u2fSignRequest{}
		err = json.Unmarshal(data, sr)
		if err != nil {
			return nil, false, err
		}
		res, err := signU2fRequest(sr)
		if err != nil {
			return nil, true, fmt.Errorf("There was an error communicating with your u2f device, or you do not have one plugged in right now.")
		}
	}
	if preferredMode != "u2f" {
		token = a.Prompts.OTP(preferredMode)
	}
	headers := a.Settings.HTTPManager.GetHeaders(a.Settings.SessionToken, a.Settings.Version, a.Settings.Pod, a.Settings.UsersID)
	b, err := json.Marshal(struct {
		OTP string `json:"otp"`
	}{OTP: token})
	if err != nil {
		return nil, err
	}
	resp, statusCode, err := a.Settings.HTTPManager.Post(b, fmt.Sprintf("%s%s/auth/signin/mfa/%s", a.Settings.AuthHost, a.Settings.AuthHostVersion, mfaID), headers)
	user := &models.User{}
	err = a.Settings.HTTPManager.ConvertResp(resp, statusCode, user)
	if err != nil {
		return nil, err
	}
	return user, err
}

// Signout signs out a user by their session token.
func (a *SAuth) Signout() error {
	headers := a.Settings.HTTPManager.GetHeaders(a.Settings.SessionToken, a.Settings.Version, a.Settings.Pod, a.Settings.UsersID)
	resp, statusCode, err := a.Settings.HTTPManager.Delete(nil, fmt.Sprintf("%s%s/auth/signout", a.Settings.AuthHost, a.Settings.AuthHostVersion), headers)
	if err != nil {
		return err
	}
	return a.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}

// Verify verifies if a given session token is still valid or not. If it is
// valid, the returned error will be nil.
func (a *SAuth) Verify() (*models.User, error) {
	headers := a.Settings.HTTPManager.GetHeaders(a.Settings.SessionToken, a.Settings.Version, a.Settings.Pod, a.Settings.UsersID)
	resp, statusCode, err := a.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/auth/verify", a.Settings.AuthHost, a.Settings.AuthHostVersion), headers)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = a.Settings.HTTPManager.ConvertResp(resp, statusCode, &user)
	if err != nil {
		return nil, err
	}
	a.Settings.UsersID = user.UsersID
	user.SessionToken = a.Settings.SessionToken
	return &user, nil
}

func signU2fRequest(req *u2fSignRequest) (*u2f.AuthenticateResponse, error) {
	devices := u2f.Devices()
	openDevices := []Device{}
	for i, device := range devices {
		err := device.Open()
		if err == nil {
			openDevices = append(openDevices, devices[i])
			defer func(i int) {
				devices[i].Close()
			}(i)
		}
	}
	if len(openDevices) == 0 {
		return nil, fmt.Errorf("no available devices")
	}
	if len(req.RegisteredKeys) == 0 {
		return nil, fmt.Errorf("no registration data")
	}
	prompted := false
	timeout := time.After(time.Second * 6)
	interval := time.NewTicker(time.Millisecond * 250)
	defer interval.Stop()
	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("no response after 6 seconds")
		case <-interval.C:
			for _, device := range openDevices {
				response, err := device.Authenticate(u2f.AuthenticateRequest{
					Challenge: req.Challenge,
					AppId:     req.AppID,
					FacetId:   req.AppID,
					KeyHandle: req.RegisteredKeys[0].KeyHandle,
				})
				if err == nil {
					return response, nil
				} else if _, ok := err.(u2f.TestOfUserPresenceRequiredError); ok && !prompted {
					fmt.Println("Touch your usb device:")
					prompted = true
				} else {
					return nil, fmt.Errorf("unknown device error")
				}
			}
		}
	}
}
