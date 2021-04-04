package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
)

// EncodeMD5 md5 encryption
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
func EncodeSha1(value string) string {
	h := sha1.New()
	io.WriteString(h, value)
	return string(h.Sum(nil))
}
