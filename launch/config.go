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
	ChromeLogLevelInfo    ChromeLogLevel = "info"
	ChromeLogLevelWarning ChromeLogLevel = "warning"
	ChromeLogLevelError   ChromeLogLevel = "error"
	ChromeLogLevelFatal   ChromeLogLevel = "fatal"

	DefaultMaxConnectWait = 30 * time.Second

	DefaultPort = 9222
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
	LauncherDebug        bool
}

func (c Config) Normalize() Config {
	normalized := c

	if normalized.StartingUrl == "" {
		normalized.StartingUrl = "about:blank"
	}
	if normalized.LogLevel == "" {
		normalized.LogLevel = ChromeLogLevelFatal
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

	flags = append(flags, c.LogFlags()...)
	flags = append(flags, c.ChromeFlags...)
	flags = append(flags, c.StartingUrl)

	return flags
}

func (c Config) LogFlags() []string {
	switch c.LogLevel {
	case ChromeLogLevelInfo:
		return []string{
			"--enable-logging",
			"--log-level=0",
		}
	case ChromeLogLevelWarning:
		return []string{
			"--enable-logging",
			"--log-level=1",
		}
	case ChromeLogLevelError:
		return []string{
			"--enable-logging",
			"--log-level=2",
		}
	case ChromeLogLevelFatal:
		return []string{
			"--enable-logging",
			"--log-level=3",
		}
	default:
		return []string{}
	}
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
