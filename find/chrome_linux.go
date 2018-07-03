package find

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"regexp"
	"strings"
)

func init() {
	priorities = map[string]int{
		"chrome-wrapper$":       51,
		"google-chrome-stable$": 50,
		"google-chrome$":        49,
		"chromium-browser$":     48,
		"chromium$":             47,
	}
	if path := os.Getenv("LIGHTHOUSE_CHROMIUM_PATH"); path != "" {
		priorities[path] = 100
	}
	if path := os.Getenv("CHROME_PATH"); path != "" {
		priorities[path] = 101
	}
}

/**
 * Look for linux executables in 3 ways
 * 1. Look into CHROME_PATH env variable
 * 2. Look into the directories where .desktop are saved on gnome based distro's
 * 3. Look for google-chrome-stable & google-chrome executables by using the which command
 */
func allChromeInstallations() []string {
	installations := installs{}

	// 1. Look into CHROME_PATH env variable
	installations = append(installations, fromEnv()...)

	// 2. Look into the directories where .desktop are saved on gnome based distro's
	installations = append(installations, fromDesktop()...)

	// 3. Look for google-chrome(-stable) & chromium(-browser) executables by using the which command
	for _, inst := range fromWhich() {
		if canAccess(inst) {
			installations = append(installations, inst)
		}
	}

	return installations.Prioritized()
}

func fromDesktop() []string {
	desktopInstallationFolders := []string{"/usr/share/applications/"}
	if user, err := user.Current(); err == nil {
		desktopInstallationFolders = append(desktopInstallationFolders, path.Join(user.HomeDir, ".local/share/applications/"))
	}
	installations := make([]string, 0)
	for _, dir := range desktopInstallationFolders {
		installations = append(installations, findChromeExecutables(dir)...)
	}
	return installations
}

func findChromeExecutables(dir string) []string {
	argumentsRegex := regexp.MustCompile(`(^[^ ]+).*`) // Take everything up to the first space
	chromeExecRegex := `^Exec=\/.*\/(google-chrome|chrome|chromium)-.*`

	if !canAccess(dir) {
		return []string{}
	}

	// Output of the grep & print looks like:
	//    /opt/google/chrome/google-chrome --profile-directory
	//    /home/user/Downloads/chrome-linux/chrome-wrapper %U
	cmd := exec.Command(fmt.Sprintf(`grep -ER %q %s | awk -F '=' '{print $2}'`, chromeExecRegex, dir))
	if err := cmd.Run(); err != nil {
		// Some systems do not support grep -R so fallback to -r.
		// See https://github.com/GoogleChrome/chrome-launcher/issues/46 for more context.
		cmd = exec.Command(fmt.Sprintf(`grep -Er %q %s | awk -F '=' '{print $2}'`, chromeExecRegex, dir))
		if err := cmd.Run(); err != nil {
			return []string{}
		}
	}

	raw, _ := cmd.Output()
	installations := make([]string, 0)
	for _, cur := range strings.Split(string(raw), `\n`) {
		installations = append(installations, argumentsRegex.ReplaceAllString(cur, "$1"))
	}
	return installations
}

func fromWhich() []string {
	executables := []string{
		"google-chrome-stable",
		"google-chrome",
		"chromium-browser",
		"chromium",
	}
	installations := make([]string, 0)
	for _, ex := range executables {
		cmd := exec.Command("which", ex)
		if err := cmd.Run(); err != nil {
			continue
		}
		raw, err := cmd.Output()
		if err != nil {
			continue
		}
		installations = append(installations, strings.SplitN(string(raw), `\n`, 2)[0])
	}
	return installations
}
