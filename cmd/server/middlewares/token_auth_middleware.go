package middlewares

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/cmd/server/handlers"
	authtoken "github.com/vovainside/vobook/domain/auth_token"
)

func TokenAuth(c *gin.Context) {
	sig := c.Request.Header.Get("Authorization")
	if sig == "" {
		handlers.Abort(c, errors.AuthTokenMissing)
		return
	}

	sig = strings.TrimPrefix(sig, "Bearer ")
	if len(sig) < 128 {
		handlers.Abort(c, errors.AuthTokenInvalid)
		return
	}

	token := sig[:64]
	elem, err := authtoken.Find(token)
	if err != nil {
		handlers.Abort(c, err)
		return
	}

	if sig != authtoken.Sign(&elem) {
		handlers.Abort(c, errors.AuthTokenInvalid)
		return
	}

	if elem.ExpiresAt.Before(time.Now()) {
		handlers.Abort(c, errors.AuthTokenExpired)
		return
	}

	c.Set("user", *elem.User)
	c.Next()
}
