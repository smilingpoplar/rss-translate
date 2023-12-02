package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

func MD5(data []byte) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, bytes.NewReader(data)); err != nil {
		return "", fmt.Errorf("error hashing data: %w", err)
	}

	md5Bytes := hash.Sum(nil)
	md5Str := hex.EncodeToString(md5Bytes)
	return md5Str, nil
}
