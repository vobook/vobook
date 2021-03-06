package contactproperty

import (
	"github.com/go-pg/pg/v9"

	"vobook/database"
	"vobook/database/models"
)

func CreateMany(elems *[]models.ContactProperty) (err error) {
	if len(*elems) == 0 {
		return
	}
	_, err = database.ORM().
		Model(elems).
		Insert()

	return
}

func Find(id string) (elem models.ContactProperty, err error) {
	err = database.ORM().
		Model(&elem).
		Where("contact_property.id = ?", id).
		Relation("Contact").
		First()

	return
}

func Update(elem *models.ContactProperty) (err error) {
	err = database.ORM().
		Update(elem)

	return
}

func Trash(userID string, ids ...string) (err error) {
	_, err = database.ORM().Exec(
		"UPDATE contact_properties SET deleted_at=NOW() WHERE id IN (?) AND contact_id IN "+
			"(SELECT id FROM contacts WHERE user_id=?)",
		pg.In(ids), userID)

	return
}

func Restore(userID string, ids ...string) (err error) {
	_, err = database.ORM().Exec(
		"UPDATE contact_properties SET deleted_at=null WHERE id IN (?) AND contact_id IN "+
			"(SELECT id FROM contacts WHERE user_id=?)",
		pg.In(ids), userID)

	return
}

func Delete(userID string, ids ...string) (err error) {
	_, err = database.ORM().Exec(
		"DELETE FROM contact_properties WHERE id IN (?) AND contact_id IN "+
			"(SELECT id FROM contacts WHERE user_id=?)",
		pg.In(ids), userID)

	return
}

func DeleteByContact(contactID string) (err error) {
	_, err = database.ORM().Exec(
		"DELETE FROM contact_properties WHERE contact_id=?", contactID)

	return
}

func Reorder(userID string, ids []string) (err error) {
	for i, id := range ids {
		_, err = database.ORM().Exec(
			"UPDATE contact_properties SET \"order\"=? WHERE id = ? AND contact_id IN "+
				"(SELECT id FROM contacts WHERE user_id=?)",
			i, id, userID)
		if err != nil {
			return
		}
	}

	return
}
