package core

import (
	exception "clicktweak/internal/pkg/error"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"clicktweak/internal/pkg/model"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	user = model.User{
		UserName: "javad",
		Password: "gotcha!!",
		Email:    "javadalipanah@gmail.com",
	}
	url = model.Url{
		Url:        "http://alipanah.me/resume",
		Suggestion: "javad",
	}
)

// user database mock
type mockUserDB struct {
	users []*model.User
}

func (db *mockUserDB) GetByEmail(email string) (*model.User, error) {
	for _, user := range db.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("not found")
}

func (db *mockUserDB) GetByUserName(userName string) (*model.User, error) {
	for _, user := range db.users {
		if user.UserName == userName {
			return user, nil
		}
	}
	return nil, errors.New("not found")
}

func (db *mockUserDB) Save(user *model.User) error {
	for _, u := range db.users {
		if u.Email == user.Email || u.UserName == user.UserName {
			return exception.ResourceAlreadyExists
		}
	}
	var temp = new(model.User)
	copier.Copy(temp, user)
	pass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	temp.Password = string(pass)
	db.users = append(db.users, temp)
	return nil
}

// url database mock
type mockUrlDB struct {
	urls []*model.Url
}

func (db *mockUrlDB) GetByID(id string) (*model.Url, error) {
	for _, url := range db.urls {
		if url.ID == id {
			return url, nil
		}
	}
	return nil, errors.New("not found")
}

func (db *mockUrlDB) GetByUserID(userID uint) ([]*model.Url, error) {
	var result = make([]*model.Url, 0)
	for _, url := range db.urls {
		if url.UserID == userID {
			result = append(result, &model.Url{})
			if err := copier.Copy(result[len(result)-1], url); err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}

func (db *mockUrlDB) Save(url *model.Url) error {
	db.urls = append(db.urls, url)
	return nil
}

func TestSignUp(t *testing.T) {
	e := echo.New()
	userJSON, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := SignUp(&mockUserDB{}, "test")

	// success
	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var result map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &result)
		assert.Contains(t, result, "token")
	}

	// failure because user exists
	req = httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusConflict, rec.Code)
	}
}

func TestLogin(t *testing.T) {
	db := &mockUserDB{}
	_ = db.Save(&user)
	e := echo.New()
	user.UserName = ""
	userJSON, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := Login(db, "test")

	// success
	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var result map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &result)
		assert.Contains(t, result, "token")
	}

	// failure, invalid credentials
	user.Password = "invalid"
	userJSON, _ = json.Marshal(user)
	req = httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}
}

func TestShorten(t *testing.T) {
	db := &mockUrlDB{}
	e := echo.New()
	urlJSON, _ := json.Marshal(url)
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(string(urlJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := Shorten(db)

	// generate token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = float64(user.ID)
	claims["exp"] = time.Now().Add(time.Hour * 24)
	c.Set("user", token)

	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var result map[string]interface{}
		_ = json.Unmarshal(rec.Body.Bytes(), &result)
		assert.Contains(t, result, "shortURL")
	}
}
