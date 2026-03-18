package connection

import (
	"net"
	"sync"
)

var conn net.Conn
var once sync.Once

func GetConn(addr string) net.Conn {
	once.Do(
		func() {
			conn = connect(addr)
		},
	)
	return conn
}

func connect(addr string) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	return conn
}