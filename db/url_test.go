package db

import (
	"fmt"
	"testing"

	"github.com/antimatter96/awter-go/constants"
	url "github.com/antimatter96/awter-go/db/url"
)

var testServiceR url.Service
var testServiceM url.Service

func init() {
	if err := constants.Init("config-test", "../", "..", "./"); err != nil {
		fmt.Printf("cant initialize constants : %v", err)
	}
	InitRedis()
	InitMySQL()

	testServiceM = NewURLInterfaceMySQL()
	testServiceR = NewURLInterfaceRedis()
}

var err error

func BenchmarkCreateRedis(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		urlObj := &url.ShortURL{
			Short:         "short",
			Salt:          "SaltIsVeryLong, VeryLong",
			Nonce:         "NonceIsVeryLong, VeryLong",
			PasswordHash:  "PasswordHashIsVeryLong, VeryLong",
			EncryptedLong: "EncryptedLongIsVeryLong, VeryLong",
		}
		err = testServiceR.Create(*urlObj)
	}
}

func BenchmarkCreateMySQL(b *testing.B) {

	for n := 0; n < b.N; n++ {
		urlObj := &url.ShortURL{
			Short:         "short",
			Salt:          "SaltIsVeryLong, VeryLong",
			Nonce:         "NonceIsVeryLong, VeryLong",
			PasswordHash:  "PasswordHashIsVeryLong, VeryLong",
			EncryptedLong: "EncryptedLongIsVeryLong, VeryLong",
		}
		err = testServiceM.Create(*urlObj)
	}
	sqlDB.Exec("DELETE FROM SHORT_URLS")
}
