package impl

import (
	"fmt"

	exception "clicktweak/internal/pkg/error"
	"clicktweak/internal/pkg/model"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const (
	usersTable  = "users"
	emailCol    = "email"
	userNameCol = "user_name"
	passwordCol = "password"
)

// User implements db.User interface
type User struct {
	db *gorm.DB
}

func NewUserDB(db *gorm.DB) (*User, error) {
	db.AutoMigrate(&model.User{})
	return &User{db}, nil
}

func (u *User) GetByEmail(email string) (*model.User, error) {
	result := new(model.User)
	err := u.db.Table(usersTable).Where(fmt.Sprintf("%s = ?", emailCol), email).First(result).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			log.Error(err)
			return nil, exception.InternalServerError
		}
		return nil, nil
	}
	return result, nil
}

func (u *User) GetByUserName(userName string) (*model.User, error) {
	result := new(model.User)
	err := u.db.Table(usersTable).Where(fmt.Sprintf("%s = ?", userNameCol), userName).First(result).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			log.Error(err)
			return nil, exception.InternalServerError
		}
		return nil, nil
	}
	return result, nil
}

func (u *User) Save(user *model.User) error {
	// begin transaction
	tx := u.db.Begin()
	if err := tx.Error; err != nil {
		log.Error(err)
		return exception.InternalServerError
	}
	// check if user exists
	temp := new(model.User)
	err := tx.Table(usersTable).Where(fmt.Sprintf("%s = ? OR %s = ?", userNameCol, emailCol), user.UserName, user.Email).First(temp).Error
	if err == nil {
		tx.Rollback()
		return exception.UserAlreadyExists
	}
	if !gorm.IsRecordNotFoundError(err) {
		log.Errorln(err)
		tx.Rollback()
		return exception.InternalServerError
	}

	err = u.db.Create(user).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return exception.InternalServerError
	}

	err = tx.Commit().Error
	if err != nil {
		log.Error(err)
		return exception.InternalServerError
	}
	return nil
}
