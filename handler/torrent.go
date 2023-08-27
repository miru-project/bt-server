package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/miru-project/bt-server/app"
	"github.com/miru-project/bt-server/config"
	models "github.com/miru-project/bt-server/model"
	"github.com/miru-project/bt-server/util"
)

func TorrentStatus(c *fiber.Ctx) error {
	return c.JSON(app.BTClient.Stats())
}

func ListTorrent(c *fiber.Ctx) error {
	var result []models.TorrentResult
	for hash, t := range app.Torrents {
		result = append(result, models.TorrentResult{
			InfoHash: hash,
			Name:     t.Name(),
		})
	}
	return c.JSON(result)
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
	hex := t.InfoHash().HexString()

	app.Torrents[hex] = t

	files := []string{}
	if len(t.Info().Files) == 0 {
		files = append(files, t.Name())
	} else {
		for _, file := range t.Info().Files {
			files = append(files, file.DisplayPath(t.Info()))
		}
	}

	return c.JSON(models.TorrentDetailResult{
		InfoHash: hex,
		Detail:   t.Info(),
		Files:    files,
	})
}

func DeleteTorrent(c *fiber.Ctx) error {
	infoHash := c.Params("infoHash")
	t, ok := app.Torrents[infoHash]
	if !ok {
		return c.Status(http.StatusNotFound).SendString("torrent not found")
	}
	t.Drop()

	if config.AUTO_DELETE_CACHE_FILE == "true" {
		if err := os.RemoveAll(path.Join(config.DATA_DIR, t.Name())); err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Internal server error")
		}
	}

	delete(app.Torrents, infoHash)
	return c.SendStatus(http.StatusOK)
}

func GetTorrentData(c *fiber.Ctx) error {
	params := c.AllParams()
	infoHash := params["infoHash"]
	filePath := params["*1"]
	t, ok := app.Torrents[infoHash]
	if !ok {
		return c.Status(http.StatusNotFound).SendString("torrent not found")
	}
	if filePath == "" {
		files := []string{}
		if len(t.Info().Files) == 0 {
			files = append(files, t.Name())
		} else {
			for _, file := range t.Info().Files {
				files = append(files, file.DisplayPath(t.Info()))
			}
		}
		return c.JSON(models.TorrentDetailResult{
			InfoHash: infoHash,
			Detail:   t.Info(),
			Files:    files,
		})
	}
	files := t.Files()
	unescape, err := url.PathUnescape(filePath)
	log.Info(unescape)
	if err != nil {
		log.Error(err.Error())
	}
	// 获取文件后缀
	fileExtension := path.Ext(unescape)
	if len(files) == 0 && unescape == t.Name() {
		return serverTorrentData(c, fileExtension, t.NewReader(), t.Length())
	}
	for _, file := range files {
		if file.DisplayPath() == unescape {
			return serverTorrentData(c, fileExtension, file.NewReader(), file.Length())
		}
	}
	return c.Status(http.StatusNotFound).SendString("file not found")
}

func serverTorrentData(c *fiber.Ctx, fileExtension string, reader torrent.Reader, fileSize int64) error {
	// 获取文件后缀
	mime, ok := util.IsMedia(fileExtension)
	log.Info(mime, fileExtension)
	if !ok {
		c.Set("Content-Type", "application/octet-stream")
		return c.SendStream(reader)
	}

	c.Set("Content-Type", mime)
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
