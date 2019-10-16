package birthdaynotifier

import (
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	checkInterval = 1 * time.Minute
)

func Start(exit <-chan bool) {
	go worker(exit)
}

func worker(exit <-chan bool) {
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
	log.Println("checking...")
}
