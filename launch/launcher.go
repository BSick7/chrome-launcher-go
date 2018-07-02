package launch

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/BSick7/chrome-launcher-go/find"
)

type Launcher interface {
	Launch() error
	Kill() error
	Pid() int
	Port() int
}

func New(cfg Config) Launcher {
	cfg = cfg.Normalize()

	return &launcher{
		cfg: cfg,
	}
}

type launcher struct {
	cmd *exec.Cmd
	cfg Config

	outFile        *os.File
	errFile        *os.File
	pidFile        string
	port           int
	inited         bool
	createdUserDir bool
}

func (l *launcher) Port() int {
	return l.port
}

func (l *launcher) Pid() int {
	if l.cmd == nil || l.cmd.Process == nil {
		return -1
	}
	return l.cmd.Process.Pid
}

func (l *launcher) init() error {
	if l.inited {
		return nil
	}

	if l.cfg.ChromePath == "" {
		installations := find.Chrome()
		if len(installations) <= 0 {
			return fmt.Errorf("could not find chrome")
		}
		l.cfg.ChromePath = installations[0]
	}

	if l.createdUserDir = l.cfg.UserDataDir == ""; l.createdUserDir {
		l.cfg.UserDataDir = os.TempDir()
		os.MkdirAll(l.cfg.UserDataDir, os.ModePerm)
	}

	var err error
	outFilename := path.Join(l.cfg.UserDataDir, "chrome-out.log")
	if l.outFile, err = os.OpenFile(outFilename, os.O_APPEND|os.O_CREATE, 0666); err != nil {
		return fmt.Errorf("error creating chrome out file %q: %s", outFilename, err)
	}
	errFilename := path.Join(l.cfg.UserDataDir, "chrome-err.log")
	if l.errFile, err = os.OpenFile(errFilename, os.O_APPEND|os.O_CREATE, 0666); err != nil {
		return fmt.Errorf("error creating chrome err file %q: %s", errFilename, err)
	}

	l.pidFile = path.Join(l.cfg.UserDataDir, "chrome.pid")
	l.inited = true
	return nil
}

func (l *launcher) Launch() error {
	if l.cfg.RequestedPort != 0 {
		l.port = l.cfg.RequestedPort
		// If an explicit port is passed, first look for an open connection...
		d := &debugger{port: l.port}
		if d.IsReady() {
			return nil
		}
	}

	if err := l.init(); err != nil {
		return fmt.Errorf("error initializing: %s", err)
	}
	if err := l.spawn(); err != nil {
		return fmt.Errorf("error spawning chrome: %s", err)
	}

	d := &debugger{port: l.port, debug: l.cfg.LauncherDebug}
	if !d.WaitUntilReady(l.cfg.MaxConnectWait) {
		return fmt.Errorf("timed out awaiting debugger connection")
	}
	return nil
}

func (l *launcher) Kill() error {
	if l.cmd == nil {
		return nil
	}

	defer func() {
		l.cmd = nil
		l.destroyTmp()
	}()
	if err := l.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("error killing chrome: %s", err)
	}
	return nil
}

func (l *launcher) destroyTmp() {
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

	os.RemoveAll(l.cfg.UserDataDir)
}

func (l *launcher) spawn() error {
	if l.cmd != nil {
		return nil
	}

	// If requested port is not set, the launcher is responsible for generating
	// We perform instead of chrome so we know which port is taken
	if l.cfg.RequestedPort == 0 {
		if l.port = getRandomPort(); l.port <= 0 {
			return fmt.Errorf("could not find available port")
		}
	}

	cmd := exec.Command(l.cfg.ChromePath, l.cfg.Flags(l.port)...)
	cmd.Stdin = nil
	cmd.Stdout = l.outFile
	cmd.Stderr = l.errFile
	cmd.Env = l.cfg.Env()
	if l.cfg.LauncherDebug {
		log.Println("exec", cmd.Path, cmd.Args)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting: %s", err)
	}
	l.cmd = cmd

	if l.cfg.LauncherDebug {
		log.Println("pid", cmd.Process.Pid)
	}

	if err := ioutil.WriteFile(l.pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0666); err != nil {
		return fmt.Errorf("error writing pid file: %s", err)
	}
	return nil
}
