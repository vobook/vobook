package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Generic map
type M map[string]interface{}

var (
	Letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	Digits  = []byte("0123456789")
	Symbols = []byte("!@#$%^&*-=_.?")
)

func Chars() []byte {
	chars := append(Letters, Digits...)
	chars = append(chars, Symbols...)
	return chars
}

func randomString(length int) string {
	alphanums := Letters
	alphanums = append(alphanums, Digits...)
	str := make([]byte, length)
	for i := range str {
		str[i] = alphanums[rand.Intn(len(alphanums))]
	}
	return string(str)
}

func randomLetters(length int) string {
	letters := make([]byte, length)
	for i := range letters {
		letters[i] = Letters[rand.Intn(len(Letters))]
	}
	return string(letters)
}
