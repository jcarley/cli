package crypto

import "github.com/catalyzeio/gcm/gcm"

// ICrypto
type ICrypto interface {
	DecryptFile(encryptedFilePath, key, iv, outputFilePath string) error
	EncryptFile(plainFilePath string, key, iv []byte) (string, error)
	NewEncryptFileReader(plainFilePath string, key, iv []byte) (*gcm.EncryptFileReader, error)
	Hex(src []byte, maxLen int) []byte
	Unhex(src []byte, maxLen int) []byte
	Base64Encode(src []byte, maxLen int) []byte
	Base64Decode(src []byte, maxLen int) []byte
}

// SCrypto is an implementor of ICrypto
type SCrypto struct{}

// New creates a new instance of ICrypto
func New() ICrypto {
	return &SCrypto{}
}
