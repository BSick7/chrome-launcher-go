package find

import (
	"os"

	"github.com/mgutz/logxi/v1"
)

func canAccess(path string) bool {
	if path == "" {
		return false
	}
	if fi, err := os.Open(path); err != nil {
		return false
	} else {
		fi.Close()
		return true
	}
}

func fromEnv() []string {
	if path := os.Getenv("CHROME_PATH"); canAccess(path) {
		return []string{path}
	}
	if path := os.Getenv("LIGHTHOUSE_CHROMIUM_PATH"); canAccess(path) {
		log.Warn("LIGHTHOUSE_CHROMIUM_PATH is deprecated, use CHROME_PATH env variable instead.")
		return []string{path}
	}
	return []string{}
}
