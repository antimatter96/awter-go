package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"../constants"
	"../db"
	"../db/url"
	"github.com/julienschmidt/httprouter"
)

// All the different errors
const (
	constErrEmailMissing        string = "Email Not Present"
	constErrPasswordMissing     string = "Password Not Present"
	constErrNotRegistered       string = "No records found"
	constErrInternalError       string = "An Error Occured"
	constErrPasswordMatchFailed string = "Passwords do not match"
	constErrEmailTaken          string = "Email Taken"
	constErrURLMissing          string = "URL Missing"
	constErrPasswordTooShort    string = "Password too short"
	constErrURLNotPresent       string = "URL not present"
)

var bcryptCost int

type key int

// The constants for context
const (
	SessionIDKey key = 1
	UserIDKey    key = 2
)

var letterRunes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var urlService url.URLService
var oneHour time.Duration = 720 * time.Minute

// Init is used to initialize all things
func Init() {
	initCookie()
	parseTemplates()

	urlService = db.NewURLInterfaceRedis()

	bcryptCost = int(constants.Value("bcrypt-cost").(float64))
	if bcryptCost > 31 {
		panic("Bcrypt Cost Exceeded")
	}
}

func generateRandomString(length int) (string, error) {
	x := make([]byte, length)
	_, err := rand.Read(x)

	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(x), nil
}

func getSalt(length int) (string, error) {
	x := make([]byte, length)
	_, err := rand.Read(x)

	if err != nil {
		return "", err
	}
	return string(x), nil
}

// HandlerWithContext is a custom request handler
type HandlerWithContext func(context.Context, http.ResponseWriter, *http.Request, httprouter.Params)

// Wrapper is used to
func Wrapper(lead HandlerWithContext) httprouter.Handle {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		lead(ctx, w, r, ps)
	}
}

// ExtractSessionID sets the userId in context
func ExtractSessionID(next HandlerWithContext) HandlerWithContext {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		httpCookie, err := r.Cookie("sessionid")
		if err == nil {
			value := make(map[string]string)
			err = cookie.Decode("sessionid", httpCookie.Value, &value)
			if err == nil {
				fmt.Println(value)
				if sessionID := value["sessionid"]; sessionID != "" {
					ctx = context.WithValue(ctx, SessionIDKey, &sessionID)
				}
			}
		}

		if _, isString := ctx.Value(SessionIDKey).(string); !isString {
			newSessionID, _ := generateRandomString(32)
			encodedValue, err := cookie.Encode("sessionid", newSessionID)
			if err == nil {
				cookie := &http.Cookie{
					Name:     "sessionid",
					Value:    encodedValue,
					Path:     "/",
					HttpOnly: true,
					Expires:  time.Now().Add(oneHour),
				}
				http.SetCookie(w, cookie)
			}
		}
		next(ctx, w, r, ps)
	}
}
