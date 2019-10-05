package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/cmd/server/responses"
	"github.com/vovainside/vobook/database/models"
	authtoken "github.com/vovainside/vobook/domain/auth_token"
	emailverification "github.com/vovainside/vobook/domain/email_verification"
	"github.com/vovainside/vobook/domain/user"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var req requests.RegisterUser
	if !bindJSON(c, &req) {
		return
	}

	u, err := req.ToUser()
	if err != nil {
		abort422(c, err)
		return
	}

	_, err = user.FindByEmail(req.Email)
	if err == nil {
		Abort(c, errors.ReqisterUserEmailExists)
		return
	}

	err = user.Create(u)
	if err != nil {
		Abort(c, err)
		return
	}

	err = emailverification.Create(u.ID, u.Email)
	if err != nil {
		Abort(c, err)
		return
	}

	// TODO send email verification email

	c.JSON(http.StatusOK, u)
}

func Login(c *gin.Context) {
	var req requests.Login
	if !bindJSON(c, &req) {
		return
	}

	elem, err := user.FindByEmail(req.Email)
	if err != nil {
		Abort(c, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(elem.Password), []byte(req.Password))
	if err != nil {
		abort422(c, errors.WrongEmailOrPassword)
		return
	}

	token := &models.AuthToken{
		UserID:    elem.ID,
		ClientID:  models.ClientID(c.GetInt("clientID")),
		ClientIP:  c.Request.RemoteAddr,
		UserAgent: c.Request.UserAgent(),
	}
	err = authtoken.Create(token)
	if err != nil {
		Abort(c, err)
		return
	}

	resp := responses.Login{
		User:      elem,
		Token:     authtoken.Sign(token),
		ExpiresAt: token.ExpiresAt,
	}

	c.JSON(http.StatusOK, resp)
}

func GetAuthUser(c *gin.Context) {
	c.JSON(http.StatusOK, authUser(c))
}
