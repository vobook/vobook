package errors

import (
	"fmt"
	"net/http"
	"strings"

	"vobook/config"
)

var (
	ReqisterUserEmailExists    = New422("User with this email already registered")
	EmailVerificationNotExists = New422("Invalid verification code")
	EmailVerificationExpired   = New422("E-mail verification code expired")
	EmailChangeEmailInUser     = New422("Email already in use by another user")
	EmailChangeSameEmail       = New422("You already using this email address")
	WrongEmailOrPassword       = New422("Wrong email or password")
	InvalidAppClient           = New400("Invalid client")
	AuthTokenMissing           = New401("Token missing")
	AuthTokenInvalidLength     = New401("Invalid token (001)")
	AuthTokenInvalidSign       = New401("Invalid token (002)")
	AuthTokenNotFound          = New401("Invalid token (003)")
	AuthTokenExpired           = New401("Token expired")
	WrongPassword              = New422("Wrong password")
	UserByEmailNotFound        = New404("User with this email not exists")
	PasswordResetTokenNotFound = New404("Invalid password reset request")
	PasswordResetTokenExpired  = New404("Password reset request expired")
	InvalidContactPropertyType = New422("Invalid contact property type")
	InvalidGender              = New422("Invalid gender")
	ContactNotFound            = New404("Contact not found")
	CreateContactNameEmpty     = New422("Contact name missing (Enter name or first name or last name)")
	ContactPhotoNotExists      = New404("Contact do not have photo")
)

func New400(message string) error {
	return Error{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}
func New401(message string) error {
	return Error{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func New422(message string) error {
	return Error{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}

func New404(message string) error {
	return Error{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

type Error struct {
	Code    int
	Message string
	Err     error
}

type List []string
type Input map[string]string

func (l List) Error() string {
	s := strings.Builder{}
	for i, v := range l {
		s.WriteString(fmt.Sprintf("%d. %s\n", i+1, v))
	}
	return s.String()
}

func (i Input) Add(k, v string) {
	i[k] = v
}
func (i Input) Has() bool {
	return len(i) > 0
}

func (i Input) Error() string {
	s := strings.Builder{}
	for k, v := range i {
		s.WriteString(k + ": " + v + "\n")
	}
	return s.String()
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
