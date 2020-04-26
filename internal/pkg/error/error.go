package error

import "errors"

var (
	InternalServerError   = errors.New("internal server error")
	ResourceAlreadyExists = errors.New("resource already exists")
	InvalidCredentials    = errors.New("invalid credentials")
	MalformedRequest      = errors.New("malformed input")
	ResourceNotFound      = errors.New("resource not found")
)

func ToJSON(err error) map[string]string {
	return map[string]string{
		"message": err.Error(),
	}
}
