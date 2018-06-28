package find

import (
	"fmt"
	"os/exec"
	"os/user"
	"path"
	"strings"
)

var (
	lsregister = `/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister`
)

func init() {
	priorities = map[string]int{
		`^/Applications/.*Chrome.app`:        100,
		`^/Applications/.*Chrome Canary.app`: 101,
		`^/Volumes/.*Chrome.app`:             -2,
		`^/Volumes/.*Chrome Canary.app`:      -1,
	}

	if user, err := user.Current(); err == nil {
		priorities[fmt.Sprintf(`^%s/Applications/.*Chrome.app`, user.HomeDir)] = 50
		priorities[fmt.Sprintf(`^%s/Applications/.*Chrome Canary.app`, user.HomeDir)] = 51
	}

	if path := os.Getenv("LIGHTHOUSE_CHROMIUM_PATH"); path != "" {
		priorities[path] = 150
	}
	if path := os.Getenv("CHROME_PATH"); path != "" {
		priorities[path] = 151
	}
}

/**
 * Look for darwin executables in 2 ways
 * 1. Look into CHROME_PATH env variable
 * 2. Look through launch services
 */
func chrome() []string {
	installations := installs{}

	// 1. Look into CHROME_PATH env variable
	installations = append(installations, fromEnv()...)

	// 2. Look through launch services
	for _, inst := range fromLaunchServices() {
		if canAccess(inst) {
			installations = append(installations, inst)
		}
	}

	return installations.Prioritized()
}

func fromLaunchServices() []string {
	suffixes := []string{
		"/Contents/MacOS/Google Chrome Canary",
		"/Contents/MacOS/Google Chrome",
	}

	cmd := exec.Command(fmt.Sprintf(`%s -dump | grep -i 'google chrome\( canary\)\?.app$' | awk '{$1=""; print $0}'`, lsregister))
	if err := cmd.Run(); err != nil {
		return []string{}
	}
	raw, _ := cmd.Output()

	installations := make([]string, 0)
	for _, cur := range strings.Split(string(raw), `\n`) {
		for _, suffix := range suffixes {
			installations = append(installations, path.Join(strings.TrimSpace(inst), suffix))
		}
	}
	return installations
}
