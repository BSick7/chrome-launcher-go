package launch

import (
	"runtime"
)

var (
	defaultLocalFlags = []string{
		// Disable built-in Google Translate service
		"--disable-translate",
		// Disable all chrome extensions entirely
		"--disable-extensions",
		// Disable various background network services, including extension updating,
		//   safe browsing service, upgrade detector, translate, UMA
		"--disable-background-networking",
		// Disable fetching safebrowsing lists, likely redundant due to disable-background-networking
		"--safebrowsing-disable-auto-update",
		// Disable syncing to a Google account
		"--disable-sync",
		// Disable reporting to UMA, but allows for collection
		"--metrics-recording-only",
		// Disable installation of default apps on first run
		"--disable-default-apps",
		// Mute any audio
		"--mute-audio",
		// Skip first run wizards
		"--no-first-run",
	}
)

func newLocalLauncher(cfg Config) Launcher {
	return &baseLauncher{
		chromePath:    cfg.ChromePath,
		requestedPort: cfg.RequestedPort,
		startingUrl:   cfg.StartingUrl,
		userDataDir:   cfg.UserDataDir,
		flags:         localFlags(cfg),
		env:           cfg.Env(),
		debug:         GetChromeLogLevelNumber(cfg.LogLevel) < 3,
	}
}

func localFlags(cfg Config) []string {
	flags := defaultLocalFlags
	if cfg.IgnoreDefaultFlags {
		flags = []string{}
	}
	if runtime.GOOS == "linux" {
		flags = append(flags, "--disable-setuid-sandbox")
	}
	flags = append(flags, cfg.ChromeFlags...)
	flags = append(flags, cfg.LogFlags()...)
	return flags
}
