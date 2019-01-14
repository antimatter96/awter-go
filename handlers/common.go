package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"../cache"
	"../constants"
	"../db"
	"github.com/julienschmidt/httprouter"
)

// All the different errors
var (
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

var letterRunes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var urlService db.URLService

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

var bcryptCost int

type key int

// The constants for context
var (
	SessionIDKey key = 1
	UserIDKey    key = 2
)

type HandlerWithContext func(context.Context, http.ResponseWriter, *http.Request, httprouter.Params)

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
				sessionId := value["sessionid"]
				if sessionId != "" {
					fmt.Println(sessionId)
					userID, err := cache.GetSessionValue(sessionId, "userId")
					if err == nil {
						ctx = context.WithValue(ctx, UserIDKey, &userID)
					}
					ctx = context.WithValue(ctx, SessionIDKey, &sessionId)
				}
			}
		}
		next(ctx, w, r, ps)
	}
}

func MiddlewareAllowOnlyAuth(next HandlerWithContext) HandlerWithContext {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		sessionId := ctx.Value(SessionIDKey).(string)
		if sessionId == "" {
			http.Redirect(w, r, "/login?url_next="+r.URL.String(), http.StatusSeeOther)
		} else {
			userID, err := cache.GetSessionValue(sessionId, "userId")
			if err != nil {
				http.Redirect(w, r, "/login?url_next="+r.URL.String(), http.StatusSeeOther)
			} else {
				ctx = context.WithValue(ctx, UserIDKey, &userID)
				next(ctx, w, r, ps)
			}
		}
	}
}
