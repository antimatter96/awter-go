package handlers

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	simplerand "math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/argon2"

	"github.com/asaskevich/govalidator"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
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
	var salt string
	var toStoreURL string = url
	if passwordProtect {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
		if err != nil {
			shortnerTemplate.Execute(w, map[string]interface{}{
				"error": constErrInternalError,
			})
			return
		}
		hashedPassword = string(hashed)

		salt, err = getSalt(8)
		if err != nil {
			shortnerTemplate.Execute(w, map[string]interface{}{
				"error": constErrInternalError,
			})
			return
		}
		fmt.Println(salt)
		fmt.Println([]byte(salt))

		ek := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
		fmt.Println("KEY", ek)
		final, err := encrypt([]byte(url), ek)
		if err != nil {
			shortnerTemplate.Execute(w, map[string]interface{}{
				"error": constErrInternalError,
			})
			return
		}
		fmt.Println(final)
		fmt.Println(string(final))
		toStoreURL = string(final)
	}

	err = urlService.CreatePassword(shortURL, toStoreURL, hashedPassword, salt)
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

	present, longURL, password, _, err := urlService.GetLong(shortURL)
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

	present, longURL, password, storedSalt, err := urlService.GetLong(shortURL)
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
	if storedSalt != "" {
		salt, err := base64.URLEncoding.DecodeString(storedSalt)
		if err != nil {
			elongateTemplate.Execute(w, map[string]interface{}{
				"shortURL":        shortURL,
				"error":           constErrInternalError,
				"passwordProtect": true,
			})
			return
		}
		fmt.Println("D SALT", salt)
		fmt.Println("D SALT DTRIG", string(salt))
		dk := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
		fmt.Println("KEY", dk)
		longURLBytes, err := base64.URLEncoding.DecodeString(longURL)
		if err != nil {
			fmt.Println("Decoding of longURL", err)
			elongateTemplate.Execute(w, map[string]interface{}{
				"shortURL":        shortURL,
				"error":           constErrInternalError,
				"passwordProtect": true,
			})
			return
		}
		fmt.Println("D longurl bytes", longURLBytes)
		fmt.Println("D longurl", string(longURLBytes))
		//fmt.Println(longURLBytes)
		final, err := decrypt(longURLBytes, dk)
		if err != nil {
			fmt.Println("Decryption", err)
			elongateTemplate.Execute(w, map[string]interface{}{
				"shortURL":        shortURL,
				"error":           constErrInternalError,
				"passwordProtect": true,
			})
			return
		}
		longURL = string(final)
		//fmt.Println(string(final))
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

func sendInternalServerError(template *template.Template, w http.ResponseWriter) {
	template.Execute(w, map[string]interface{}{
		"error": constErrInternalError,
	})
	return
}

func encrypt(plaintext []byte, key []byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}
