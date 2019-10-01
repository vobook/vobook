package services

import (
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
)

func CreateUser(elem *models.User) (err error) {
	db := database.Conn()
	_, err = db.Model(elem).Insert()
	return
}

func FindUserByEmail(email string) (elem models.User, err error) {
	err = database.ORM().Model(&elem).Where("email=?", email).First()
	return
}
