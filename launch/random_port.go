package launch

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

var (
	portRand    = rand.New(rand.NewSource(time.Now().UnixNano()))
	initialPool []int
)

func init() {
	minPort := 1024
	maxPort := 65535

	initialPool = make([]int, 0)
	for i := minPort; i <= maxPort; i++ {
		initialPool = append(initialPool, i)
	}
}

func getRandomPort() int {
	isAvailable := func(port int) bool {
		conn, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}

	pool := make([]int, len(initialPool))
	copy(pool, initialPool)

	for len(pool) > 0 {
		i := portRand.Intn(len(pool))
		cur := pool[i]
		pool = append(pool[:i], pool[i+1:]...)
		if isAvailable(cur) {
			return cur
		}
	}
	return -1
}
