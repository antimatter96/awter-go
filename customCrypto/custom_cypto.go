// Pacakge customcrypto stores all the constants used all over the service

package customcrypto

import (
	"crypto/rand"
	simplerand "math/rand"
	"time"
)

// CustomCrypto encapsulates all the functions performed by
// our crypto services
type CustomCrypto interface {
	Encrypt(string, string) (string, string, string, error)
	Decrypt(string, string, string, string) (string, error)
}

type PasswordChecker interface {
	IsSame([]byte, []byte) error
	GetHash([]byte) ([]byte, error)
	NoMatch(error) bool
}

var letterRunes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// GenerateRandomString is used to generate a string of given length
// String will follow the pattern [0-9a-zA-Z]+
func GenerateRandomString(length int) (string, error) {

	randomFactor := make([]byte, 2)
	_, err := rand.Read(randomFactor)
	if err != nil {
		return "", err
	}

	simplerand.Seed(time.Now().UnixNano() * int64(randomFactor[0]) * int64(randomFactor[1]))

	arr := make([]byte, length)
	for i := range arr {
		arr[i] = letterRunes[simplerand.Intn(len(letterRunes))]
	}
	return string(arr), nil
}
