// Pacakge server stores all the constants used all over the service

package server

import (
	"fmt"
	"net/http"

	"github.com/antimatter96/awter-go/customcrypto"
	"github.com/antimatter96/awter-go/db/url"
	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

func (server *server) mainGet(w http.ResponseWriter, r *http.Request) {
	renderParams := r.Context().Value(ctxKeyRenderParms).(*map[string]interface{})
	server.shortnerTemplate.Execute(w, renderParams)
}

func (server *server) shortPost(w http.ResponseWriter, r *http.Request) {
	renderParams := r.Context().Value(ctxKeyRenderParms).(*map[string]interface{})

	errParseForm := r.ParseForm()

	if errParseForm != nil {
		fmt.Println("parse error", errParseForm)
		(*renderParams)["error"] = ErrInternalError
		server.shortnerTemplate.Execute(w, renderParams)
		return
	}

	urlNext := "/"
	if len(r.Form["url_next"]) > 0 {
		urlNext = r.Form["url_next"][0]
		(*renderParams)["url_next"] = urlNext
	}

	link := r.FormValue("url")
	passwordProtect := r.FormValue("passwordProtect") == "on"
	password := r.FormValue("password")

	if link == "" || !govalidator.IsURL(link) {
		(*renderParams)["error"] = ErrURLMissing
		server.shortnerTemplate.Execute(w, renderParams)
		return
	}

	if passwordProtect && password == "" {
		(*renderParams)["error"] = ErrPasswordMissing
		(*renderParams)["url_next"] = urlNext
		server.shortnerTemplate.Execute(w, renderParams)
		return
	} else if !passwordProtect {
		password = "default"
	}

	var shortURL string
	var err error
	for {
		shortURL, err = customcrypto.GenerateRandomString(6)
		if err != nil {
			(*renderParams)["error"] = ErrInternalError
			server.shortnerTemplate.Execute(w, renderParams)
			return
		}
		present, err := server.urlService.Present(shortURL)
		if err != nil {
			(*renderParams)["error"] = ErrInternalError
			server.shortnerTemplate.Execute(w, renderParams)
			return
		}
		if !present {
			break
		}
	}

	var hashedPassword string

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), server.BcryptCost)
	if err != nil {
		(*renderParams)["error"] = ErrInternalError
		server.shortnerTemplate.Execute(w, renderParams)
		return
	}
	hashedPassword = string(hashed)

	nonce, salt, encryptedLong, err := customcrypto.Encrypt(password, link)
	if err != nil {
		(*renderParams)["error"] = ErrInternalError
		server.shortnerTemplate.Execute(w, renderParams)
		return
	}

	urlObj := &url.ShortURL{Short: shortURL, Nonce: nonce, Salt: salt, EncryptedLong: encryptedLong, PasswordHash: hashedPassword}

	err = server.urlService.Create(*urlObj)
	if err != nil {
		(*renderParams)["error"] = ErrInternalError
		server.shortnerTemplate.Execute(w, renderParams)
		return
	}

	(*renderParams)["shortURL"] = shortURL
	(*renderParams)["passwordProtect"] = passwordProtect
	(*renderParams)["password"] = password
	(*renderParams)["longURL"] = link

	server.createdTemplate.Execute(w, renderParams)
}

func (server *server) elongateGet(w http.ResponseWriter, r *http.Request) {
	server.elongateTemplate.Execute(w, nil)
}

func (server *server) elongatePost(w http.ResponseWriter, r *http.Request) {
	server.shortnerTemplate.Execute(w, nil)
}
