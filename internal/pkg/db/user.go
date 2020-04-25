package db

import "clicktweak/internal/pkg/model"

// User is an abstraction for users database
type User interface {
	// GetByEmail retrieves user by email address
	//
	// returns (user, nil) on success, (nil, err) on failure, and (nil, nil) when not found
	GetByEmail(email string) (*model.User, error)

	// GetByUserName retrieves user by user name
	//
	// returns (user, nil) on success, (nil, err) on failure, and (nil, nil) when not found
	GetByUserName(userName string) (*model.User, error)

	// Saves given user in database
	//
	// returns nil on success and err on failure
	Save(user *model.User) error
}
