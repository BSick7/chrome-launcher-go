package main

import (
	"time"

	"github.com/BSick7/chrome-launcher-go/find"
	"github.com/BSick7/chrome-launcher-go/launch"
)

func main() {
	// This is a simple test
	// This should launch google in a chrome browser and close after 3s
	launcher := launch.New(launch.Config{
		StartingUrl: "https://www.google.com",
		ChromePath:  find.Chrome(),
	})
	defer launcher.Kill()
	if err := launcher.Launch(); err != nil {
		panic(err)
	}
	if err := launch.ProbeUntil(launcher.Port(), 2*time.Second); err != nil {
		panic("could not connect to remote debugger")
	}
	<-time.After(3 * time.Second)
}
