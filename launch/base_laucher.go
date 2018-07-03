package launch

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
)

type baseLauncher struct {
	// config
	chromePath    string
	requestedPort int
	startingUrl   string
	userDataDir   string
	flags         []string
	env           []string
	debug         bool

	// resource tracking
	cmd            *exec.Cmd
	outFile        *os.File
	errFile        *os.File
	pidFile        string
	port           int
	inited         bool
	createdUserDir bool
}

func (l *baseLauncher) Port() int {
	return l.port
}

func (l *baseLauncher) Pid() int {
	if l.cmd == nil || l.cmd.Process == nil {
		return -1
	}
	return l.cmd.Process.Pid
}

func (l *baseLauncher) Launch() error {
	if l.requestedPort != 0 {
		// If an explicit port is passed, first look for an open connection...
		if l.port = l.requestedPort; Probe(l.port) {
			return nil
		}
	}

	if err := l.init(); err != nil {
		return fmt.Errorf("error initializing: %s", err)
	}
	if err := l.spawn(); err != nil {
		l.dumpChromeLogs()
		return fmt.Errorf("error spawning chrome: %s", err)
	}

	return nil
}

func (l *baseLauncher) Kill() error {
	if l.cmd == nil {
		return nil
	}

	defer func() {
		l.cmd = nil
		l.destroy()
	}()
	if err := l.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("error killing chrome: %s", err)
	}
	return nil
}

func (l *baseLauncher) init() error {
	if l.inited {
		return nil
	}

	if l.createdUserDir = l.userDataDir == ""; l.createdUserDir {
		l.userDataDir = os.TempDir()
		os.MkdirAll(l.userDataDir, os.ModePerm)
	}

	var err error
	outFilename := path.Join(l.userDataDir, "chrome-out.log")
	if l.outFile, err = os.OpenFile(outFilename, os.O_APPEND|os.O_CREATE, 0666); err != nil {
		return fmt.Errorf("error creating chrome out file %q: %s", outFilename, err)
	}
	errFilename := path.Join(l.userDataDir, "chrome-err.log")
	if l.errFile, err = os.OpenFile(errFilename, os.O_APPEND|os.O_CREATE, 0666); err != nil {
		return fmt.Errorf("error creating chrome err file %q: %s", errFilename, err)
	}

	l.pidFile = path.Join(l.userDataDir, "chrome.pid")
	l.inited = true
	return nil
}

func (l *baseLauncher) spawn() error {
	if l.cmd != nil {
		return nil
	}

	// If requested port is not set, the launcher is responsible for generating
	// We perform instead of chrome so we know which port is taken
	if l.requestedPort == 0 {
		if l.port = getRandomPort(); l.port <= 0 {
			return fmt.Errorf("could not find available port")
		}
	}

	args := append(l.flags, fmt.Sprintf(`--remote-debugging-port=%d`, l.port))
	if l.createdUserDir {
		args = append(args, fmt.Sprintf(`--user-data-dir=%s`, l.userDataDir))
	}
	args = append(args, l.startingUrl)

	cmd := exec.Command(l.chromePath, args...)
	cmd.Stdin = nil
	cmd.Stdout = l.outFile
	cmd.Stderr = l.errFile
	cmd.Env = l.env
	if l.debug {
		log.Println("exec", cmd.Path, cmd.Args)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting: %s", err)
	}
	l.cmd = cmd

	if l.debug {
		log.Println("pid", cmd.Process.Pid)
	}

	if err := ioutil.WriteFile(l.pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0666); err != nil {
		return fmt.Errorf("error writing pid file: %s", err)
	}
	return nil
}

func (l *baseLauncher) destroy() {
	// Only clean up the tmp dir if we created it
	if !l.createdUserDir {
		return
	}

	if l.outFile != nil {
		l.outFile.Close()
		l.outFile = nil
	}

	if l.errFile != nil {
		l.errFile.Close()
		l.errFile = nil
	}

	os.RemoveAll(l.userDataDir)
	l.inited = false
}

func (l *baseLauncher) dumpChromeLogs() {
	if !l.debug {
		return
	}
	outlogs, _ := ioutil.ReadFile(l.outFile.Name())
	log.Println("out", l.outFile.Name(), string(outlogs))
	errlogs, _ := ioutil.ReadFile(l.errFile.Name())
	log.Println("err", l.errFile.Name(), string(errlogs))
}
