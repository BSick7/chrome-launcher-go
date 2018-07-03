package launch

import (
	"fmt"
	"net"
	"time"
)

// Probes chrome remote debugger to verify it's listening
func Probe(port int) bool {
	return ProbeUntil(port, 500*time.Millisecond) == nil
}

// Probes chrome remote debugger until a connection is verified or wait expires
func ProbeUntil(port int, wait time.Duration) error {
	deadline := time.Now().Add(wait)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), time.Second)
		if err != nil {
			if deadline.Before(time.Now()) {
				return fmt.Errorf("debugger dial error: %s", err)
			}
		} else if conn != nil {
			conn.Close()
			return nil
		}
	}
}
