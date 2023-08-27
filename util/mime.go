package util

var mediaMimeType = map[string]string{
	".mp4":  "video/mp4",
	".webm": "video/webm",
	".ogg":  "video/ogg",
	".flv":  "video/x-flv",
	".mov":  "video/quicktime",
	".ts":   "video/MP2T",
	".3gp":  "video/3gpp",
	".avi":  "video/x-msvideo",
	".wmv":  "video/x-ms-wmv",
	".mkv":  "video/x-matroska",
	".mpg":  "video/mpeg",
	".mpeg": "video/mpeg",
	".m4v":  "video/x-m4v",
	".mvp":  "video/x-matroska",
	".flac": "audio/flac",
	".mp3":  "audio/mpeg",
	".wav":  "audio/wav",
	".m4a":  "audio/x-m4a",
}

// 传入文件后缀
func IsMedia(fileExtension string) (string, bool) {
	mime, ok := mediaMimeType[fileExtension]
	return mime, ok
}
