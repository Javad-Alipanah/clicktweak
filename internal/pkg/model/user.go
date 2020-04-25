package model

import (
	"errors"
	"regexp"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserName string `gorm:"INDEX:user_name;Column:user_name;NOT NULL;UNIQUE" json:"user_name,omitempty"`
	Email    string `gorm:"INDEX:email;Column:email;NOT NULL;UNIQUE" json:"email,omitempty"`
	Password string `gorm:"Column:password;NOT NULL" json:"password"`
}

func (u *User) Validate() error {
	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return err
	}

	validUserName := regexp.MustCompile(`^[a-zA-Z0-9_]{4,32}$`)
	validPassword := regexp.MustCompile(`^.{8,}$`)

	if !validUserName.MatchString(u.UserName) {
		return errors.New("username must be alphanumeric between 4 and 32 characters")
	}

	if !validPassword.MatchString(u.Password) {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}
