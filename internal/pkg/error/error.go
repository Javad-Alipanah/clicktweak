package error

import "errors"

var (
	InternalServerError = errors.New("internal server error")
	UserAlreadyExists   = errors.New("user already exists")
	InvalidCredentials  = errors.New("invalid credentials")
)

func ToJSON(err error) map[string]string {
	return map[string]string{
		"message": err.Error(),
	}
}
