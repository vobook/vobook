package birthdaynotifier

import (
	"fmt"
	"time"

	"github.com/vovainside/vobook/enum/gender"
	"github.com/vovainside/vobook/utils"

	"github.com/davecgh/go-spew/spew"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/vovainside/vobook/config"
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"

	log "github.com/sirupsen/logrus"
)

const (
	checkInterval = 1 * time.Minute
)

var tbot *tb.Bot

func Start(exit <-chan bool) {
	var err error
	tbot, err = tb.NewBot(tb.Settings{
		Token:  config.Get().TelegramBotAPI,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		panic(err)
	}

	go worker(exit)
}

func worker(exit <-chan bool) {
	go check()
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	log.Println("Birthday Checker started")
loop:
	for {
		select {
		case <-exit:
			log.Println("Birthday Checker stopped")
			break loop
		case <-ticker.C:
			go check()
		}
	}
}

func check() {
	// birthday in exactly 10 days
	var elems []models.Contact
	err := database.Conn().Model(&elems).
		Where("contact.birthday is not null").
		// TODO
		// this shit is ugly oO (and probably slow on large datasets)
		// make it better if you can
		Where("((extract('year', now())::text || '-' || extract('month', birthday)::text || '-' || extract('day', birthday)::text)::date - interval '10 days') = now()::date").
		Relation("Props").
		Relation("User").
		Select()
	if err != nil {
		log.Error(err)
		return
	}

	// send to telegram
	for _, el := range elems {
		err = sendToTelegram(el)
		log.Error(err)
	}

	//SELECT * FROM public."studentData" where date_part('day',TO_DATE("DOB", 'DD/MM/YYYY'))='20'
	//AND date_part('month',TO_DATE("DOB", 'DD/MM/YYYY'))='04';

	//where age(cd.birthdate) - (extract(year from age(cd.birthdate)) || ' years')::interval = '0'::interval

	spew.Dump(elems)
}

func sendToTelegram(c models.Contact) (err error) {
	if c.User.TelegramID == 0 {
		return
	}

	user := &tb.User{
		ID: c.User.TelegramID,
	}

	dd := utils.DateDiff(c.Birthday, time.Now())

	spew.Dump(dd)

	message := fmt.Sprintf("Через %d дней день рождения у %s\n", dd.Day, c.Name)

	switch c.Gender {
	case gender.Male:
		message += "Ему исполняется"
	case gender.Female:
		message += "Ей исполняется"
	default:
		message += "Исполняется"
	}

	message += fmt.Sprintf(" %d лет\n", dd.Year)

	message += c.Birthday.Format("2 Jan 2006")

	_, err = tbot.Send(user, message)
	return
}
