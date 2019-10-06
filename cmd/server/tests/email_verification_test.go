package tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/responses"
	"github.com/vovainside/vobook/config"
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/factories"
	. "github.com/vovainside/vobook/tests/apitest"
	"github.com/vovainside/vobook/tests/assert"
	"github.com/vovainside/vobook/utils"
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

func TestVerifyNEmail_ShouldBeExpired(t *testing.T) {
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
