package crypto

import (
	"fmt"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/gcm/gcm"
)

// EncryptFile takes in an open plaintext file and encrypts it to a temporary
// location based on the key and IV. It is up to the caller to ensure the
// encrypted file is deleted after it's used. The passed in key and iv should
// *NOT* be base64 encoded or hex encoded.
func (c *SCrypto) EncryptFile(plainFilePath string, key, iv []byte) (string, error) {
	if len(key) != KeySize {
		return "", fmt.Errorf("Invalid key length. Keys must be %d bytes", KeySize)
	}
	if len(iv) != IVSize {
		return "", fmt.Errorf("Invalid IV length. IVs must be %d bytes", IVSize)
	}
	outputFile, err := ioutil.TempFile("", "encr")
	if err != nil {
		return "", err
	}
	outputFile.Close()
	logrus.Debugf("temp file %s", outputFile.Name())

	err = gcm.EncryptFile(plainFilePath, outputFile.Name(), key, iv, c.Unhex([]byte(gcm.AAD), AADSize))
	if err != nil {
		return "", err
	}
	return outputFile.Name(), nil
}
