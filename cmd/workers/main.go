package main

import (
	"os"
	"os/signal"

	"vobook/config"
	"vobook/logger"
	"vobook/services/mail"
	"vobook/workers"

	log "github.com/sirupsen/logrus"
)

func main() {
	conf := config.Get()
	logger.Setup()

	mail.InitDrivers()

	quit := make(chan os.Signal)
	exit := make(chan bool)
	signal.Notify(quit, os.Interrupt)

	workers.Start(exit)

	go func() {
		<-quit
		close(exit)
	}()

	<-exit
	log.Println(conf.App.Name + "." + conf.App.Env + " workers stopped")
}
