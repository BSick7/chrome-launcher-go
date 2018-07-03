package launch

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	StartingUrl          string
	RequestedPort        int
	ChromePath           string
	ChromeFlags          []string
	LogLevel             ChromeLogLevel
	EnvironmentVariables map[string]string
	UserDataDir          string
	IgnoreDefaultFlags   bool
	IsLambdaEnv          bool
}

func (c Config) Normalize() Config {
	normalized := c

	normalized.LogLevel = NormalizeChromeLogLevel(normalized.LogLevel)
	if normalized.StartingUrl == "" {
		normalized.StartingUrl = "about:blank"
	}
	if normalized.LogLevel == "" {
		normalized.LogLevel = ChromeLogLevelFatal
	}
	if normalized.ChromeFlags == nil {
		normalized.ChromeFlags = []string{}
	}
	if normalized.EnvironmentVariables == nil {
		normalized.EnvironmentVariables = envMapFromProcess()
	}

	return normalized
}

func (c Config) LogFlags() []string {
	num := GetChromeLogLevelNumber(c.LogLevel)
	// No need to set logging flags if requesting 'fatal'
	if num < 3 {
		return []string{
			"--enable-logging",
			fmt.Sprintf("--log-level=%d", num),
			"--v=99",
		}
	}
	return []string{}
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
