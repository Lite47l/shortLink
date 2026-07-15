package service

import (
	"crypto/rand"
	"math/big"
)

const (
	chars  = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890_"
	lenght = 10
)

func Generate() (string, error) {
	result := make([]byte, lenght)
	charsLenght := big.NewInt(int64(len(chars)))

	for i := 0; i < lenght; i++ {
		n, err := rand.Int(rand.Reader, charsLenght)
		if err != nil {
			return "", err
		}
		result[i] = chars[n.Int64()]
	}
	return string(result), nil
}
