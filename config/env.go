package config

import "os"

var (
	SECRET        = getVar("SECRET", "adsads")
	SERVER_LISTEN = getVar("SERVER_LISTEN", "127.0.0.1:3000")
)

func getVar(Key string, defaultVal string) string {
	if val, ok := os.LookupEnv(Key); ok {
		return val
	}
	return defaultVal
}
