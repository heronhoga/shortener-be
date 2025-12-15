package util

import (
	"crypto/rand"
	"math/big"
	"regexp"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomName(n int) (string, error) {
	result := make([]byte, n)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[num.Int64()]
	}
	return string(result), nil
}

func CheckValidLinkName(link string) bool {
		re := regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]*$`)
		return re.MatchString(link)
}
