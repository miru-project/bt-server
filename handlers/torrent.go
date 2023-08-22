package handlers

import (
	"bytes"
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/miru-project/bt-server/app"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func TorrentStatus(c *fiber.Ctx) error {
	return c.JSON(app.BTClient.Stats())
}

func ListTorrent(c *fiber.Ctx) error {
	return c.JSON(app.Torrents)
}

func AddTorrent(c *fiber.Ctx) error {
	body := c.Body()
	mediaInfo, err := metainfo.Load(bytes.NewReader(body))
	if err != nil {
		return err
	}
	t, err := app.BTClient.AddTorrent(mediaInfo)
	if err != nil {
		return err
	}
	app.Torrents[t.InfoHash().HexString()] = t
	return c.JSON(t.Info())
}

func GetTorrentData(c *fiber.Ctx) error {
	params := c.AllParams()
	infoHash := params["infoHash"]
	path := params["*1"]
	t, ok := app.Torrents[infoHash]
	if !ok {
		return c.SendString("torrent not found")
	}
	if path == "" {
		return c.JSON(t.Info())
	}
	files := t.Files()
	unescape, err := url.PathUnescape(path)
	log.Info(unescape)
	if err != nil {
		log.Error(err.Error())
	}
	if len(files) == 0 && unescape == t.Name() {
		return serverTorrentData(c, t.NewReader(), t.Length())
	}
	for _, file := range files {
		if file.DisplayPath() == unescape {
			return serverTorrentData(c, file.NewReader(), file.Length())
		}
	}
	return c.SendString("file not found")
}

func serverTorrentData(c *fiber.Ctx, reader torrent.Reader, fileSize int64) error {
	c.Set("Content-Type", "video/mp4")

	defer func(reader torrent.Reader) {
		err := reader.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}(reader)
	reader.SetResponsive()

	rangeHeader := c.Get("Range")
	if rangeHeader != "" {
		ranges := strings.Split(rangeHeader, "=")[1]
		rangeParts := strings.Split(ranges, "-")
		start, _ := strconv.ParseInt(rangeParts[0], 10, 64)
		end := fileSize - 1
		if rangeParts[1] != "" {
			end, _ = strconv.ParseInt(rangeParts[1], 10, 64)
		}

		c.Status(http.StatusPartialContent)
		c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		c.Set("Accept-Ranges", "bytes")
		c.Set("Content-Length", strconv.FormatInt(end-start+1, 10))

		log.Info(fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))

		_, err := reader.Seek(start, 0)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Internal server error")
		}
		return c.SendStream(reader)
	}

	return c.SendStream(reader)
}
