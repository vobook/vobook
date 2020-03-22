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
	resp := responses.Error{}
	code := http.StatusInternalServerError

	switch err.(type) {
	case errors.Error:
		e := err.(errors.Error)
		code = e.Code
		resp.Error = e.Error()
	case errors.List:
		resp.Errors = err.(errors.List)
	case errors.Input:
		resp.InputErrors = err.(errors.Input)
	default:
		resp.Error = err.Error()
		switch err {
		case pg.ErrNoRows:
			code = http.StatusNotFound
		}
	}

	if !config.IsReleaseEnv() {
		log.Println(string(debug.Stack()))
	}

	c.AbortWithStatusJSON(code, resp)
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
