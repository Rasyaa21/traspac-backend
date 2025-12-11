package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

func GenerateRandomToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", nil
	}
	return hex.EncodeToString(b), nil
}


func GenerateUppercaseSixDigitOTP() (string, error) {
	const otpCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)

	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(otpCharset))))
		if err != nil {
			return "", err
		}
		result[i] = otpCharset[num.Int64()]
	}
	return string(result), nil
}