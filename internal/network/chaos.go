package network

import (
	"math/rand"
	"net"
	"time"
)

func InjectChaos(conn net.Conn) {
	rand.Seed(time.Now().UnixNano())

	threshold := rand.Intn(100)

	if threshold < 5 {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		conn.Close()
	} else if threshold < 10 {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
	}
}
