package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/responses"
)

func abort400(c *gin.Context, err error) {
	abort(c, errors.New400(err.Error()))
}

func abort422(c *gin.Context, err error) {
	abort(c, errors.New422(err.Error()))
}

func abort(c *gin.Context, err error) {
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
