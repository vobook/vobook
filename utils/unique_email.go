package utils

import (
	"strings"

	fake "github.com/brianvoe/gofakeit"
	"github.com/go-pg/pg"
	"github.com/vovainside/vobook/database"
)

// UniqueEmail creates random unique email for given table & column
func UniqueEmail(table string, col ...string) (string, error) {
	length := 2
	column := "email"
	if len(col) == 1 {
		column = col[0]
	}

	for {
		email := RandomString(length) + "@" + RandomString(length) + "." + fake.DomainSuffix()
		email = strings.ToLower(email)
		_, err := database.ORM().ExecOne("SELECT * FROM ? WHERE ? = ?", pg.Ident(table), pg.Ident(column), email)
		if err == pg.ErrNoRows {
			return email, nil
		}

		length++
	}
}
