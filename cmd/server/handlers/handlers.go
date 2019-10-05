package handlers

import (
	"net/http"

	"github.com/vovainside/vobook/database/models"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/responses"
)

type Validatable interface {
	Validate() error
}

func abort400(c *gin.Context, err error) {
	Abort(c, errors.New400(err.Error()))
}

func abort422(c *gin.Context, err error) {
	Abort(c, errors.New422(err.Error()))
}

func Abort(c *gin.Context, err error) {
	e, ok := err.(errors.Error)
	if !ok {
		e = errors.Error{
			Message: err.Error(),
			Err:     err,
		}
		switch err {
		case pg.ErrNoRows:
			e.Code = http.StatusNotFound
			e.Message = "not found"
		default:
			e.Code = http.StatusInternalServerError
		}
	}

	resp := responses.Error{
		Error: e.Error(),
	}

	c.AbortWithStatusJSON(e.Code, resp)
}

func bindJSON(c *gin.Context, req Validatable) (ok bool) {
	err := c.ShouldBindJSON(req)
	if err != nil {
		abort400(c, err)
		return
	}

	err = req.Validate()
	if err != nil {
		abort422(c, err)
		return
	}

	return true
}

func authUser(c *gin.Context) models.User {
	elem := c.MustGet("user")
	return elem.(models.User)
}
