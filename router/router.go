package router

import (
	"github.com/miru-project/bt-server/app"
	handlers "github.com/miru-project/bt-server/handler"
)

func init() {
	app.App.Get("/", handlers.Hello)
	app.App.Get("/version", handlers.Version)
	app.App.Get("/status", handlers.TorrentStatus)
	app.App.Get("/torrent", handlers.ListTorrent)
	app.App.Post("/torrent", handlers.AddTorrent)
	app.App.Delete("/torrent/:infoHash/*", handlers.DeleteTorrent)
	app.App.Get("/torrent/:infoHash/*", handlers.GetTorrentData)
}
