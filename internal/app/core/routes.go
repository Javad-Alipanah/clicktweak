package core

import (
	"errors"
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

func Login(user db.User, secret string) echo.HandlerFunc {
	return func(context echo.Context) (err error) {
		u := new(model.User)
		if err = context.Bind(u); err != nil {
			log.Error(err)
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(exception.InternalServerError))
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
		password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(exception.InternalServerError))
		}
		if string(password) != result.Password {
			return context.JSON(http.StatusUnauthorized, exception.ToJSON(exception.InvalidCredentials))
		}

		// generate JWT
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = result.ID
		claims["exp"] = time.Now().Add(time.Hour * 24)

		// generate encoded token and send to client
		et, err := token.SignedString([]byte(secret))
		if err != nil {
			log.Errorln(err)
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(exception.InternalServerError))
		}

		return context.JSON(http.StatusOK, map[string]string{
			"token": et,
		})
	}
}
