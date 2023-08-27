package app

import (
	"github.com/anacrolix/torrent"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/miru-project/bt-server/config"
)

var (
	App      *fiber.App
	BTClient *torrent.Client
	Torrents = make(map[string]*torrent.Torrent)
)

func init() {
	var err error
	App = fiber.New()

	App.Use(logger.New())

	App.Use(func(c *fiber.Ctx) error {
		secret := c.Get("Authorization")
		if secret == "" {
			secret = c.Query("secret")
		}
		if secret != config.SECRET {
			return c.SendStatus(403)
		}
		return c.Next()
	})

	log.Info("Starting BT client")

	cc := torrent.NewDefaultClientConfig()
	cc.DataDir = config.DATA_DIR

	BTClient, err = torrent.NewClient(cc)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("BT client started")
}
