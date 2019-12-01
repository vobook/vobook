package contact

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/filters"
	"github.com/vovainside/vobook/database/models"
	contactproperty "github.com/vovainside/vobook/domain/contact_property"
)

func Create(m *models.Contact) (err error) {
	_, err = database.ORM().
		Model(m).
		Insert()
	if err != nil {
		return
	}

	for i := range m.Props {
		m.Props[i].ContactID = m.ID
	}

	err = contactproperty.CreateMany(&m.Props)
	return
}

func Search(userID string, req requests.SearchContact) (data []models.Contact, count int, err error) {
	q := database.ORM().Model(&data)
	q.Apply(filters.PageFilter(req.Page, req.Limit))
	q.Apply(filters.TrashedFilter(req.Trashed, "contact"))

	q.Where("user_id = ?", userID)

	if len(req.Query) > 2 {
		q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q.Join("JOIN contact_properties cp ON cp.contact_id=contact.id")
			q.WhereOr("first_name ilike ?", "%"+req.Query+"%")
			q.WhereOr("last_name ilike ?", "%"+req.Query+"%")
			q.WhereOr("contact.name ilike ?", "%"+req.Query+"%")
			q.WhereOr("cp.value ilike ?", "%"+req.Query+"%")
			return q, nil
		})
	}
	q.Relation("Props", func(q *orm.Query) (*orm.Query, error) {
		q.Order("order ASC")
		return q, nil
	})

	q.Order("contact.created_at DESC")
	count, err = q.SelectAndCount()
	return
}

func Find(id string) (m models.Contact, err error) {
	err = database.ORM().
		Model(&m).
		Where("id = ?", id).
		First()

	return
}

func Update(m *models.Contact) (err error) {
	err = database.ORM().
		Update(m)

	return
}

func Trash(userID string, ids ...string) (err error) {
	_, err = database.ORM().
		Model(&models.Contact{}).
		Where("user_id = ?", userID).
		WhereIn("id IN (?)", ids).
		Set("deleted_at = ?", time.Now()).
		Update()

	return
}

func Restore(userID string, ids ...string) (err error) {
	_, err = database.ORM().
		Model(&models.Contact{}).
		Where("user_id = ?", userID).
		WhereIn("id IN (?)", ids).
		Set("deleted_at = null").
		Update()

	return
}

func Delete(userID string, ids ...string) (err error) {
	_, err = database.ORM().Exec(`
	DELETE FROM contact_properties WHERE contact_id IN (
		SELECT id FROM contacts WHERE user_id = ?0 AND contact_id IN (?1)
	);
	DELETE FROM contacts WHERE user_id = ?0 AND id IN (?1);
`, userID, pg.In(ids))

	return
}

func Props(id string) (elems []models.ContactProperty, err error) {
	err = database.ORM().
		Model(&elems).
		Where("contact_id = ?", id).
		Order("order ASC").
		Select()

	return
}
