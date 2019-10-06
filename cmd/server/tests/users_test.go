package tests

import (
	"net/http"
	"testing"

	fake "github.com/brianvoe/gofakeit"
	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/cmd/server/responses"
	"github.com/vovainside/vobook/config"
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/factories"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/domain/user"
	"github.com/vovainside/vobook/services/mail"
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
		Path:         "register-user",
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

	sentMail := mail.TestRepo.GetMail(req.Email)
	assert.Equals(t, config.Get().Mail.From, sentMail.From)
	assert.Equals(t, []string{req.Email}, sentMail.To)
	assert.True(t, sentMail.Subject != "")
	assert.True(t, sentMail.Body != "")
}

func TestChangeEmail(t *testing.T) {
	apitest.Login(t)
	email, err := utils.UniqueEmail("users")
	assert.NotError(t, err)

	elem := apitest.User(t)
	password := elem.Password
	passwordHash, err := utils.HashPassword(password)
	assert.NotError(t, err)

	err = user.UpdatePassword(elem.ID, passwordHash)
	assert.NotError(t, err)

	req := requests.ChangeEmail{
		Email:    email,
		Password: password,
	}

	var resp responses.Success
	apitest.POST(t, apitest.Request{
		Path:         "change-email",
		Body:         req,
		AssertStatus: http.StatusOK,
		BindResponse: &resp,
	})

	assert.DatabaseHas(t, "email_verifications", utils.M{
		"user_id": apitest.AuthUser.ID,
		"email":   req.Email,
	})

	sentMail := mail.TestRepo.GetMail(req.Email)
	assert.Equals(t, config.Get().Mail.From, sentMail.From)
	assert.Equals(t, []string{req.Email}, sentMail.To)
	assert.True(t, sentMail.Subject != "")
	assert.True(t, sentMail.Body != "")
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	userEl, err := factories.CreateUser()
	assert.NotError(t, err)

	req := requests.RegisterUser{
		FirstName: fake.FirstName(),
		LastName:  fake.LastName(),
		Email:     userEl.Email,
		Password:  utils.RandomString(),
	}

	var resp responses.Error
	apitest.POST(t, apitest.Request{
		Path:         "register-user",
		Body:         req,
		AssertStatus: http.StatusUnprocessableEntity,
		BindResponse: &resp,
		IsPublic:     true,
	})

	assert.Equals(t, resp.Error, errors.ReqisterUserEmailExists.Error())
}

func TestUserLogin(t *testing.T) {
	userEl, err := factories.MakeUser()
	assert.NotError(t, err)

	password := utils.RandomString(10)
	passwordHash, err := utils.HashPassword(password)
	assert.NotError(t, err)

	userEl.Password = passwordHash
	err = database.ORM().Insert(&userEl)
	assert.NotError(t, err)

	req := requests.Login{
		Email:    userEl.Email,
		Password: password,
	}

	var resp responses.Login
	apitest.POST(t, apitest.Request{
		Path:         "login",
		Body:         req,
		AssertStatus: http.StatusOK,
		BindResponse: &resp,
		IsPublic:     true,
	})

	assert.DatabaseHas(t, "auth_tokens", utils.M{
		"user_id": userEl.ID,
	})
}

func TestViewCurrentUser(t *testing.T) {
	var resp models.User
	apitest.GET(t, apitest.Request{
		Path:         "user",
		AssertStatus: http.StatusOK,
		BindResponse: &resp,
	})

	assert.Equals(t, apitest.AuthUser.ID, resp.ID)
	assert.Equals(t, apitest.AuthUser.FirstName, resp.FirstName)
	assert.Equals(t, apitest.AuthUser.LastName, resp.LastName)
	assert.Equals(t, "", resp.Password)
}

func TestChangeUserPassword(t *testing.T) {
	apitest.Login(t)

	elem := apitest.User(t)
	oldPassword := elem.Password
	oldPasswordHash, err := utils.HashPassword(oldPassword)
	assert.NotError(t, err)

	err = user.UpdatePassword(elem.ID, oldPasswordHash)
	assert.NotError(t, err)

	newPassword := utils.RandomString(10)
	req := requests.ChangeUserPassword{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	var resp responses.Success
	apitest.POST(t, apitest.Request{
		Path:         "change-password",
		Body:         req,
		AssertStatus: http.StatusOK,
		BindResponse: &resp,
	})

	// login with new password
	req2 := requests.Login{
		Email:    elem.Email,
		Password: newPassword,
	}

	var resp2 responses.Login
	apitest.POST(t, apitest.Request{
		Path:         "login",
		Body:         req2,
		AssertStatus: http.StatusOK,
		BindResponse: &resp2,
		IsPublic:     true,
	})

	assert.Equals(t, apitest.AuthUser.ID, elem.ID)
}

func TestPasswordReset(t *testing.T) {
	userEl, err := factories.CreateUser()
	assert.NotError(t, err)

	req := requests.ResetPasswordStart{
		Email: userEl.Email,
	}

	// password reset start
	var resp responses.Success
	apitest.POST(t, apitest.Request{
		Path:         "reset-password",
		Body:         req,
		AssertStatus: http.StatusOK,
		BindResponse: &resp,
	})

	token := &models.PasswordReset{}
	err = database.ORM().Model(token).Where("user_id = ?", userEl.ID).First()
	assert.NotError(t, err)

	// check token
	var resp2 string
	apitest.POST(t, apitest.Request{
		Path:         "reset-password/" + token.Token,
		AssertStatus: http.StatusOK,
		BindResponse: &resp2,
	})

	// change password
	req3 := requests.ResetPassword{
		Token:    token.Token,
		Password: utils.RandomString(10),
	}
	var resp3 responses.Success
	apitest.PUT(t, apitest.Request{
		Path:         "reset-password",
		Body:         req3,
		AssertStatus: http.StatusOK,
		BindResponse: &resp3,
	})

	// login with new password
	req4 := requests.Login{
		Email:    userEl.Email,
		Password: req3.Password,
	}
	var resp4 responses.Login
	apitest.POST(t, apitest.Request{
		Path:         "login",
		Body:         req4,
		AssertStatus: http.StatusOK,
		BindResponse: &resp4,
		IsPublic:     true,
	})

	assert.Equals(t, userEl.ID, resp4.User.ID)
}
