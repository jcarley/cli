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

	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
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
	signinResp, err := f()
	if err != nil {
		return nil, err
	}

	var user *models.User

	if signinResp.MFAID != "" {
		user, err = a.mfaSignin(signinResp.MFAID, signinResp.MFAPreferredMode)
		if err != nil {
			return nil, err
		}
	} else {
		user = signinResp.toUser()
	}

	a.Settings.UsersID = user.UsersID
	a.Settings.Username = user.Username
	a.Settings.SessionToken = user.SessionToken
	return user, nil
}

type signinResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	SessionToken     string `json:"sessionToken"`
	MFAID            string `json:"mfaID"`
	MFAPreferredMode string `json:"mfaPreferredType"`
}

func (sr *signinResponse) toUser() *models.User {
	return &models.User{
		UsersID:      sr.ID,
		Username:     sr.Name,
		Email:        sr.Email,
		SessionToken: sr.SessionToken,
	}
}

func (a *SAuth) signInWithCredentials() (*signinResponse, error) {
	login := models.Login{
		Identifier: a.Settings.Username,
		Password:   a.Settings.Password,
	}
	if a.Settings.Username == "" || a.Settings.Password == "" {
		username, password, err := a.Prompts.UsernamePassword()
		if err != nil {
			return nil, err
		}
		login = models.Login{
			Identifier: username,
			Password:   password,
		}
	}

	b, err := json.Marshal(login)
	if err != nil {
		return nil, err
	}
	headers := httpclient.GetHeaders(a.Settings.SessionToken, a.Settings.Version, a.Settings.Pod, a.Settings.UsersID)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/auth/signin", a.Settings.AuthHost, a.Settings.AuthHostVersion), headers)
	if err != nil {
		return nil, err
	}
	signinResp := &signinResponse{}
	return signinResp, httpclient.ConvertResp(resp, statusCode, signinResp)
}

func (a *SAuth) signInWithKey() (*signinResponse, error) {
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

	headers := httpclient.GetHeaders(a.Settings.SessionToken, a.Settings.Version, a.Settings.Pod, a.Settings.UsersID)
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
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/auth/signin/key", a.Settings.AuthHost, a.Settings.AuthHostVersion), headers)
	if err != nil {
		return nil, err
	}
	signinResp := &signinResponse{}
	return signinResp, httpclient.ConvertResp(resp, statusCode, signinResp)
}

func (a *SAuth) mfaSignin(mfaID string, preferredMode string) (*models.User, error) {
	token := a.Prompts.OTP(preferredMode)
	headers := httpclient.GetHeaders(a.Settings.SessionToken, a.Settings.Version, a.Settings.Pod, a.Settings.UsersID)
	b, err := json.Marshal(struct {
		OTP string `json:"otp"`
	}{OTP: token})
	if err != nil {
		return nil, err
	}
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/auth/signin/mfa/%s", a.Settings.AuthHost, a.Settings.AuthHostVersion, mfaID), headers)
	user := &models.User{}
	err = httpclient.ConvertResp(resp, statusCode, user)
	if err != nil {
		return nil, err
	}
	return user, err
}

// Signout signs out a user by their session token.
func (a *SAuth) Signout() error {
	headers := httpclient.GetHeaders(a.Settings.SessionToken, a.Settings.Version, a.Settings.Pod, a.Settings.UsersID)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/auth/signout", a.Settings.AuthHost, a.Settings.AuthHostVersion), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}

// Verify verifies if a given session token is still valid or not. If it is
// valid, the returned error will be nil.
func (a *SAuth) Verify() (*models.User, error) {
	headers := httpclient.GetHeaders(a.Settings.SessionToken, a.Settings.Version, a.Settings.Pod, a.Settings.UsersID)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/auth/verify", a.Settings.AuthHost, a.Settings.AuthHostVersion), headers)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = httpclient.ConvertResp(resp, statusCode, &user)
	if err != nil {
		return nil, err
	}
	a.Settings.UsersID = user.UsersID
	return &user, nil
}
