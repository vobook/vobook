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

func (r *RegisterUser) Validate() (user *models.User, err error) {
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
	if err != nil {
		return
	}

	password, err := utils.HashPassword(r.Password)
	if err != nil {
		return
	}

	user = &models.User{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Password:  password,
	}
	return
}
