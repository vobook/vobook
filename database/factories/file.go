package factories

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"vobook/database"
	"vobook/database/models"
	"vobook/services/fs"
	"vobook/tests/fake"
	"vobook/utils"

	"github.com/brianvoe/gofakeit"
)

func MakeFile(mOpt ...models.File) (m models.File, err error) {
	if len(mOpt) == 1 {
		m = mOpt[0]
	}

	if m.UserID == "" {
		var userEl models.User
		userEl, err = CreateUser()
		if err != nil {
			return
		}
		m.UserID = userEl.ID
	}

	if m.Name == "" {
		m.Name = gofakeit.Name()
	}
	if m.Description == "" {
		m.Description = gofakeit.Name()
	}

	if len(m.Bytes) == 0 {
		m.Bytes, err = fake.Picture()
		if err != nil {
			return
		}
	}

	m.Size = int64(len(m.Bytes))

	sha := sha256.New()
	sha.Write(m.Bytes)
	m.Hash = hex.EncodeToString(sha.Sum(nil))

	if m.Filename == "" {
		m.Filename = utils.RandomString() + ".png"
	}

	m.Type = fs.Type(m.Filename)
	m.Base64 = base64.StdEncoding.EncodeToString(m.Bytes)

	return
}

func CreateFile(mOpt ...models.File) (m models.File, err error) {
	m, err = MakeFile(mOpt...)
	if err != nil {
		return
	}

	savedAs, err := fs.Save(m.Filename, m.Bytes)
	if err != nil {
		return
	}

	m.Path = savedAs

	err = database.ORM().Insert(&m)
	return
}
