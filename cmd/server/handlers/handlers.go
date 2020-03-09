package handlers

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"

	"vobook/cmd/server/errors"
	"vobook/cmd/server/responses"
	"vobook/config"
	"vobook/database/models"
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

	if !config.IsReleaseEnv() {
		log.Println(string(debug.Stack()))
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

func bindQuery(c *gin.Context, req Validatable) (ok bool) {
	err := c.ShouldBindQuery(req)
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

func AuthUser(c *gin.Context) models.User {
	elem := c.MustGet("user")
	return elem.(models.User)
}
