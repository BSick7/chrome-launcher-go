package main

import (
	"time"

	"github.com/BSick7/chrome-launcher-go/launch"
)

func main() {
	// This is a simple test
	// This should launch google in a chrome browser and close after 5s
	launcher := launch.New(launch.Config{
		StartingUrl: "https://www.google.com",
	})
	defer launcher.Kill()
	if err := launcher.Launch(); err != nil {
		panic(err)
	}
	<-time.After(5 * time.Second)
}
