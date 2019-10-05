package responses

import (
	"time"

	"github.com/vovainside/vobook/database/models"
)

type Login struct {
	User      models.User `json:"user"`
	Token     string      `json:"token"`
	ExpiresAt time.Time   `json:"expires"`
}
