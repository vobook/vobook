package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vobook/config"
	"vobook/database"

	"github.com/davecgh/go-spew/spew"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	userID             = 0
	checkEventInterval = time.Hour * 24
)

type Event struct {
	ID          int64
	Name        string
	Description string
	Date        time.Time
}

func main() {
	quit := make(chan os.Signal, 1)
	exit := make(chan struct{})
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	b, err := tb.NewBot(tb.Settings{
		Token:  config.Get().TelegramBotAPI,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/ping", func(m *tb.Message) {
		message, err := b.Send(m.Sender, "pong")
		if err != nil {
			println(err.Error())
		}
		spew.Dump(message)
	})

	go CheckEventsStartWorker(b, exit)

	go func() {
		<-quit
		close(exit)
		b.Stop()
	}()

	b.Start()
	println("Telegram bot stopped")
}

func CheckEvents(b *tb.Bot) (err error) {
	conn := database.Conn()

	events := make([]Event, 0)
	err = conn.Model(&events).Column("*").Where("date >= NOW() - INTERVAL '30 days'").Select()
	if err != nil {
		return
	}

	for _, elem := range events {
		//Mon Jan 2 15:04:05 -0700 MST 2006
		message := elem.Date.Format("2 Jan 2006 " + elem.Name)
		err = send(b, message)
		if err != nil {
			println("err sending message ", err.Error())
		}
	}

	return err
}

func CheckEventsStartWorker(b *tb.Bot, exit <-chan struct{}) {
	ticker := time.NewTicker(checkEventInterval)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-exit:
			break loop
		case <-ticker.C:
			go func() {
				err := CheckEvents(b)
				if err != nil {
					err = send(b, err.Error())
					if err != nil {
						println(err.Error())
					}
				}
			}()
		}
	}

	println("Worker stopped")
}

func send(b *tb.Bot, m string) (err error) {
	user := &tb.User{
		ID: userID,
	}
	_, err = b.Send(user, m)
	return
}
