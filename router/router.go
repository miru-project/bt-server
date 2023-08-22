package router

import (
	"github.com/miru-project/bt-server/app"
	"github.com/miru-project/bt-server/handlers"
)

func init() {
	app.App.Get("/", handlers.Hello)
	app.App.Get("/status", handlers.TorrentStatus)
	app.App.Get("/torrent", handlers.ListTorrent)
	app.App.Post("/torrent", handlers.AddTorrent)
	app.App.Get("/torrent/:infoHash/*", handlers.GetTorrentData)
}
