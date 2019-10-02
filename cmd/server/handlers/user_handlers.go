package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/services"
)

func RegisterUser(c *gin.Context) {
	var req requests.RegisterUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		abort400(c, err)
		return
	}

	user, err := req.Validate()
	if err != nil {
		abort422(c, err)
		return
	}

	_, err = services.FindUserByEmail(req.Email)
	if err == nil {
		abort(c, errors.ReqisterUserEmailExists)
		return
	}

	err = services.CreateUser(user)
	if err != nil {
		abort(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func SearchUsers(c *gin.Context) {

}

func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	c.JSON(200, "user id is "+id)
}
