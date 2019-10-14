package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/cmd/server/responses"
	"github.com/vovainside/vobook/config"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/domain/auth_token"
	"github.com/vovainside/vobook/domain/email_verification"
	"github.com/vovainside/vobook/domain/password_reset"
	"github.com/vovainside/vobook/domain/user"
	"github.com/vovainside/vobook/services/mail"
	"github.com/vovainside/vobook/utils"
)

func RegisterUser(c *gin.Context) {
	var req requests.RegisterUser
	if !bindJSON(c, &req) {
		return
	}

	elem, err := req.ToModel()
	if err != nil {
		abort422(c, err)
		return
	}

	_, err = user.FindByEmail(req.Email)
	if err == nil {
		Abort(c, errors.ReqisterUserEmailExists)
		return
	}

	err = user.Create(elem)
	if err != nil {
		Abort(c, err)
		return
	}

	token, err := emailverification.Create(elem.ID, elem.Email)
	if err != nil {
		Abort(c, err)
		return
	}

	err = mail.SendTemplate(elem.Email, "email-verification", mail.Replace{
		"link": fmt.Sprintf(config.Get().WebClientAddr+"/verify-email/%s", token.Token),
	})
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, elem)
}

func ChangeEmail(c *gin.Context) {
	var req requests.ChangeEmail
	if !bindJSON(c, &req) {
		return
	}

	userEl := AuthUser(c)
	err := bcrypt.CompareHashAndPassword([]byte(userEl.Password), []byte(req.Password))
	if err != nil {
		abort422(c, errors.WrongPassword)
		return
	}

	existing, err := user.FindByEmail(req.Email)
	if err != nil && err != errors.UserByEmailNotFound {
		Abort(c, err)
		return
	}
	if err == nil {
		if existing.ID != userEl.ID {
			err = errors.EmailChangeEmailInUser
		} else {
			err = errors.EmailChangeSameEmail
		}
		if err != nil {
			Abort(c, err)
			return
		}
	}

	token, err := emailverification.Create(userEl.ID, req.Email)
	if err != nil {
		Abort(c, err)
		return
	}

	err = mail.SendTemplate(req.Email, "change-email-verification", mail.Replace{
		"link": fmt.Sprintf(config.Get().WebClientAddr+"/verify-email/%s", token.Token),
	})
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("To change your email please click the link we just sent to your new email address"))
}

func ChangePassword(c *gin.Context) {
	var req requests.ChangeUserPassword
	if !bindJSON(c, &req) {
		return
	}

	elem := AuthUser(c)
	err := bcrypt.CompareHashAndPassword([]byte(elem.Password), []byte(req.OldPassword))
	if err != nil {
		abort422(c, errors.WrongPassword)
		return
	}

	password, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		Abort(c, err)
		return
	}

	err = user.UpdatePassword(elem.ID, password)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Your password successfully changed"))
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
	c.JSON(http.StatusOK, AuthUser(c))
}

func UpdateAuthUser(c *gin.Context) {
	var req requests.UpdateUser
	if !bindJSON(c, &req) {
		return
	}

	userEl := AuthUser(c)
	req.ToModel(&userEl)

	err := user.Update(&userEl)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Saved"))
}

func DeleteAuthUser(c *gin.Context) {
	var req requests.DeleteUser
	if !bindJSON(c, &req) {
		return
	}

	userEl := AuthUser(c)

	err := bcrypt.CompareHashAndPassword([]byte(userEl.Password), []byte(req.Password))
	if err != nil {
		abort422(c, errors.WrongEmailOrPassword)
		return
	}

	err = user.Delete(userEl.ID)
	if err != nil {
		Abort(c, err)
		return
	}

	err = authtoken.DeleteByUser(userEl.ID)
	if err != nil {
		Abort(c, err)
		return
	}

	err = mail.SendTemplate(userEl.Email, "user-account-deleted", nil)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Your account deleted"))
}

func ResetPasswordStart(c *gin.Context) {
	var req requests.ResetPasswordStart
	if !bindJSON(c, &req) {
		return
	}

	elem, err := user.FindByEmail(req.Email)
	if err != nil {
		Abort(c, err)
		return
	}

	token := &models.PasswordReset{
		UserID: elem.ID,
	}
	err = passwordreset.Create(token)
	if err != nil {
		Abort(c, err)
		return
	}

	err = mail.SendTemplate(elem.Email, "password-reset", mail.Replace{
		"link": fmt.Sprintf(config.Get().WebClientAddr+"/reset-password/%s", token.Token),
	})
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Password change confirmation sent to your email"))
}

func ResetPasswordCheckToken(c *gin.Context) {
	token := c.Param("token")
	elem, err := passwordreset.Find(token)
	if err != nil {
		Abort(c, err)
		return
	}

	if elem.ExpiresAt.Before(time.Now()) {
		Abort(c, errors.PasswordResetTokenExpired)
		return
	}

	c.JSON(http.StatusOK, "OK")
}

func ResetPassword(c *gin.Context) {
	var req requests.ResetPassword
	if !bindJSON(c, &req) {
		return
	}

	elem, err := passwordreset.Find(req.Token)
	if err != nil {
		Abort(c, err)
		return
	}

	if elem.ExpiresAt.Before(time.Now()) {
		Abort(c, errors.PasswordResetTokenExpired)
		return
	}

	userEl, err := user.FindByID(elem.UserID)
	if err != nil {
		Abort(c, errors.PasswordResetTokenExpired)
		return
	}

	password, err := utils.HashPassword(req.Password)
	if err != nil {
		Abort(c, err)
		return
	}

	err = user.UpdatePassword(userEl.ID, password)
	if err != nil {
		Abort(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.OK("Your password successfully changed"))
}
