package utils

import (
	"fmt"
	"math/rand"

	"github.com/go-pg/pg"
	"github.com/vovainside/vobook/database"
)

const (
	UniqueTokenColumnDefault = "token"
	UniqueTokenLengthDefault = 64
)

type UniqueTokenOpts struct {
	Column string
	Length int
}

// UniqueToken creates unique token for given table.column
func UniqueToken(table string, opts ...UniqueTokenOpts) (string, error) {
	var opt UniqueTokenOpts
	if len(opts) == 1 {
		opt = opts[0]
	} else {
		opt = UniqueTokenOpts{
			Column: UniqueTokenColumnDefault,
			Length: UniqueTokenLengthDefault,
		}
	}

	if opt.Length < 1 {
		return "", fmt.Errorf("cannot make token with %d chars", opt.Length)
	}

	maxTries := len(Chars()) * opt.Length
	for i := 0; i < maxTries; i++ {
		token := RandomToken(opt.Length)
		_, err := database.ORM().ExecOne("SELECT * FROM ? WHERE ? = ?", pg.F(table), pg.F(opt.Column), token)
		if err == pg.ErrNoRows {
			return token, nil
		}

		if err != nil {
			return "", err
		}
	}

	return "", fmt.Errorf("unable to find unique token within %d attempts", maxTries)
}

func RandomToken(length int) string {
	chars := Chars()
	token := make([]byte, length)
	for i := 0; i < length; i++ {
		token[i] = chars[rand.Intn(len(chars))]
	}
	return string(token)
}
