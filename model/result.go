package models

import "github.com/anacrolix/torrent/metainfo"

type TorrentDetailResult struct {
	InfoHash string         `json:"infoHash"`
	Detail   *metainfo.Info `json:"detail"`
	Files    []string       `json:"files"`
}

type TorrentResult struct {
	InfoHash string `json:"infoHash"`
	Name     string `json:"name"`
}
