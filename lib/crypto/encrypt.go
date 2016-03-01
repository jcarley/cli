package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
)

// EncryptFile takes in an open plaintext file and encrypts it to a temporary
// location based on the key and IV. It is up to the caller to ensure the
// encrypted file is deleted after it's used. The passed in key and iv should
// *NOT* be base64 encoded or hex encoded.
/*func (c *SCrypto) EncryptFile(plainFilePath string, key, iv []byte) (string, error) {
	plainFile, err := os.Open(plainFilePath)
	defer plainFile.Close()
	if err != nil {
		return "", err
	}

	file, err := os.OpenFile(fmt.Sprintf("%s.encr", plainFilePath), os.O_CREATE|os.O_RDWR, 0600)
	defer file.Close()
	if err != nil {
		return "", err
	}

	aes, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	stream := cipher.NewCTR(aes, iv)
	for {
		chunk := make([]byte, ChunkSize)
		read, _ := plainFile.Read(chunk)
		if read == 0 {
			break
		}
		ct := make([]byte, ChunkSize)
		stream.XORKeyStream(ct, chunk)
		file.Write(ct)
	}
	return file.Name(), nil
}*/

// GCM encryption version
func (c *SCrypto) EncryptFile(plainFilePath string, key []byte) (string, [][]byte, error) {
	if len(key) != 32 {
		return "", nil, errors.New("Invalid key length. Keys must be 256 bits")
	}
	plainFile, err := os.Open(plainFilePath)
	defer plainFile.Close()
	if err != nil {
		return "", nil, err
	}

	file, err := os.OpenFile(fmt.Sprintf("%s.encr", plainFilePath), os.O_CREATE|os.O_RDWR, 0600)
	defer file.Close()
	if err != nil {
		return "", nil, err
	}

	aes, err := aes.NewCipher(key)
	if err != nil {
		return "", nil, err
	}
	aesgcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", nil, err
	}

	ivs := [][]byte{}
	// this isn't great, but it'll work
	// might be a better idea to take a look at the source and come up with a streamable version
	// or perhaps use AES CTR with an HMAC on top of it.
	// just be sure to use two different keys for these
	for {
		iv := make([]byte, IVSize)
		rand.Read(iv)
		ivs = append(ivs, iv)
		chunk := make([]byte, ChunkSize)
		read, _ := plainFile.Read(chunk)
		if read == 0 {
			break
		}
		encrChunk := aesgcm.Seal(nil, iv, chunk, nil)
		file.Write(encrChunk)
	}
	return file.Name(), ivs, nil
}
