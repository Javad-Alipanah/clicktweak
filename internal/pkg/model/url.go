package model

import (
	"clicktweak/internal/pkg/util"
	"errors"
	"net/url"
	"regexp"
	"time"
)

type Url struct {
	ID         string    `gorm:"primary_key" json:"-"`
	Url        string    `gorm:"Column:url;NOT NULL" json:"url"`
	Suggestion string    `sql:"-" json:"suggestion"`
	CreatedAt  time.Time `json:"-"`
	UserID     uint      `gorm:"INDEX:user_id" json:"-"`
}

// Validate sanitizes given url and suggestion
func (u *Url) Validate() error {
	if _, err := url.ParseRequestURI(u.Url); err != nil {
		if _, err2 := url.ParseRequestURI("http://" + u.Url); err2 != nil {
			return err
		}
		u.Url = "http://" + u.Url
	}

	if !regexp.MustCompile(util.DefaultRegex).MatchString(u.Suggestion) {
		return errors.New("invalid suggestion")
	}
	return nil
}
