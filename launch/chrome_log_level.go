package launch

import (
	"strings"
)

type ChromeLogLevel string

var (
	ChromeLogLevelInfo    ChromeLogLevel = "info"
	ChromeLogLevelWarning ChromeLogLevel = "warning"
	ChromeLogLevelError   ChromeLogLevel = "error"
	ChromeLogLevelFatal   ChromeLogLevel = "fatal"

	levels = map[ChromeLogLevel]int{
		ChromeLogLevelInfo:    0,
		ChromeLogLevelWarning: 1,
		ChromeLogLevelError:   2,
		ChromeLogLevelFatal:   3,
	}
)

// See https://peter.sh/experiments/chromium-command-line-switches/#log-level
func GetChromeLogLevelNumber(level ChromeLogLevel) int {
	if num, ok := levels[level]; ok {
		return num
	}
	return levels[ChromeLogLevelFatal]
}

func NormalizeChromeLogLevel(level ChromeLogLevel) ChromeLogLevel {
	safe := ChromeLogLevel(strings.ToLower(string(level)))
	if _, ok := levels[safe]; ok {
		return safe
	}
	return ChromeLogLevelFatal
}
