package model

import (
	"clicktweak/internal/pkg/util"
	"errors"
	"regexp"
	"strings"
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
	if !strings.HasPrefix(u.Url, "http://") && !strings.HasPrefix(u.Url, "https://") {
		u.Url = "http://" + u.Url
	}

	validURL := regexp.MustCompile(`^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`)
	if !validURL.MatchString(u.Url) {
		return errors.New("invalid url")
	}

	if !regexp.MustCompile(util.DefaultRegex).MatchString(u.Suggestion) {
		return errors.New("invalid suggestion")
	}
	return nil
}
