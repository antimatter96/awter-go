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
	createdTemplate = template.Must(template.ParseFiles("./template/created.html"))
	elongateTemplate = template.Must(template.ParseFiles("./template/elongate.html"))
	r.ParseForm()
	var urlNext string
	if len(r.Form["url_next"]) > 0 {
		urlNext = r.Form["url_next"][0]
	}
	shortnerTemplate.Execute(w, map[string]interface{}{
		"url_next": urlNext,
	})
}

// NewLoginHandlerPost us
func ShortnerPost(ctx context.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

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
	password_protect := r.FormValue("password_protect") == "on"
	password := r.FormValue("password")

	if url == "" || !govalidator.IsURL(url) {
		shortnerTemplate.Execute(w, map[string]interface{}{
			"error":    constErrURLMissing,
			"url_next": urlNext,
		})
		return
	}

	if password_protect && password == "" {
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
	if password_protect {
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
		"shortURL":         shortURL,
		"password_protect": password_protect,
		"password":         password,
		"longURL":          url,
	})

	// present, userID, passwordHash, errGetPassword := userService.GetPasswordHash("email")

	// if errGetPassword != nil {
	// 	loginTemplate.Execute(w, map[string]interface{}{
	// 		"error":    constErrInternalError,
	// 		"url_next": urlNext,
	// 	})
	// 	return
	// }

	// if !present {
	// 	loginTemplate.Execute(w, map[string]interface{}{
	// 		"error":    constErrNotRegistered,
	// 		"url_next": urlNext,
	// 	})
	// 	return
	// }

	// verified := checkPassword(&password, &passwordHash)

	// if !verified {
	// 	fmt.Println("Wrong passwword")
	// 	loginTemplate.Execute(w, map[string]interface{}{
	// 		"error":    constErrPasswordMatchFailed,
	// 		"url_next": urlNext,
	// 	})
	// 	//loginTemplate.Execute(w, )
	// 	return
	// }

	// fmt.Println(userID)

	// newSessionID, errGen := generateRandomString(48)
	// if errGen != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// plainValue := map[string]string{"sessionid": newSessionID}
	// encodedValue, errCookieEncode := cookie.Encode("sessionid", plainValue)

	// if errCookieEncode != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// } else {
	// 	newCookie := &http.Cookie{
	// 		Name:     "sessionid",
	// 		Value:    encodedValue,
	// 		Path:     "/",
	// 		HttpOnly: true,
	// 	}
	// 	http.SetCookie(w, newCookie)

	// 	errSetInCache := cache.SetSessionValue(newSessionID, "userId", userID)
	// 	if errSetInCache != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}
	// 	http.Redirect(w, r, urlNext, http.StatusSeeOther)
	// }

	// loginTemplate.Execute(w, nil)
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
			"shortURL": shortURL,
			"error":    constErrPasswordMissing,
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
