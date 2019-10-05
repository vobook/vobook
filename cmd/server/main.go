package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/vovainside/vobook/cmd/server/routes"
	"github.com/vovainside/vobook/config"
)

func main() {
	conf := config.Get()
	r := gin.Default()

	routes.Register(r)

	server := &http.Server{
		Addr:    conf.Server.Host + ":" + conf.Server.Port,
		Handler: r,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		if err := server.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err.Error())
	}

	log.Println(conf.App.Name + "." + conf.App.Env + " server closed")
}
