package handlers

import (
	"net/http"
	"time"

	"vobook/cmd/server/errors"
	"vobook/cmd/server/responses"
	"vobook/config"
	emailverification "vobook/domain/email_verification"
	"vobook/domain/user"

	"github.com/gin-gonic/gin"
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
