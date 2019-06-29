package shortner

import (
	"fmt"
	"net/http"

	"github.com/antimatter96/awter-go/utils"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/antimatter96/awter-go/db"
	"github.com/antimatter96/awter-go/db/url"
	. "github.com/antimatter96/awter-go/handlers/common"
)

var urls url.Service

func InitShortner(store string) {
	switch store {
	case "mysql":
		urls = db.NewURLInterfaceMySQL()
	case "redis":
		urls = db.NewURLInterfaceRedis()
	}

	parseTemplates()
}

// Get renders the basic form
func Get(w http.ResponseWriter, r *http.Request) {
	renderParams, _ := r.Context().Value(CtxKeyRenderParms).(map[string]interface{})
	shortnerTemplate.Execute(w, renderParams)
}

// Post handles creation of a short URL
func Post(w http.ResponseWriter, r *http.Request) {
	renderParams, _ := r.Context().Value(CtxKeyRenderParms).(map[string]interface{})
	errParseForm := r.ParseForm()

	if errParseForm != nil {
		fmt.Println("parse error", errParseForm)
		renderParams["error"] = ConstErrInternalError
		shortnerTemplate.Execute(w, renderParams)
		return
	}

	urlNext := "/"
	if len(r.Form["url_next"]) > 0 {
		urlNext = r.Form["url_next"][0]
		renderParams["url_next"] = urlNext
	}

	link := r.FormValue("url")
	passwordProtect := r.FormValue("passwordProtect") == "on"
	password := r.FormValue("password")

	if link == "" || !govalidator.IsURL(link) {
		renderParams["error"] = ConstErrURLMissing
		shortnerTemplate.Execute(w, renderParams)
		return
	}

	if passwordProtect && password == "" {
		renderParams["error"] = ConstErrPasswordMissing
		renderParams["url_next"] = urlNext
		shortnerTemplate.Execute(w, renderParams)
		return
	} else if !passwordProtect {
		password = "default"
	}

	var shortURL string
	var err error
	for {
		shortURL, err = generateRandomString2(6)
		if err != nil {
			renderParams["error"] = ConstErrInternalError
			shortnerTemplate.Execute(w, renderParams)
			return
		}
		present, err := urls.Present(shortURL)
		if err != nil {
			renderParams["error"] = ConstErrInternalError
			shortnerTemplate.Execute(w, renderParams)
			return
		}
		if !present {
			break
		}
	}

	var hashedPassword string

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		renderParams["error"] = ConstErrInternalError
		shortnerTemplate.Execute(w, renderParams)
		return
	}
	hashedPassword = string(hashed)

	nonce, salt, encryptedLong, err := encryptNew(password, link)
	if err != nil {
		renderParams["error"] = ConstErrInternalError
		shortnerTemplate.Execute(w, renderParams)
		return
	}

	urlObj := &url.ShortURL{Short: shortURL, Nonce: nonce, Salt: salt, EncryptedLong: encryptedLong, PasswordHash: hashedPassword}

	err = urls.Create(*urlObj)
	if err != nil {
		renderParams["error"] = ConstErrInternalError
		shortnerTemplate.Execute(w, renderParams)
		return
	}

	renderParams["shortURL"] = shortURL
	renderParams["passwordProtect"] = passwordProtect
	renderParams["password"] = password
	renderParams["longURL"] = link

	createdTemplate.Execute(w, renderParams)

}

// ElongateGet is
func ElongateGet(w http.ResponseWriter, r *http.Request) {
	renderParams, _ := r.Context().Value(CtxKeyRenderParms).(map[string]interface{})

	vars := mux.Vars(r)
	shortURL := vars["id"]

	longURL, _renderParams := checkShorURLAndPassword(shortURL, r)
	if longURL == "" {
		renderParams = utils.MapMerge(renderParams, _renderParams)
		elongateTemplate.Execute(w, renderParams)
		return
	}

	http.Redirect(w, r, longURL, http.StatusSeeOther)
}

// ElongatePost is
func ElongatePost(w http.ResponseWriter, r *http.Request) {
	renderParams, _ := r.Context().Value(CtxKeyRenderParms).(map[string]interface{})
	errParseForm := r.ParseForm()

	if errParseForm != nil {
		fmt.Println("parse error", errParseForm)
		renderParams["error"] = ConstErrInternalError
		shortnerTemplate.Execute(w, renderParams)
		return
	}

	ElongateGet(w, r)
}

func checkShorURLAndPassword(shortURL string, r *http.Request) (longURL string, _renderParams map[string]interface{}) {
	_renderParams = make(map[string]interface{})
	urlObj, err := urls.GetLong(shortURL)

	if err != nil {
		fmt.Println(err)
		_renderParams["error"] = ConstErrInternalError
	}

	if urlObj.Short == "" {
		fmt.Println("Empty Map")
		_renderParams["error"] = ConstErrURLNotPresent
	}

	password := "default"
	err = bcrypt.CompareHashAndPassword([]byte(urlObj.PasswordHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			_renderParams["shortURL"] = shortURL
			_renderParams["passwordProtect"] = true

			password = r.FormValue("password")
			if password == "" {
				_renderParams["error"] = ConstErrPasswordMissing
				return
			}
			err = bcrypt.CompareHashAndPassword([]byte(urlObj.PasswordHash), []byte(password))
			if err != nil {
				if err == bcrypt.ErrMismatchedHashAndPassword {
					_renderParams["error"] = ConstErrPasswordMatchFailed
					return
				}
			}
		} else {
			fmt.Printf("_%v_\n", err.Error())
			_renderParams["error"] = ConstErrInternalError
			return
		}
	}

	longURL, err = decryptNew(password, urlObj.EncryptedLong, urlObj.Nonce, urlObj.Salt)
	if err != nil {
		_renderParams["error"] = ConstErrInternalError
	}

	return longURL, nil
}
