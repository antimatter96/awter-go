package handlers

import (
	"context"
	"crypto/rand"
	"fmt"
	"html/template"
	simplerand "math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
)

// NewLoginHandlerGet is userd asd
func ShortnerGet(ctx context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	shortnerTemplate = template.Must(template.ParseFiles("./template/shortner.html"))
	shortnerTemplate.Execute(w, nil)
}

// NewLoginHandlerPost us
func ShortnerPost(ctx context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	createdTemplate = template.Must(template.ParseFiles("./template/created.html"))
	errParseForm := r.ParseForm()

	if errParseForm != nil {
		fmt.Println("parse error", errParseForm)
		shortnerTemplate.Execute(w, map[string]interface{}{
			"error": constErrInternalError,
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
			"error":    constErrURLMissing,
			"url_next": urlNext,
		})
		return
	}

	if passwordProtect && password == "" {
		shortnerTemplate.Execute(w, map[string]interface{}{
			"error":    constErrPasswordMissing,
			"url_next": urlNext,
		})
		return
	}

	var shortURL string
	var err error
	for true {
		shortURL, err = generateRandomString2(6)
		if err != nil {
			shortnerTemplate.Execute(w, map[string]interface{}{
				"error": constErrInternalError,
			})
			return
		}
		present, err := urlService.PresentShort(shortURL)
		if err != nil {
			shortnerTemplate.Execute(w, map[string]interface{}{
				"error": constErrInternalError,
			})
			return
		}
		if !present {
			break
		}
	}

	var hashedPassword string
	if passwordProtect {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
		if err != nil {
			shortnerTemplate.Execute(w, map[string]interface{}{
				"error": constErrInternalError,
			})
			return
		}
		hashedPassword = string(hashed)
	}

	err = urlService.CreatePassword(shortURL, url, hashedPassword)
	if err != nil {
		shortnerTemplate.Execute(w, map[string]interface{}{
			"error": constErrInternalError,
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

func ElongateGet(ctx context.Context, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	elongateTemplate = template.Must(template.ParseFiles("./template/elongate.html"))
	shortURL := ps.ByName("id")

	present, longURL, password, err := urlService.GetLong(shortURL)
	if err != nil {
		elongateTemplate.Execute(w, map[string]interface{}{
			"error": constErrInternalError,
		})
		return
	}
	if !present {
		elongateTemplate.Execute(w, map[string]interface{}{
			"error": constErrURLMissing,
		})
		return
	}
	if password != "" {
		elongateTemplate.Execute(w, map[string]interface{}{
			"shortURL":        shortURL,
			"error":           constErrPasswordMissing,
			"passwordProtect": true,
		})
		return
	}
	http.Redirect(w, r, longURL, http.StatusSeeOther)
}

func ElongatePost(ctx context.Context, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	errParseForm := r.ParseForm()

	if errParseForm != nil {
		fmt.Println("parse error", errParseForm)
		shortnerTemplate.Execute(w, map[string]interface{}{
			"error": constErrInternalError,
		})
		return
	}

	elongateTemplate = template.Must(template.ParseFiles("./template/elongate.html"))
	shortURL := ps.ByName("id")

	fmt.Println(shortURL, len(shortURL))
	present, longURL, password, err := urlService.GetLong(shortURL)
	if err != nil {
		elongateTemplate.Execute(w, map[string]interface{}{
			"error": constErrInternalError,
		})
		return
	}
	if !present {
		elongateTemplate.Execute(w, map[string]interface{}{
			"error": constErrURLMissing,
		})
		return
	}
	if password == "" {
		http.Redirect(w, r, longURL, http.StatusSeeOther)
		return
	}
	user_password := r.FormValue("password")
	if user_password == "" {
		elongateTemplate.Execute(w, map[string]interface{}{
			"shortURL":        shortURL,
			"error":           constErrPasswordMissing,
			"passwordProtect": true,
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user_password))
	if err != nil {
		elongateTemplate.Execute(w, map[string]interface{}{
			"shortURL":        shortURL,
			"error":           constErrPasswordMatchFailed,
			"passwordProtect": true,
		})
		return
	}
	http.Redirect(w, r, longURL, http.StatusSeeOther)
}

func generateRandomString2(length int) (string, error) {

	randomFactor := make([]byte, 2)
	_, err := rand.Read(randomFactor)
	if err != nil {
		return "", err
	}

	simplerand.Seed(time.Now().UnixNano() * int64(randomFactor[0]) * int64(randomFactor[1]))

	arr := make([]byte, length)
	for i := range arr {
		arr[i] = letterRunes[simplerand.Intn(len(letterRunes))]
	}
	return string(arr), nil
}
