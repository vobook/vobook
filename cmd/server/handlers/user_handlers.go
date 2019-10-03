package handlers

import (
	"net/http"

	emailverification "github.com/vovainside/vobook/domain/email_verification"

	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/domain/user"
)

func RegisterUser(c *gin.Context) {
	var req requests.RegisterUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		abort400(c, err)
		return
	}

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

	// TODO send email confirmation email

	c.JSON(http.StatusOK, u)
}

func SearchUsers(c *gin.Context) {

}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	c.JSON(200, "user id is "+id)
}
