package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/services"
)

func RegisterUser(c *gin.Context) {
	var req requests.RegisterUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := req.Validate()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	_, err = services.FindUserByEmail(req.Email)
	if err == nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "user with email "+req.Email+" already exists")
		return
	}

	err = services.CreateUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
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
