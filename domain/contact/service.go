package contact

import (
	"github.com/go-pg/pg/orm"

	"github.com/vovainside/vobook/cmd/server/requests"
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/filters"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/domain/contact_property"
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
	q.Apply(filters.TrashedFilter(req.Trashed))

	q.Where("user_id = ?", userID)

	if len(req.Query) > 2 {
		q.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q.WhereOr("first_name ilike ?", "%"+req.Query+"%")
			q.WhereOr("last_name ilike ?", "%"+req.Query+"%")
			q.WhereOr("name ilike ?", "%"+req.Query+"%")
			return q, nil
		})
		q.Relation("Props", func(q *orm.Query) (*orm.Query, error) {
			q.Where("value ?", "%"+req.Query+"%")
			q.Order("order ASC")
			return q, nil
		})
	} else {
		q.Relation("Props")
	}

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

func Props(id string) (elems []models.ContactProperty, err error) {
	err = database.ORM().
		Model(&elems).
		Where("contact_id = ?", id).
		Order("order ASC").
		Select()

	return
}
