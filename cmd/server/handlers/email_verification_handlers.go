package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/responses"
	"github.com/vovainside/vobook/config"
	emailverification "github.com/vovainside/vobook/domain/email_verification"
	"github.com/vovainside/vobook/domain/user"
)

func VerifyEmail(c *gin.Context) {
	token := c.Param("token")

	model, err := emailverification.FindByToken(token)
	if err != nil {
		Abort(c, err)
		return
	}

	if model.CreatedAt.Add(config.Get().EmailVerificationLifetime).Before(time.Now()) {
		Abort(c, errors.EmailVerificationExpired)
		return
	}

	userEl, err := user.FindByID(model.UserID)
	if err != nil {
		Abort(c, err)
		return
	}

	err = user.EmailVerified(userEl.ID, model.Email)
	if err != nil {
		Abort(c, err)
		return
	}

	err = emailverification.DeleteByToken(token)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Your email successfully verified"))
}
