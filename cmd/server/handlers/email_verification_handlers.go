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
	id := c.Param("id")

	model, err := emailverification.FindByID(id)
	if err != nil {
		abort(c, err)
		return
	}

	if model.CreatedAt.Add(config.Get().EmailVerificationLifetime).Before(time.Now()) {
		abort(c, errors.EmailVerificationExpired)
		return
	}

	userEl, err := user.FindByID(model.UserID)
	if err != nil {
		abort(c, err)
		return
	}

	err = user.EmailVerified(userEl.ID, model.Email)
	if err != nil {
		abort(c, err)
		return
	}

	err = emailverification.Delete(id)
	if err != nil {
		abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Your email successfully verified"))
}
