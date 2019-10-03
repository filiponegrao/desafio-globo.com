package main

import (
	"os"
	"strconv"

	"github.com/filiponegrao/convivva-server/config"
	"github.com/filiponegrao/desafio-globo.com/controllers"
	"github.com/filiponegrao/desafio-globo.com/db"
	"github.com/filiponegrao/desafio-globo.com/server"
)

var conf config.Configuration
var confEmail controllers.EmailConfiguration

func main() {

	configFileName := "config.json"
	initialConfigure(configFileName)

	database := db.Connect()
	s := server.Setup(database)

	port := conf.ApiPort

	if p := os.Getenv("PORT"); p != "" {
		if _, err := strconv.Atoi(p); err == nil {
			port = p
		}
	}

	s.Run(":" + port)
}

func initialConfigure(path string) {
	conf = config.Get(path)

	confEmail = controllers.EmailConfiguration{
		Mail:     conf.Email,
		Password: conf.EmailPassword,
		Server:   conf.EmailServer,
		Port:     conf.EmailPort,
		Site:     conf.Site,
	}
	controllers.ConfigEmailEngine(confEmail)
}
