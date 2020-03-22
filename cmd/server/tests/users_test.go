package tests

import (
	"net/http"
	"testing"

	fake "github.com/brianvoe/gofakeit"

	"vobook/cmd/server/errors"
	"vobook/cmd/server/requests"
	"vobook/cmd/server/responses"
	"vobook/config"
	"vobook/database"
	"vobook/database/factories"
	"vobook/database/models"
	"vobook/domain/user"
	"vobook/services/mail"
	. "vobook/tests/apitest"
	"vobook/tests/assert"
	"vobook/utils"
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
	POST(t, Request{
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

func TestUpdateUser(t *testing.T) {
	ReLogin(t)
	firstName := fake.FirstName()
	lastName := fake.LastName()
	req := requests.UpdateUser{
		FirstName: &firstName,
		LastName:  &lastName,
	}

	var resp responses.Success
	PUT(t, Request{
		Path:         "user",
		Body:         req,
		AssertStatus: http.StatusOK,
		BindResponse: &resp,
	})

	assert.DatabaseHas(t, "users", utils.M{
		"id":         User(t).ID,
		"first_name": firstName,
		"last_name":  lastName,
	})
}

func TestChangeEmail(t *testing.T) {
	Login(t)
	email, err := utils.UniqueEmail("users")
	assert.NotError(t, err)

	elem := User(t)
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
	POST(t, Request{
		Path:         "change-email",
		Body:         req,
		AssertStatus: http.StatusOK,
		BindResponse: &resp,
	})

	assert.DatabaseHas(t, "email_verifications", utils.M{
		"user_id": AuthUser.ID,
		"email":   req.Email,
	})

	sentMail := mail.TestRepo.GetMail(req.Email)
	assert.Equals(t, config.Get().Mail.From, sentMail.From)
	assert.Equals(t, []string{req.Email}, sentMail.To)
	assert.True(t, sentMail.Subject != "")
	assert.True(t, sentMail.Body != "")
}

func TestDeleteAccount(t *testing.T) {
	ReLogin(t)
	elem := User(t)
	password := elem.Password
	passwordHash, err := utils.HashPassword(password)
	assert.NotError(t, err)

	err = user.UpdatePassword(elem.ID, passwordHash)
	assert.NotError(t, err)

	req := requests.DeleteUser{
		Password: password,
	}

	var resp responses.Success
	PUT(t, Request{
		Path:         "user/delete",
		Body:         req,
		AssertStatus: http.StatusOK,
		BindResponse: &resp,
	})

	assert.DatabaseMissing(t, "auth_tokens", utils.M{
		"user_id": AuthUser.ID,
	})

	assert.DatabaseHasDeleted(t, "users", AuthUser.ID)

	sentMail := mail.TestRepo.GetMail(AuthUser.Email)
	assert.Equals(t, config.Get().Mail.From, sentMail.From)
	assert.Equals(t, []string{AuthUser.Email}, sentMail.To)
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
	POST(t, Request{
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
	POST(t, Request{
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
	ReLogin(t)
	var resp models.User
	GET(t, Request{
		Path:         "user",
		AssertStatus: http.StatusOK,
		BindResponse: &resp,
	})

	assert.Equals(t, AuthUser.ID, resp.ID)
	assert.Equals(t, AuthUser.FirstName, resp.FirstName)
	assert.Equals(t, AuthUser.LastName, resp.LastName)
	assert.Equals(t, "", resp.Password)
}

func TestChangeUserPassword(t *testing.T) {
	Login(t)

	elem := User(t)
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
	POST(t, Request{
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
	POST(t, Request{
		Path:         "login",
		Body:         req2,
		AssertStatus: http.StatusOK,
		BindResponse: &resp2,
		IsPublic:     true,
	})

	assert.Equals(t, AuthUser.ID, elem.ID)
}

func TestPasswordReset(t *testing.T) {
	userEl, err := factories.CreateUser()
	assert.NotError(t, err)

	req := requests.ResetPasswordStart{
		Email: userEl.Email,
	}

	// password reset start
	var resp responses.Success
	POST(t, Request{
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
	POST(t, Request{
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
	PUT(t, Request{
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
	POST(t, Request{
		Path:         "login",
		Body:         req4,
		AssertStatus: http.StatusOK,
		BindResponse: &resp4,
		IsPublic:     true,
	})

	assert.Equals(t, userEl.ID, resp4.User.ID)
}
