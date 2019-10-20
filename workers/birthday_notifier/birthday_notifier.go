package birthdaynotifier

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/vovainside/vobook/config"
	"github.com/vovainside/vobook/database"
	"github.com/vovainside/vobook/database/models"
	"github.com/vovainside/vobook/domain/birthday_notification_log"
	"github.com/vovainside/vobook/enum/contact_property_type"
	"github.com/vovainside/vobook/enum/gender"
	"github.com/vovainside/vobook/utils"
)

const (
	checkInterval = 1 * time.Hour
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
	var elems []models.Contact
	err := database.Conn().Model(&elems).
		Join("LEFT JOIN birthday_notification_logs bnl ON bnl.contact_id=contact.id").
		Where("bnl.created_at is null OR (extract('year', now())::text || '-' || extract('month', bnl.created_at)::text || '-' || extract('day', bnl.created_at)::text)::date - now()::date NOT IN (?)", pg.In([]int{10, 3, 0})).
		//Where("bnl.created_at is null OR bnl.created_at::date < now()").
		Where("contact.birthday is not null").
		// TODO
		// this shit is ugly oO (and probably slow on large datasets)
		// make it better if you can
		WhereIn("(extract('year', now())::text || '-' || extract('month', birthday)::text || '-' || extract('day', birthday)::text)::date - now()::date IN (?)", []int{10, 3, 0}).
		Relation("Props").
		Relation("User").
		Select()
	if err != nil {
		log.Error(err)
		return
	}

	println("Записей:", len(elems))

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

	_ = buildMessage(c)

	//user := &tb.User{
	//	ID: c.User.TelegramID,
	//}

	//_, err = tbot.Send(user, msg)
	return
}

func buildMessage(c models.Contact) (msg string) {
	dayAt, days, age := utils.BirthdayInfo(c.Birthday)

	if days == 0 {
		msg = "Сегодня " + c.Name + " отмечает день рождения\n"
		switch c.Gender {
		case gender.Male:
			msg += "Ему исполнилось"
		case gender.Female:
			msg += "Ей исполнилось"
		default:
			msg += "Исполнилось"
		}
	} else {
		age++
		msg = fmt.Sprintf("Через %d %s, %s, %s отмечает день рождения\n", days, wordForm(days, "день", "дня", "дней"), atWeekday(dayAt.Weekday()), c.Name)
		switch c.Gender {
		case gender.Male:
			msg += "Ему исполнится"
		case gender.Female:
			msg += "Ей исполнится"
		default:
			msg += "Исполнится"
		}
	}

	msg += fmt.Sprintf(" %d лет\n", age)
	msg += "Дата рождения: " + c.Birthday.Format("2 January 2006") + "\n"

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

	println(msg)
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
