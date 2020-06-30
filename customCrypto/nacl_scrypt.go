package customcrypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

const (
	naclScryptKeyLen   = 32
	naclScryptNonceLen = 24
	naclScryptSaltLen  = 12
)

// NaclScrypt is all that you get
type NaclScrypt struct{}

// Encrypt is
func (cc *NaclScrypt) Encrypt(password, data string) (nonceToSave, saltToSave, encryptedToSave string, err error) {
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

	salt := make([]byte, naclScryptSaltLen)
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}

	ek, err := scrypt.Key(pwdBuff, salt, 16384, 8, 1, naclScryptKeyLen)
	if err != nil {
		panic(err)
	}

	var ekLimited [naclScryptKeyLen]byte
	copy(ekLimited[:], ek)

	nonce := make([]byte, naclScryptNonceLen)
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	var secretKeyLimited [naclScryptNonceLen]byte
	copy(secretKeyLimited[:], nonce)

	encrypted := secretbox.Seal(nil, []byte(data), &secretKeyLimited, &ekLimited)

	nonceToSave = hex.EncodeToString(nonce)
	saltToSave = hex.EncodeToString(salt)
	encryptedToSave = hex.EncodeToString(encrypted)

	return
}

// Decrypt is
func (cc *NaclScrypt) Decrypt(password, data, nonceString, saltString string) (long string, err error) {
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
	dk, err := scrypt.Key(pwdBuff, saltBytes, 16384, 8, 1, naclScryptKeyLen)
	if err != nil {
		panic(err)
	}

	nonceBytes, err := hex.DecodeString(nonceString)
	if err != nil {
		panic(err)
	}

	var dkLimited [naclScryptKeyLen]byte
	copy(dkLimited[:], dk)

	var secretKeyLimited [naclScryptNonceLen]byte
	copy(secretKeyLimited[:], nonceBytes)

	encBytes, err := hex.DecodeString(data)
	if err != nil {
		panic(err)
	}

	decrypted, ok := secretbox.Open(nil, []byte(encBytes), &secretKeyLimited, &dkLimited)
	if !ok {
		err = fmt.Errorf("cant decode")
		panic(err)
	}
	long = string(decrypted)
	return
}
