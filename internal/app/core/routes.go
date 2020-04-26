package core

import (
	"clicktweak/internal/pkg/util"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"clicktweak/internal/pkg/db"
	exception "clicktweak/internal/pkg/error"
	"clicktweak/internal/pkg/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Login receives user credentials and if correct responds with fresh jwt token
func Login(user db.User, secret string) echo.HandlerFunc {
	return func(context echo.Context) (err error) {
		u := new(model.User)
		if err = context.Bind(u); err != nil {
			return context.JSON(http.StatusBadRequest, exception.ToJSON(exception.MalformedRequest))
		}

		// only provide one of email or user_name fields
		if len(u.Email) > 0 && len(u.UserName) > 0 {
			err = errors.New("provide only one of fields user_name or email")
			return context.JSON(http.StatusBadRequest, exception.ToJSON(err))
		}

		// provide one of specified fields
		if len(u.Email) == 0 && len(u.UserName) == 0 {
			err = errors.New("you must provide email or user_name")
			return context.JSON(http.StatusBadRequest, exception.ToJSON(err))
		}

		// retrieve from db
		var result = new(model.User)
		if len(u.Email) > 0 {
			result, err = user.GetByEmail(u.Email)
		} else {
			result, err = user.GetByUserName(u.UserName)
		}

		if err == exception.InternalServerError {
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(err))
		}

		if result == nil {
			return context.JSON(http.StatusUnauthorized, exception.ToJSON(exception.InvalidCredentials))
		}

		// check password
		err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(u.Password))
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				return context.JSON(http.StatusUnauthorized, exception.ToJSON(exception.InvalidCredentials))
			}
			log.Error(err)
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(exception.InternalServerError))
		}

		et, err := generateToken(result, secret)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(err))
		}

		return context.JSON(http.StatusOK, map[string]string{
			"token": et,
		})
	}
}

// SingUp receives user info, validates and if correct adds user to database
//
// if user already exists, returns http status 409 (conflict)
// on success returns a fresh token to user
func SignUp(user db.User, secret string) echo.HandlerFunc {
	return func(context echo.Context) (err error) {
		u := new(model.User)
		if err = context.Bind(u); err != nil {
			log.Error(err)
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(exception.InternalServerError))
		}

		if err = u.Validate(); err != nil {
			return context.JSON(http.StatusBadRequest, exception.ToJSON(err))
		}

		// generate password hash
		password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(exception.InternalServerError))
		}
		u.Password = string(password)

		if err = user.Save(u); err != nil {
			var status int
			if err == exception.InternalServerError {
				status = http.StatusInternalServerError
			} else if err == exception.ResourceAlreadyExists {
				status = http.StatusConflict
			}
			return context.JSON(status, exception.ToJSON(err))
		}

		et, err := generateToken(u, secret)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(err))
		}

		return context.JSON(http.StatusCreated, map[string]string{
			"token": et,
		})
	}
}

// Shorten returns the shortened representation of given url
func Shorten(url db.Url) echo.HandlerFunc {
	return func(context echo.Context) error {
		u := new(model.Url)
		if err := context.Bind(u); err != nil {
			return context.JSON(http.StatusBadRequest, exception.ToJSON(exception.MalformedRequest))
		}

		user := context.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		userID := claims["id"].(float64)
		u.UserID = uint(userID)

		if err := u.Validate(); err != nil {
			return context.JSON(http.StatusBadRequest, exception.ToJSON(err))
		}

		var unique string
		encoder := util.NewEncoder62()
		if len(u.Suggestion) > 0 {
			unique = u.Suggestion
		} else {
			unique = encoder.Encode(rand.Uint32())
		}

		// do while
		var err error
		for {
			u.ID = unique
			err = url.Save(u)
			if err == nil || err != exception.ResourceAlreadyExists {
				break
			}
			// do generate unique
			unique, _ = encoder.SimilarSuggestion(unique)
		}

		if err != nil {
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(err))
		}
		return context.JSON(http.StatusCreated, map[string]string{
			"shortURL": context.Request().Host + "/" + u.ID,
		})
	}
}

func generateToken(u *model.User, secret string) (string, error) {
	// generate JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 24)

	// generate encoded token and send to client
	return token.SignedString([]byte(secret))
}
