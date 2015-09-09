package helpers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
)

// DecryptFile takes in an ecrypted file and decrypts it to the given
// output path based on the Key and IV. The Key and IV should be the hex and
// base64 encoded version
func DecryptFile(encryptedFilePath string, key string, iv string, outputFilePath string) {
	encryptedFile, encrFileErr := os.Open(encryptedFilePath)
	defer encryptedFile.Close()
	if encrFileErr != nil {
		fmt.Println(encrFileErr.Error())
		os.Exit(1)
	}
	sizeBytes := make([]byte, 8)
	binary.Read(encryptedFile, binary.LittleEndian, sizeBytes)
	var origSize int64
	binary.Read(bytes.NewBuffer(sizeBytes), binary.LittleEndian, &origSize)
	block, err := aes.NewCipher(Unhex(Base64Decode([]byte(key))))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	file, fileErr := os.OpenFile(outputFilePath, os.O_RDWR, 0644)
	if os.IsNotExist(fileErr) {
		file, fileErr = os.OpenFile(outputFilePath, os.O_CREATE|os.O_RDWR, 0644)
	}
	if fileErr != nil {
		fmt.Println(fileErr.Error())
		os.Exit(1)
	}
	defer file.Close()

	decrypter := cipher.NewCBCDecrypter(block, Unhex(Base64Decode([]byte(iv))))
	chunkSize := 24 * 1024
	for {
		chunk := make([]byte, chunkSize)
		read, _ := encryptedFile.Read(chunk)
		if read == 0 {
			break
		}
		plainChunk := make([]byte, read)
		decrypter.CryptBlocks(plainChunk, chunk[:read]) // only decrypt the amount we read
		file.Write(plainChunk)
	}
	file.Truncate(origSize)
}

// EncryptFile takes in an open plaintext file and encrypts it to a temporary
// location based on the key and IV. It is up to the caller to ensure the
// encrypted file is deleted after it's used. The passed in key and iv should
// *NOT* be base64 encoded or hex encoded.
func EncryptFile(plainFilePath string, key []byte, iv []byte, importRequiresLength bool) string {
	plainFile, plainFileErr := os.Open(plainFilePath)
	defer plainFile.Close()
	if plainFileErr != nil {
		fmt.Println(plainFileErr.Error())
		os.Exit(1)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// TODO think of something more unique than this without opening a temp file,
	// taking the name, then opening it with this method.
	// In the end I need to open a temp file with 0644
	file, fileErr := os.OpenFile(fmt.Sprintf("%s.encr", plainFilePath), os.O_CREATE|os.O_RDWR, 0644)
	defer file.Close()
	if fileErr != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if importRequiresLength {
		// pack the file size in the first 8 bytes
		// would be nice to figure out how to do this and the Decrypt at least in a
		// similar manner
		stat, _ := plainFile.Stat()
		var origSize = uint64(stat.Size())
		sizeBytes := make([]byte, 8)
		binary.PutUvarint(sizeBytes, origSize)
		binary.Write(file, binary.LittleEndian, sizeBytes)
	}

	encrypter := cipher.NewCBCEncrypter(block, iv)
	chunkSize := 24 * 1024
	for {
		chunk := make([]byte, chunkSize)
		read, _ := plainFile.Read(chunk)
		if read == 0 {
			break
		}
		encrChunk := make([]byte, chunkSize)
		encrypter.CryptBlocks(encrChunk, chunk)
		file.Write(encrChunk)
	}
	return file.Name()
}

// Hex encode bytes
func Hex(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return bytes.Trim(dst, "\x00") // go is broken, need to figure out where exactly
}

// Unhex bytes
func Unhex(src []byte) []byte {
	dst := make([]byte, hex.DecodedLen(len(src)))
	hex.Decode(dst, src)
	return bytes.Trim(dst, "\x00") // go is broken, need to figure out where exactly
}

// Base64Encode bytes
func Base64Encode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return bytes.Trim(dst, "\x00") // go is broken, need to figure out where exactly
}

// Base64Decode bytes
func Base64Decode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	base64.StdEncoding.Decode(dst, src)
	return bytes.Trim(dst, "\x00") // go is broken, need to figure out where exactly
}
