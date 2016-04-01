package crypto

import (
	"fmt"

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
	outputFilePath := fmt.Sprintf("%s.encr", plainFilePath)
	err := gcm.EncryptFile(plainFilePath, outputFilePath, key, iv, c.Unhex([]byte(gcm.AAD), AADSize))
	if err != nil {
		return "", err
	}
	return outputFilePath, nil
}
