package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRand, err := rand.Read(b)
	if err != nil {
		return []byte{}, fmt.Errorf("creating: %v", err)
	}

	if n > nRand {
		return []byte{}, fmt.Errorf("if nRand equals n: %v", err)
	}

	return b, nil
}

func ToString(n int) (string, error) {
	s, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("creating: %v", err)
	}

	return base64.URLEncoding.EncodeToString(s), nil
}
