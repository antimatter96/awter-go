package common

import (
	"net/http"

	"github.com/antimatter96/awter-go/constants"
	"github.com/gorilla/csrf"
)

var CSRFMiddleware func(http.Handler) http.Handler

var BcryptCost int

type key int

// The constants for context
const (
	SessionIDKey key = 1
	UserIDKey    key = 2
)

// Init is used to initialize all things
func InitCommon() {
	BcryptCost = int(constants.Value("bcrypt-cost").(float64))
	if BcryptCost > 31 {
		panic("Bcrypt Cost Exceeded")
	}
}

func InitCSRF(errorHandler http.HandlerFunc) {
	CSRFMiddleware = csrf.Protect(
		[]byte(constants.Value("csrf-auth-key").(string)),
		csrf.FieldName("_csrf_token"),
		csrf.CookieName("_csrf_token"),
		csrf.Secure(constants.ENVIRONMENT != "dev"),
		csrf.ErrorHandler(errorHandler),
	)
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
