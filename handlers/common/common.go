package common

import (
	"net/http"

	"github.com/antimatter96/awter-go/constants"
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

// HandlerWithContext is a custom request handler
// type HandlerWithContext func(context.Context, http.ResponseWriter, *http.Request)

// Wrapper is used to
// func Wrapper(lead HandlerWithContext) http.Handler {
// 	ctx := context.Background()
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		lead(ctx, w, r)
// 	}
// }
