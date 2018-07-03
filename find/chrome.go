package find

func AllChromeInstallations() []string {
	return allChromeInstallations()
}

func Chrome() string {
	chromes := allChromeInstallations()
	if len(chromes) > 0 {
		return chromes[0]
	}
	return ""
}
