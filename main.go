package main

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/miru-project/bt-server/app"
	"github.com/miru-project/bt-server/config"
	_ "github.com/miru-project/bt-server/router"
)

func main() {
	err := app.App.Listen(config.SERVER_LISTEN)
	if err != nil {
		log.Error(err.Error())
	}
}
