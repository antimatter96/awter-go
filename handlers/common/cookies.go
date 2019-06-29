package common

import (
	"github.com/antimatter96/awter-go/constants"
	"github.com/gorilla/securecookie"
)

var Cookie *securecookie.SecureCookie

func initCookie() {
	cookieConfig, _ := constants.Value("cookie").(map[string]interface{})

	hashKeyString, _ := cookieConfig["hashkey"].(string)
	blockKeyString, _ := cookieConfig["blockkey"].(string)

	hashKey := []byte(hashKeyString)
	blockKey := []byte(blockKeyString)

	Cookie = securecookie.New(hashKey, blockKey)
	Cookie.MaxAge(43200)
}
