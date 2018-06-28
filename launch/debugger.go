package launch

import (
	"fmt"
	"net"
	"time"
)

type debugger struct {
	port int
}

func (d *debugger) IsReady() bool {
	return d.WaitUntilReady(500 * time.Millisecond)
}

func (d *debugger) WaitUntilReady(maxWait time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf(":%d", d.port), time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
