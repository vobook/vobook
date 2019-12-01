package birthdaynotifier

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"github.com/vovainside/vobook/config"
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
	birthdaynotificationlog "github.com/vovainside/vobook/domain/birthday_notification_log"
	contactpropertytype "github.com/vovainside/vobook/enum/contact_property_type"
	"github.com/vovainside/vobook/enum/gender"
	"github.com/vovainside/vobook/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	checkInterval = 1 * time.Hour
)

var tbot *tb.Bot

// days to notify before birthday
// todo: should be configurable
var notifyDaysBefore = []int{10, 3, 0}

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
	var elems []models.Contact
	err := database.Conn().Model(&elems).
		Join("LEFT JOIN birthday_notification_logs bnl ON bnl.contact_id=contact.id").
		Where("bnl.created_at is NULL OR (date_part('year', now())::text || '-' || date_part('month', bnl.created_at)::text || '-' || date_part('day', bnl.created_at)::text)::date - now()::date NOT IN (?)", pg.In(notifyDaysBefore)).
		Where("bnl.created_at IS NULL OR bnl.created_at::date < now()").
		Where("contact.dob_month IS NOT NULL").
		Where("contact.dob_day IS NOT NULL").
		Where("contact.deleted_at IS NULL").
		// TODO
		// this shit is ugly (and probably slow on large datasets)
		// make it better if you can
		WhereIn("(date_part('year', now())::TEXT || '-' || dob_month::TEXT || '-' || dob_day::TEXT)::DATE - now()::DATE IN (?)", []int{10, 3, 0}).
		Relation("Props").
		Relation("User").
		Select()
	if err != nil {
		log.Error(err)
		return
	}

	// send to telegram
	var wg sync.WaitGroup
	for _, el := range elems {
		wg.Add(1)
		go func(m models.Contact) {
			defer wg.Done()
			err := sendToTelegram(m)
			if err != nil {
				log.Error(err)
				return
			}

			err = birthdaynotificationlog.Create(m.ID)
			if err != nil {
				log.Error(err)
				return
			}
		}(el)
	}
	wg.Wait()
}

func sendToTelegram(c models.Contact) (err error) {
	if c.User.TelegramID == 0 {
		return
	}

	msg := buildMessage(c)

	user := &tb.User{
		ID: c.User.TelegramID,
	}

	_, err = tbot.Send(user, msg)
	return
}

func buildMessage(c models.Contact) (msg string) {
	year := c.DOBYear
	if year == 0 {
		year = time.Now().Year()
	}

	bornAt := time.Date(year, c.DOBMonth, c.DOBDay, 0, 0, 0, 0, time.UTC)
	nextAt, daysLeft, age := utils.BirthdayInfo(bornAt)

	if daysLeft == 0 {
		msg = "Сегодня " + c.Name + " отмечает день рождения\n"
		if c.DOBYear != 0 {
			switch c.Gender {
			case gender.Male:
				msg += "Ему исполнилось"
			case gender.Female:
				msg += "Ей исполнилось"
			default:
				msg += "Исполнилось"
			}
			msg += fmt.Sprintf(" %d лет\n", age)
		}
	} else {
		age++
		msg = fmt.Sprintf("Через %d %s, %s, %s отмечает день рождения\n", daysLeft, wordForm(daysLeft, "день", "дня", "дней"), atWeekday(nextAt.Weekday()), c.Name)
		if c.DOBYear != 0 {
			switch c.Gender {
			case gender.Male:
				msg += "Ему исполнится"
			case gender.Female:
				msg += "Ей исполнится"
			default:
				msg += "Исполнится"
			}
			msg += fmt.Sprintf(" %d лет\n", age)
		}
	}

	dateFormat := "2 January"
	if c.DOBYear != 0 {
		dateFormat += " 2006"
	}

	msg += "Дата рождения: " + bornAt.Format(dateFormat) + "\n"

	phone := ""
	for _, prop := range c.Props {
		if prop.Type == contactpropertytype.Phone {
			phone = prop.Value
			break
		}
	}
	if phone != "" {
		msg += "Тел: " + phone + "\n"
	}

	return
}

func atWeekday(w time.Weekday) string {
	switch w {
	case time.Monday:
		return "в понедельник"
	case time.Tuesday:
		return "во вторник"
	case time.Wednesday:
		return "в среду"
	case time.Thursday:
		return "в четверг"
	case time.Friday:
		return "в пятницу"
	case time.Saturday:
		return "в субботу"
	case time.Sunday:
		return "в воскресенье"
	}

	return ""
}

func wordForm(n int, f1, f2, f5 string) string {
	n = int(math.Abs(float64(n))) % 100
	x := n % 10
	if n > 10 && n < 20 {
		return f5
	}
	if x > 1 && x < 5 {
		return f2
	}
	if x == 1 {
		return f1
	}

	return f5
}
