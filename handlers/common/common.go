package common

import (
	"net/http"
)

var CSRFMiddleware func(http.Handler) http.Handler

type key int

// The constants for context
const (
	SessionIDKey key = 1
	UserIDKey    key = 2
)

// HandlerWithContext is a custom request handler
// type HandlerWithContext func(context.Context, http.ResponseWriter, *http.Request)

// Wrapper is used to
// func Wrapper(lead HandlerWithContext) http.Handler {
// 	ctx := context.Background()
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		lead(ctx, w, r)
// 	}
// }
