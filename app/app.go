package app

import (
	"context"
	"net"

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
	// 修复 Android 的 DNS 问题
	// https://github.com/YouROK/TorrServer/blob/34634649024b06921e7f2bb068da4a4ad35fabe4/server/cmd/main.go#L101-L120
	dnsResolve()
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

func dnsResolve() {
	addrs, err := net.LookupHost("www.google.com")
	if len(addrs) == 0 {
		log.Error("Check dns failed", addrs, err)

		fn := func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", "1.1.1.1:53")
		}

		net.DefaultResolver = &net.Resolver{
			Dial: fn,
		}

		addrs, err = net.LookupHost("www.google.com")
		log.Info("Check cloudflare dns", addrs, err)
	} else {
		log.Info("Check dns OK", addrs, err)
	}
}
