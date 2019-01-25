package common

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"../../constants"
)

// All the different errors
const (
	ConstErrEmailMissing        string = "Email Not Present"
	ConstErrPasswordMissing     string = "Password Not Present"
	ConstErrNotRegistered       string = "No records found"
	ConstErrInternalError       string = "An Error Occured"
	ConstErrPasswordMatchFailed string = "Passwords do not match"
	ConstErrEmailTaken          string = "Email Taken"
	ConstErrURLMissing          string = "URL Missing"
	ConstErrPasswordTooShort    string = "Password too short"
	ConstErrURLNotPresent       string = "URL not present"
)

var BcryptCost int

type key int

// The constants for context
const (
	SessionIDKey key = 1
	UserIDKey    key = 2
)

var oneHour time.Duration = 720 * time.Minute

// Init is used to initialize all things
func InitCommon() {
	initCookie()

	BcryptCost = int(constants.Value("bcrypt-cost").(float64))
	if BcryptCost > 31 {
		panic("Bcrypt Cost Exceeded")
	}
}

func GenerateRandomString(length int) (string, error) {
	x := make([]byte, length)
	_, err := rand.Read(x)

	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(x), nil
}

func GetSalt(length int) (string, error) {
	x := make([]byte, length)
	_, err := rand.Read(x)

	if err != nil {
		return "", err
	}
	return string(x), nil
}

// HandlerWithContext is a custom request handler
// type HandlerWithContext func(context.Context, http.ResponseWriter, *http.Request)

// Wrapper is used to
// func Wrapper(lead HandlerWithContext) http.Handler {
// 	ctx := context.Background()
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		lead(ctx, w, r)
// 	}
// }
