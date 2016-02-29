package crypto

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
)

// Hex encode bytes
func (c *SCrypto) Hex(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	if len(dst) != 32 && len(dst) != 64 {
		dst = bytes.Trim(dst, "\x00")
	}
	return dst
}

// Unhex bytes
func (c *SCrypto) Unhex(src []byte) []byte {
	dst := make([]byte, hex.DecodedLen(len(src)))
	hex.Decode(dst, src)
	if len(dst) != 32 && len(dst) != 64 {
		dst = bytes.Trim(dst, "\x00")
	}
	return dst
}

// Base64Encode bytes
func (c *SCrypto) Base64Encode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	if len(dst) != 32 && len(dst) != 64 {
		dst = bytes.Trim(dst, "\x00")
	}
	return dst
}

// Base64Decode bytes
func (c *SCrypto) Base64Decode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	base64.StdEncoding.Decode(dst, src)
	if len(dst) != 32 && len(dst) != 64 {
		dst = bytes.Trim(dst, "\x00")
	}
	return dst
}
