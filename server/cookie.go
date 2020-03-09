package server

import "github.com/gorilla/securecookie"

func (server *server) initCookie(hashKeyString, blockKeyString string) {
	hashKey := []byte(hashKeyString)
	blockKey := []byte(blockKeyString)

	server.Cookie = securecookie.New(hashKey, blockKey)
	server.Cookie.MaxAge(43200)
}
