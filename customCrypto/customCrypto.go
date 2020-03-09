package lib

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	simplerand "math/rand"
	"time"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

var letterRunes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomString2(length int) (string, error) {

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

func encrypt(password, data string) (nonceToSave, saltToSave, encryptedToSave string, err error) {
	defer func() {
		if errRecovered := recover(); errRecovered != nil {
			if value, isError := errRecovered.(error); isError {
				nonceToSave = ""
				saltToSave = ""
				encryptedToSave = ""
				fmt.Println("Error", value)
			}
		}
	}()
	pwdBuff := []byte(password)

	salt := make([]byte, 12)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}

	ek, err := scrypt.Key(pwdBuff, salt, 16384, 8, 1, 32)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, 24)
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	var ekLimited [32]byte
	copy(ekLimited[:], ek)

	var secretKeyLimited [24]byte
	copy(secretKeyLimited[:], nonce)

	encrypted := secretbox.Seal(nil, []byte(data), &secretKeyLimited, &ekLimited)

	nonceToSave = hex.EncodeToString(nonce)
	saltToSave = hex.EncodeToString(salt)
	encryptedToSave = hex.EncodeToString(encrypted)

	return
}

func decrypt(password, data, nonceString, saltString string) (long string, err error) {
	defer func() {
		if errRecovered := recover(); errRecovered != nil {
			if value, isError := errRecovered.(error); isError {
				long = ""
				fmt.Println("Error", value)
			}
		}
	}()

	saltBytes, err := hex.DecodeString(saltString)
	if err != nil {
		panic(err)
	}

	pwdBuff := []byte(password)
	dk, err := scrypt.Key(pwdBuff, saltBytes, 16384, 8, 1, 32)
	if err != nil {
		panic(err)
	}

	nonceBytes, err := hex.DecodeString(nonceString)
	if err != nil {
		panic(err)
	}

	var dkLimited [32]byte
	copy(dkLimited[:], dk)

	var secretKeyLimited [24]byte
	copy(secretKeyLimited[:], nonceBytes)

	encBytes, err := hex.DecodeString(data)
	if err != nil {
		panic(err)
	}

	decrypted, ok := secretbox.Open(nil, []byte(encBytes), &secretKeyLimited, &dkLimited)
	if !ok {
		err = fmt.Errorf("Cant decode")
		panic(err)
	}
	long = string(decrypted)
	return
}
