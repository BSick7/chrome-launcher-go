package find

import (
	"os"
	"path"
)

/**
 * Look for windows executables in 2 ways
 * 1. Look into CHROME_PATH env variable
 * 2. Look through app data, program files, program files (x86)
 */
func allChromeInstallations() []string {
	// 1. Look into CHROME_PATH env variable
	installations := fromEnv()

	// 2. Look through app data, program files, program files (x86)
	for _, inst := range fromPrograms() {
		if canAccess(inst) {
			installations = append(installations, inst)
		}
	}

	return installations
}

func fromPrograms() []string {
	suffixes := []string{
		path.Join("", "Google", "Chrome SxS", "Application", "chrome.exe"),
		path.Join("", "Google", "Chrome", "Application", "chrome.exe"),
	}
	prefixes := []string{
		os.Getenv("LOCALAPPDATA"),
		os.Getenv("PROGRAMFILES"),
		os.Getenv("PROGRAMFILES(X86)"),
	}

	paths := make([]string, 0)
	for _, prefix := range prefixes {
		if prefix == "" {
			continue
		}
		for _, suffix := range suffixes {
			paths = append(paths, path.Join(prefix, suffix))
		}
	}
	return paths
}
