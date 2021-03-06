// Pacakge server stores all the constants used all over the service

package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/antimatter96/awter-go/customcrypto"
	"github.com/antimatter96/awter-go/db/url"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/hlog"
)

func (server *Server) mainGet(w http.ResponseWriter, r *http.Request) {
	renderParams := r.Context().Value(ctxKeyRenderParms).(*map[string]interface{})
	server.shortnerTemplate.Execute(w, renderParams)
}

func (server *Server) shortPost(w http.ResponseWriter, r *http.Request) {
	renderParams := r.Context().Value(ctxKeyRenderParms).(*map[string]interface{})

	errParseForm := r.ParseForm()

	if errParseForm != nil {
		fmt.Println("parse error", errParseForm)
		(*renderParams)["error"] = ErrInternalError
		server.shortnerTemplate.Execute(w, renderParams)
		return
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
		server.shortnerTemplate.Execute(w, renderParams)
		return
	} else if !passwordProtect {
		password = "default"
	}

	var shortURL string
	var err error
	attempt := 0
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
		attempt++
		if attempt > 3-1 {
			(*renderParams)["error"] = ErrInternalError
			server.shortnerTemplate.Execute(w, renderParams)
			return
		}
	}

	hlog.FromRequest(r).Info().Msg("bcrypt encrypt start")
	hashed, err := server.passwordChecker.GetHash([]byte(password))
	hlog.FromRequest(r).Info().Msg("bcrypt encrypt end")
	if err != nil {
		(*renderParams)["error"] = ErrInternalError
		server.shortnerTemplate.Execute(w, renderParams)
		return
	}
	hashedPassword := string(hashed)

	hlog.FromRequest(r).Info().Msg("secret box encrypt start")
	nonce, salt, encryptedLong, err := server.customcrypto.Encrypt(password, link)
	if err != nil {
		(*renderParams)["error"] = ErrInternalError
		server.shortnerTemplate.Execute(w, renderParams)
		return
	}
	hlog.FromRequest(r).Info().Msg("secret box encrypt end")

	urlObj := &url.ShortURL{Short: shortURL, Nonce: nonce, Salt: salt, EncryptedLong: encryptedLong, PasswordHash: hashedPassword}

	hlog.FromRequest(r).Info().Msg("saving to redis start")
	err = server.urlService.Create(*urlObj)
	if err != nil {
		hlog.FromRequest(r).Err(err).Msg("Shit broke")
		(*renderParams)["error"] = ErrInternalError
		server.shortnerTemplate.Execute(w, renderParams)
		return
	}
	hlog.FromRequest(r).Info().Msg("saving to redis end")

	(*renderParams)["shortURL"] = shortURL
	(*renderParams)["passwordProtect"] = passwordProtect
	(*renderParams)["password"] = password
	(*renderParams)["longURL"] = link

	server.createdTemplate.Execute(w, renderParams)
}

func (server *Server) elongateGet(w http.ResponseWriter, r *http.Request) {
	server.checkShortURLAndPassword(w, r, false)
}

func (server *Server) elongatePost(w http.ResponseWriter, r *http.Request) {
	server.checkShortURLAndPassword(w, r, true)
}

func (server *Server) urlCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID := chi.URLParam(r, "id")
		URLObject, err := server.urlService.GetLong(ID)
		if err != nil {
			if err.Error() == url.ErrorNotFound {
				http.Error(w, http.StatusText(404), 404)
				return
			}
			http.Error(w, http.StatusText(500), 500)
			return
		}
		ctx := context.WithValue(r.Context(), ctxKeyURLObject, URLObject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (server *Server) checkShortURLAndPassword(w http.ResponseWriter, r *http.Request, isPost bool) {
	renderParams := r.Context().Value(ctxKeyRenderParms).(*map[string]interface{})

	ctx := r.Context()

	URLObject, ok := ctx.Value(ctxKeyURLObject).(*url.ShortURL)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	password := "default"
	canPass := false
	hlog.FromRequest(r).Info().Msg("bcrypt decrypt start")
	err := server.passwordChecker.IsSame([]byte(URLObject.PasswordHash), []byte(password))
	hlog.FromRequest(r).Info().Msg("bcrypt decrypt end")
	if err != nil {
		if server.passwordChecker.NoMatch(err) {
			(*renderParams)["shortURL"] = URLObject.Short
			(*renderParams)["passwordProtect"] = true
			(*renderParams)["error"] = ErrPasswordMissing

			if isPost {
				password = r.FormValue("password")
				if password != "" {
					hlog.FromRequest(r).Info().Msg("bcrypt decrypt start")
					err = server.passwordChecker.IsSame([]byte(URLObject.PasswordHash), []byte(password))
					hlog.FromRequest(r).Info().Msg("bcrypt decrypt end")
					if err != nil {
						if server.passwordChecker.NoMatch(err) {
							(*renderParams)["error"] = ErrPasswordMatchFailed
						} else {
							fmt.Printf("_%v_\n", err.Error())
							(*renderParams)["error"] = ErrInternalError
						}
					} else {
						(*renderParams)["error"] = nil
						canPass = true
					}
				}
			}
		} else {
			fmt.Printf("_%v_\n", err.Error())
			(*renderParams)["error"] = ErrInternalError
		}
	}

	if !canPass {
		server.elongateTemplate.Execute(w, renderParams)
		return
	}

	longURL, err := server.customcrypto.Decrypt(password, URLObject.EncryptedLong, URLObject.Nonce, URLObject.Salt)
	if err != nil {
		(*renderParams)["error"] = ErrInternalError
		server.elongateTemplate.Execute(w, renderParams)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
