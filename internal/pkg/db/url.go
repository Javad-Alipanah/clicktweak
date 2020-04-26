package db

import "clicktweak/internal/pkg/model"

// Url is an abstraction for urls database
type Url interface {
	// GetByID retrieves url from database by ID
	//
	// returns (url, nil) on success and (nil, err) on failure
	GetByID(id string) (*model.Url, error)

	// GetByUserID retrieves urls corresponding to given user ID
	GetByUserID(userID uint) ([]*model.Url, error)

	// Save stores given url in database
	Save(url *model.Url) error
}
