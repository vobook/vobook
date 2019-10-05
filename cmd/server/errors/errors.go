package errors

import (
	"net/http"

	"github.com/vovainside/vobook/config"
)

var (
	ReqisterUserEmailExists       = New422("User with this email already registered")
	EmailVerificationNotExists    = New422("Invalid verification code")
	EmailVerificationExpired      = New422("E-mail verification code expired")
	EmailChangeEmailAlreadyExists = New422("Email already in use by another user")
	WrongEmailOrPassword          = New422("Wrong email or password")
	InvalidAppClient              = New400("Invalid client")
	AuthTokenMissing              = New401("Token missing")
	AuthTokenInvalid              = New401("Invalid token")
	AuthTokenExpired              = New401("Token expired")
)

func New400(message string) error {
	return Error{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}
func New401(message string) error {
	return Error{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func New422(message string) error {
	return Error{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}

func New422Err(message string, err error) error {
	return Error{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
		Err:     err,
	}
}

type Error struct {
	Code    int
	Message string
	Err     error
}

func (e Error) Error() string {
	message := e.Message
	if !config.IsReleaseEnv() && e.Err != nil {
		message = e.Err.Error()
	}

	if message == "" {
		message = http.StatusText(e.Code)
	}

	return message
}
