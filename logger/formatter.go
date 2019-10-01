package logger

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

var (
	lastLogDate string
)

type Formatter struct {
}

func (f *Formatter) Format(entry *log.Entry) ([]byte, error) {
	var m string
	t := entry.Time

	dateFormatted := fmt.Sprintf("%s %d, %d", t.Month(), t.Day(), t.Year())

	if lastLogDate != dateFormatted {
		m = "---------------------------------------\n" +
			dateFormatted +
			"\n---------------------------------------\n"
	}

	tf := fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
	m = fmt.Sprintf("%s[%s] %s\n", m, tf, entry.Message)

	lastLogDate = dateFormatted

	return []byte(m), nil
}
