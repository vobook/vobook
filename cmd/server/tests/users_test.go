package tests

import (
	"net/http"
	"testing"

	fake "github.com/brianvoe/gofakeit"
	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/cmd/server/responses"
	"github.com/vovainside/vobook/database/factories"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/tests/apitest"
	"github.com/vovainside/vobook/tests/assert"
	"github.com/vovainside/vobook/utils"
)

func TestRegisterUser(t *testing.T) {
	email, err := utils.UniqueEmail("users")
	assert.NotError(t, err)

	req := requests.RegisterUser{
		FirstName: fake.FirstName(),
		LastName:  fake.LastName(),
		Email:     email,
		Password:  utils.RandomString(),
	}

	var resp models.User
	apitest.POST(t, apitest.Request{
		Path:         "users/register",
		Body:         req,
		AssertStatus: http.StatusOK,
		BindResponse: &resp,
		IsPublic:     true,
	})

	assert.Equals(t, req.FirstName, resp.FirstName)
	assert.Equals(t, req.LastName, resp.LastName)
	assert.Equals(t, req.Email, resp.Email)

	assert.DatabaseHas(t, "users", utils.M{
		"id":             resp.ID,
		"first_name":     req.FirstName,
		"last_name":      req.LastName,
		"email":          req.Email,
		"email_verified": false,
	})

	assert.DatabaseHas(t, "email_verifications", utils.M{
		"user_id": resp.ID,
		"email":   req.Email,
	})
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	user, err := factories.CreateUser()
	assert.NotError(t, err)

	req := requests.RegisterUser{
		FirstName: fake.FirstName(),
		LastName:  fake.LastName(),
		Email:     user.Email,
		Password:  utils.RandomString(),
	}

	var resp responses.Error
	apitest.POST(t, apitest.Request{
		Path:         "users/register",
		Body:         req,
		AssertStatus: http.StatusUnprocessableEntity,
		BindResponse: &resp,
		IsPublic:     true,
	})

	assert.Equals(t, resp.Error, errors.ReqisterUserEmailExists.Error())
}
