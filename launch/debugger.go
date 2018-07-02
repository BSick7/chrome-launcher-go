package launch

import (
	"fmt"
	"log"
	"net"
	"time"
)

type debugger struct {
	port  int
	debug bool
}

func (d *debugger) IsReady() bool {
	return d.WaitUntilReady(500 * time.Millisecond)
}

func (d *debugger) WaitUntilReady(maxWait time.Duration) bool {
	deadline := time.Now().Add(maxWait)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", d.port), time.Second)
		if err != nil {
			if deadline.Before(time.Now()) {
				if d.debug {
					log.Println("debugger dial error", err)
				}
				return false
			}
		} else if conn != nil {
			conn.Close()
			return true
		}
	}
}
