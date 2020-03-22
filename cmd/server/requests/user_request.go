package requests

import (
	"vobook/database/models"
	"vobook/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
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

type ChangeEmail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *ChangeEmail) Validate() (err error) {
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

func (r *RegisterUser) ToModel() (m *models.User, err error) {
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
			validation.Required.Error("That's ridiculous fail: where's your email?"),
			// TODO change to "login" to allow login with email or phone
			//validation.Required.Error("Don't act like Mister Bean, just enter your login"),
			is.Email,
		),
		validation.Field(
			&r.Password,
			validation.Required.Error("That sounds like a prank cause your password is blank"),
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

type DeleteUser struct {
	Password string `json:"password"`
}

func (r *DeleteUser) Validate() (err error) {
	err = validation.ValidateStruct(r,
		validation.Field(
			&r.Password,
			validation.Required,
		),
	)

	return
}

type UpdateUser struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

func (r *UpdateUser) Validate() (err error) {
	return
}

func (r *UpdateUser) ToModel(m *models.User) {
	if r.LastName != nil {
		m.LastName = *r.LastName
	}
	if r.FirstName != nil {
		m.FirstName = *r.FirstName
	}
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
