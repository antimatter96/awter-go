package server

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/securecookie"
)

func (server *Server) initCookie(hashKeyString, blockKeyString string) {
	hashKey := []byte(hashKeyString)
	blockKey := []byte(blockKeyString)

	server.cookie = securecookie.New(hashKey, blockKey)
	server.cookie.MaxAge(43200)
}

func (server *Server) initCSRF(csrfAuthKey string, dev bool) {
	server.csrfMiddleware = csrf.Protect(
		[]byte(csrfAuthKey),
		csrf.FieldName("_csrf_token"),
		csrf.CookieName("_csrf_token"),
		csrf.Secure(!dev),
		//csrf.ErrorHandler(errorHandler),
	)
}
