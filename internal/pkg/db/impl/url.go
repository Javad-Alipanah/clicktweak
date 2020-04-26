package impl

import (
	"fmt"

	exception "clicktweak/internal/pkg/error"
	"clicktweak/internal/pkg/model"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const (
	urlsTable = "urls"
	idCol     = "id"
	userIdCol = "user_id"
)

// Url implements db.Url interface
type Url struct {
	db *gorm.DB
}

func NewUrlDB(db *gorm.DB) (*Url, error) {
	db.AutoMigrate(&model.Url{})
	return &Url{db}, nil
}

func (u *Url) GetByID(id string) (*model.Url, error) {
	result := new(model.Url)
	err := u.db.Table(urlsTable).Where(fmt.Sprintf("%s = ?", idCol), id).First(result).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			log.Error(err)
			return nil, exception.InternalServerError
		}
		return nil, nil
	}
	return result, nil
}

func (u *Url) GetByUserID(id string) ([]*model.Url, error) {
	var result []*model.Url
	err := u.db.Table(urlsTable).Where(fmt.Sprintf("%s = ?", userIdCol), id).Find(result).Error
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			log.Error(err)
			return nil, exception.InternalServerError
		}
		return nil, nil
	}
	return result, nil
}

func (u *Url) Save(url *model.Url) error {
	// begin transaction
	tx := u.db.Begin()
	if err := tx.Error; err != nil {
		log.Error(err)
		return exception.InternalServerError
	}
	// check if url exists
	temp := new(model.Url)
	err := tx.Table(urlsTable).Where(fmt.Sprintf("%s = ?", idCol), url.ID).First(temp).Error
	if err == nil {
		tx.Rollback()
		return exception.ResourceAlreadyExists
	}
	if !gorm.IsRecordNotFoundError(err) {
		log.Errorln(err)
		tx.Rollback()
		return exception.InternalServerError
	}

	err = u.db.Create(url).Error
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
