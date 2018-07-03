package launch

var (
	defaultLambdaFlags = []string{
		"--disable-gpu",
		"--single-process", // Currently wont work without this :-(
		// https://groups.google.com/a/chromium.org/d/msg/headless-dev/qqbZVZ2IwEw/Y95wJUh2AAAJ
		"--no-zygote", // helps avoid zombies
		"--no-sandbox",
		"--disable-setuid-sandbox",

		// disable shared memory in lambda
		//   - https://medium.com/@marco.luethy/running-headless-chrome-on-aws-lambda-fa82ad33a9eb
		//   - https://github.com/GoogleChrome/puppeteer/issues/1834
		"--disable-dev-shm-usage",
	}
)

func newLambdaLauncher(cfg Config) Launcher {
	return &baseLauncher{
		chromePath:    cfg.ChromePath,
		requestedPort: cfg.RequestedPort,
		startingUrl:   cfg.StartingUrl,
		userDataDir:   cfg.UserDataDir,
		flags:         lambdaFlags(cfg),
		env:           cfg.Env(),
		debug:         GetChromeLogLevelNumber(cfg.LogLevel) < 3,
	}
}

func lambdaFlags(cfg Config) []string {
	flags := defaultLambdaFlags
	if cfg.IgnoreDefaultFlags {
		flags = []string{}
	}
	flags = append(flags, cfg.ChromeFlags...)
	flags = append(flags, cfg.LogFlags()...)
	return flags
}
