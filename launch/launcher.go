package launch

type Launcher interface {
	// Launches chrome
	// If a debugger port is requested and chrome is listening on that port,
	//   this will not create any resources or launch chrome
	Launch() error

	// Stops chrome and destroys resources
	// This will do nothing if Launch did not launch a new chrome
	Kill() error

	// Chrome process id
	Pid() int

	// Chrome remote debugger port
	Port() int
}

func New(cfg Config) Launcher {
	if cfg.IsLambdaEnv {
		return newLambdaLauncher(cfg.Normalize())
	}
	return newLocalLauncher(cfg.Normalize())
}
