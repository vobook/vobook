package authtoken

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/go-pg/pg"

	"github.com/vovainside/vobook/cmd/server/errors"
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/utils"
)

func Create(elem *models.AuthToken) (err error) {
	elem.Token = fmt.Sprintf("%x", sha256.Sum256([]byte(utils.RandomString(32))))

	_, err = database.ORM().
		Model(elem).
		Insert()

	return
}

func Find(token string) (elem models.AuthToken, err error) {
	err = database.ORM().
		Model(&elem).
		Relation("User").
		Where("token = ?", token).
		First()
	if err == pg.ErrNoRows {
		err = errors.AuthTokenNotFound
	}

	return
}

func Sign(elem *models.AuthToken) string {
	sigData := []string{
		elem.ID,
		elem.UserID,
		fmt.Sprintf("%d", elem.ClientID),
		elem.UserAgent,
		elem.ClientIP,
		elem.Token,
		elem.CreatedAt.UTC().Format(time.RFC3339),
		elem.ExpiresAt.UTC().Format(time.RFC3339),
	}

	sig := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(sigData, "+"))))
	return elem.Token + sig
}

func DeleteByUser(id string) (err error) {
	_, err = database.ORM().
		Model(&models.AuthToken{}).
		Where("user_id = ?", id).
		Delete()

	return
}
