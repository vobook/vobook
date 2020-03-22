package tests

import (
	"net/http"
	"testing"
	"time"

	"vobook/cmd/server/errors"
	"vobook/cmd/server/responses"
	"vobook/config"
	"vobook/database"
	"vobook/database/factories"
	. "vobook/tests/apitest"
	"vobook/tests/assert"
	"vobook/utils"
)

func TestVerifyExistingEmail(t *testing.T) {
	ev, err := factories.CreateEmailVerification()
	assert.NotError(t, err)

	ev.Email = ""
	err = database.ORM().Update(&ev)
	assert.NotError(t, err)

	var resp responses.Success
	GET(t, Request{
		Path:         "verify-email/" + ev.Token,
		AssertStatus: http.StatusOK,
		BindResponse: &resp,
		IsPublic:     true,
	})

	assert.DatabaseHas(t, "users", utils.M{
		"id":             ev.UserID,
		"email_verified": true,
	})
	assert.DatabaseMissing(t, "email_verifications", utils.M{
		"id": ev.ID,
	})
}

func TestVerifyNewEmail(t *testing.T) {
	ev, err := factories.CreateEmailVerification()
	assert.NotError(t, err)

	var resp responses.Success
	GET(t, Request{
		Path:         "verify-email/" + ev.Token,
		AssertStatus: http.StatusOK,
		BindResponse: &resp,
		IsPublic:     true,
	})

	assert.DatabaseHas(t, "users", utils.M{
		"id":             ev.UserID,
		"email":          ev.Email,
		"email_verified": true,
	})
	assert.DatabaseMissing(t, "email_verifications", utils.M{
		"id": ev.ID,
	})
}

func TestVerifyEmail_ShouldBeExpired(t *testing.T) {
	ev, err := factories.CreateEmailVerification()
	assert.NotError(t, err)

	ev.CreatedAt = time.Now().Add(-(config.Get().EmailVerificationLifetime + time.Minute))
	err = database.ORM().Update(&ev)
	assert.NotError(t, err)

	var resp responses.Error
	GET(t, Request{
		Path:         "verify-email/" + ev.Token,
		AssertStatus: http.StatusUnprocessableEntity,
		BindResponse: &resp,
		IsPublic:     true,
	})

	assert.Equals(t, errors.EmailVerificationExpired.Error(), resp.Error)
}
