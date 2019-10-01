package tests

import (
	"testing"

	fake "github.com/brianvoe/gofakeit"
	"github.com/vovainside/vobook/cmd/server/requests"
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
		AssertStatus: 200,
		BindResponse: &resp,
		IsPublic:     true,
	})

	assert.Equals(t, req.FirstName, resp.FirstName)
	assert.Equals(t, req.LastName, resp.LastName)
	assert.Equals(t, req.Email, resp.Email)
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

	var resp models.User
	apitest.POST(t, apitest.Request{
		Path:         "users/register",
		Body:         req,
		AssertStatus: 200,
		BindResponse: &resp,
		IsPublic:     true,
	})

	assert.Equals(t, req.FirstName, resp.FirstName)
	assert.Equals(t, req.LastName, resp.LastName)
	assert.Equals(t, req.Email, resp.Email)
}
