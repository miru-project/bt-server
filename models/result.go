package models

import "github.com/anacrolix/torrent/metainfo"

type TorrentDetailResult struct {
	InfoHash string         `json:"infoHash"`
	Detail   *metainfo.Info `json:"detail"`
}

type TorrentResult struct {
	InfoHash string `json:"infoHash"`
	Name     string `json:"name"`
}
