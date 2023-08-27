package config

import "os"

const VERSION = "0.0.1"

var (
	SECRET                 = getVar("SECRET", "")
	SERVER_LISTEN          = getVar("SERVER_LISTEN", "127.0.0.1:3000")
	DATA_DIR               = getVar("DATA_DIR", "./bg-server-data")
	AUTO_DELETE_CACHE_FILE = getVar("AUTO_DELETE_CACHE_FILE", "true")
)

func getVar(Key string, defaultVal string) string {
	if val, ok := os.LookupEnv(Key); ok {
		return val
	}
	return defaultVal
}
