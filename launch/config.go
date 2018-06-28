package launch

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type ChromeLogLevel string

var (
	ChromeLogLevelVerbose ChromeLogLevel = "verbose"
	ChromeLogLevelInfo    ChromeLogLevel = "info"
	ChromeLogLevelError   ChromeLogLevel = "error"
	ChromeLogLevelSilent  ChromeLogLevel = "silent"

	DefaultMaxConnectWait = 30 * time.Second
)

type Config struct {
	StartingUrl          string
	LogLevel             ChromeLogLevel
	ChromeFlags          []string
	RequestedPort        int
	ChromePath           string
	IgnoreDefaultFlags   bool
	MaxConnectWait       time.Duration
	EnvironmentVariables map[string]string
	UserDataDir          string
}

func (c Config) Normalize() Config {
	normalized := c

	if normalized.StartingUrl == "" {
		normalized.StartingUrl = "about:blank"
	}
	if normalized.LogLevel == "" {
		normalized.LogLevel = ChromeLogLevelSilent
	}
	if normalized.ChromeFlags == nil {
		normalized.ChromeFlags = []string{}
	}
	if normalized.MaxConnectWait == 0 {
		normalized.MaxConnectWait = DefaultMaxConnectWait
	}
	if normalized.EnvironmentVariables == nil {
		normalized.EnvironmentVariables = envMapFromProcess()
	}

	return normalized
}

func (c Config) Flags(port int) []string {
	flags := DefaultChromeFlags
	if c.IgnoreDefaultFlags {
		flags = []string{}
	}
	flags = append(flags, fmt.Sprintf(`--remote-debugging-port=%d`, port))

	if runtime.GOOS == "linux" {
		flags = append(flags, "--disable-setuid-sandbox")
	}

	if c.UserDataDir != "" {
		udd, _ := filepath.Abs(c.UserDataDir)
		flags = append(flags, fmt.Sprintf(`--user-data-dir=%s`, udd))
	}

	flags = append(flags, c.ChromeFlags...)
	flags = append(flags, c.StartingUrl)

	return flags
}

func (c Config) Env() []string {
	e := make([]string, 0)
	for k, v := range c.EnvironmentVariables {
		e = append(e, fmt.Sprintf(`%s=%s`, k, v))
	}
	return e
}

func envMapFromProcess() map[string]string {
	env := map[string]string{}
	for _, v := range os.Environ() {
		tokens := strings.SplitN(v, "=", 2)
		if len(tokens) != 2 {
			continue
		}
		env[tokens[0]] = tokens[1]
	}
	return env
}
