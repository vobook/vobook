package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/utils"
)

type RegisterUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (r *RegisterUser) Validate() (err error) {
	err = validation.ValidateStruct(r,
		validation.Field(
			&r.Email,
			validation.Required,
			is.Email,
		),
		validation.Field(
			&r.Password,
			validation.Required,
		),
		validation.Field(
			&r.FirstName,
			validation.Required,
		),
		validation.Field(
			&r.LastName,
			validation.Required,
		),
	)

	return
}

func (r *RegisterUser) ToUser() (m *models.User, err error) {
	password, err := utils.HashPassword(r.Password)
	if err != nil {
		return
	}

	m = &models.User{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Password:  password,
	}

	return
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *Login) Validate() (err error) {
	err = validation.ValidateStruct(r,
		validation.Field(
			&r.Email,
			validation.Required,
			is.Email,
		),
		validation.Field(
			&r.Password,
			validation.Required,
		),
	)

	return
}

type ChangeUserPassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (r *ChangeUserPassword) Validate() (err error) {
	err = validation.ValidateStruct(r,
		validation.Field(
			&r.OldPassword,
			validation.Required,
		),
		validation.Field(
			&r.NewPassword,
			validation.Required,
		),
	)

	return
}

type ResetPasswordStart struct {
	Email string `json:"email"`
}

func (r *ResetPasswordStart) Validate() (err error) {
	err = validation.ValidateStruct(r,
		validation.Field(
			&r.Email,
			validation.Required,
			is.Email,
		),
	)

	return
}

type ResetPassword struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (r *ResetPassword) Validate() (err error) {
	err = validation.ValidateStruct(r,
		validation.Field(
			&r.Token,
			validation.Required,
		),
		validation.Field(
			&r.Password,
			validation.Required,
		),
	)

	return
}
