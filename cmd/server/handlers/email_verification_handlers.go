package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/requests"
	emailverification "github.com/vovainside/vobook/domain/email_verification"
	"github.com/vovainside/vobook/domain/user"
)

func VerifyEmail(c *gin.Context) {
	var req requests.RegisterUser

	u, err := req.Validate()
	if err != nil {
		abort422(c, err)
		return
	}

	_, err = user.FindByEmail(req.Email)
	if err == nil {
		abort(c, errors.ReqisterUserEmailExists)
		return
	}

	err = user.Create(u)
	if err != nil {
		abort(c, err)
		return
	}

	err = emailverification.Create(u.ID, u.Email)
	if err != nil {
		abort(c, err)
		return
	}

	// TODO send email verification email

	c.JSON(http.StatusOK, u)
}
