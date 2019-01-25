package common

import (
	"../../constants"
	"github.com/gorilla/securecookie"
)

var cookie *securecookie.SecureCookie

func initCookie() {
	cookieConfig, _ := constants.Value("cookie").(map[string]interface{})

	hashKeyString, _ := cookieConfig["hashkey"].(string)
	blockKeyString, _ := cookieConfig["blockkey"].(string)

	hashKey := []byte(hashKeyString)
	blockKey := []byte(blockKeyString)

	cookie = securecookie.New(hashKey, blockKey)
	cookie.MaxAge(43200)
}
