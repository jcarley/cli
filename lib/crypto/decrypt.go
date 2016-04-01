package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/gcm/gcm"
)

// DecryptFile takes in an ecrypted file and decrypts it to the given
// output path based on the Key and IV. The Key and IV should be the hex and
// base64 encoded version
func (c *SCrypto) DecryptFile(encryptedFilePath, key, iv, outputFilePath string) error {
	legacy, err := isLegacy(encryptedFilePath)
	if err != nil {
		return err
	}
	logrus.Debugf("Legacy encryption scheme detected? %t", legacy)
	if legacy {
		return c.decryptLegacy(encryptedFilePath, key, iv, outputFilePath)
	}
	return gcm.DecryptFile(encryptedFilePath, outputFilePath, c.Unhex([]byte(key), KeySize), c.Unhex([]byte(iv), IVSize), c.Unhex([]byte(gcm.AAD), AADSize))
}

func (c *SCrypto) decryptLegacy(encryptedFilePath, key, iv, outputFilePath string) error {
	encryptedFile, err := os.Open(encryptedFilePath)
	defer encryptedFile.Close()
	if err != nil {
		return err
	}

	sizeBytes := make([]byte, 8)
	binary.Read(encryptedFile, binary.LittleEndian, sizeBytes)
	var origSize int64
	binary.Read(bytes.NewBuffer(sizeBytes), binary.LittleEndian, &origSize)

	block, err := aes.NewCipher(c.Unhex(c.Base64Decode([]byte(key), KeySize*2), KeySize))
	if err != nil {
		return err
	}

	file, err := os.OpenFile(outputFilePath, os.O_RDWR, 0600)
	if os.IsNotExist(err) {
		file, err = os.OpenFile(outputFilePath, os.O_CREATE|os.O_RDWR, 0600)
	}
	if err != nil {
		return err
	}
	defer file.Close()

	decrypter := cipher.NewCBCDecrypter(block, c.Unhex(c.Base64Decode([]byte(iv), LegacyIVSize*2), LegacyIVSize))
	chunkSize := aes.BlockSize
	var previousBlock []byte
	for {
		chunk := make([]byte, chunkSize)
		read, errRead := encryptedFile.Read(chunk)
		if errRead != nil && errRead != io.EOF {
			return errRead
		}
		if read%chunkSize != 0 {
			return fmt.Errorf("Logs unavailable for this job")
		}
		if read == 0 {
			break
		}
		plainChunk := make([]byte, read)
		decrypter.CryptBlocks(plainChunk, chunk[:read])
		if previousBlock != nil {
			file.Write(previousBlock)
		}
		previousBlock = plainChunk
	}
	_, err = file.Write(previousBlock[:origSize%aes.BlockSize])
	return err
}

func isLegacy(encryptedFilePath string) (bool, error) {
	stat, err := os.Stat(encryptedFilePath)
	if err != nil {
		return false, nil
	}
	encryptedFile, err := os.Open(encryptedFilePath)
	if err != nil {
		return false, err
	}
	defer encryptedFile.Close()
	sizeBytes := make([]byte, 8)
	err = binary.Read(encryptedFile, binary.LittleEndian, sizeBytes)
	if err != nil {
		return false, err
	}
	var origSize int64
	err = binary.Read(bytes.NewBuffer(sizeBytes), binary.LittleEndian, &origSize)
	if err != nil {
		return false, err
	}
	if origSize+8+(aes.BlockSize-origSize%aes.BlockSize) == stat.Size() {
		return true, nil
	}
	return false, nil
}
