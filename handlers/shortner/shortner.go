package shortner

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"../../db"
	"../../db/url"
	. "../common"
)

var urls url.Service

func InitShortner() {
	urls = db.NewURLInterfaceMySQL()
	parseTemplates()
}

// ShortnerGet is userd asd
func Get(w http.ResponseWriter, r *http.Request) {
	shortnerTemplate.Execute(w, nil)
}

// ShortnerPost us
func Post(w http.ResponseWriter, r *http.Request) {
	errParseForm := r.ParseForm()

	if errParseForm != nil {
		fmt.Println("parse error", errParseForm)
		shortnerTemplate.Execute(w, map[string]interface{}{
			"error": ConstErrInternalError,
		})
		return
	}

	urlNext := "/"
	if len(r.Form["url_next"]) > 0 {
		urlNext = r.Form["url_next"][0]
	}

	url := r.FormValue("url")
	passwordProtect := r.FormValue("passwordProtect") == "on"
	password := r.FormValue("password")

	if url == "" || !govalidator.IsURL(url) {
		shortnerTemplate.Execute(w, map[string]interface{}{
			"error":    ConstErrURLMissing,
			"url_next": urlNext,
		})
		return
	}

	if passwordProtect && password == "" {
		shortnerTemplate.Execute(w, map[string]interface{}{
			"error":    ConstErrPasswordMissing,
			"url_next": urlNext,
		})
		return
	} else if !passwordProtect {
		password = "default"
	}

	var shortURL string
	var err error
	for true {
		shortURL, err = generateRandomString2(2)
		if err != nil {
			shortnerTemplate.Execute(w, map[string]interface{}{
				"error": ConstErrInternalError,
			})
			return
		}
		present, err := urls.Present(shortURL)
		if err != nil {
			shortnerTemplate.Execute(w, map[string]interface{}{
				"error": ConstErrInternalError,
			})
			return
		}
		if !present {
			break
		}
	}

	var hashedPassword string

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		shortnerTemplate.Execute(w, map[string]interface{}{
			"error": ConstErrInternalError,
		})
		return
	}
	hashedPassword = string(hashed)

	x, y, z, err := encryptNew(password, url)
	if err != nil {
		shortnerTemplate.Execute(w, map[string]interface{}{
			"error": ConstErrInternalError,
		})
		return
	}

	err = urls.Create(shortURL, x, y, z, hashedPassword)
	if err != nil {
		shortnerTemplate.Execute(w, map[string]interface{}{
			"error": ConstErrInternalError,
		})
		return
	}

	createdTemplate.Execute(w, map[string]interface{}{
		"shortURL":        shortURL,
		"passwordProtect": passwordProtect,
		"password":        password,
		"longURL":         url,
	})

}

// ElongateGet is
func ElongateGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["id"]

	longURL, renderParams := checkShorURLAndPassword(shortURL, r)
	if longURL == "" {
		elongateTemplate.Execute(w, renderParams)
		return
	}

	http.Redirect(w, r, longURL, http.StatusSeeOther)
}

// ElongatePost is
func ElongatePost(w http.ResponseWriter, r *http.Request) {
	errParseForm := r.ParseForm()

	if errParseForm != nil {
		fmt.Println("parse error", errParseForm)
		shortnerTemplate.Execute(w, map[string]interface{}{
			"error": ConstErrInternalError,
		})
		return
	}

	vars := mux.Vars(r)
	shortURL := vars["id"]

	longURL, renderParams := checkShorURLAndPassword(shortURL, r)
	if longURL == "" {
		elongateTemplate.Execute(w, renderParams)
		return
	}

	http.Redirect(w, r, longURL, http.StatusSeeOther)

}

func checkShorURLAndPassword(shortURL string, r *http.Request) (string, map[string]interface{}) {
	mp, err := urls.GetLong(shortURL)

	if err != nil {
		return "", map[string]interface{}{
			"error": ConstErrInternalError,
		}
	}

	if len(mp) == 0 {
		return "", map[string]interface{}{
			"error": ConstErrURLMissing,
		}
	}

	password := "default"
	err = bcrypt.CompareHashAndPassword([]byte(mp["passwordHash"]), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			password = r.FormValue("password")
			if password == "" {
				return "", map[string]interface{}{
					"shortURL":        shortURL,
					"error":           ConstErrPasswordMissing,
					"passwordProtect": true,
				}
			}
			err = bcrypt.CompareHashAndPassword([]byte(mp["passwordHash"]), []byte(password))
			if err != nil {
				if err == bcrypt.ErrMismatchedHashAndPassword {
					return "", map[string]interface{}{
						"shortURL":        shortURL,
						"error":           ConstErrPasswordMatchFailed,
						"passwordProtect": true,
					}
				}
			}
		} else {
			fmt.Printf("_%v_\n", err.Error())
			return "", map[string]interface{}{
				"error": ConstErrInternalError,
			}
		}
	}

	longURL, err := decryptNew(password, mp["encrypted"], mp["nonce"], mp["salt"])
	if err != nil {
		return "", map[string]interface{}{
			"error": ConstErrInternalError,
		}
	}

	return longURL, nil
}
