package connection

import (
	"net"
	"sync"
)

var conn any
var once sync.Once

func GetConn() any {
	once.Do(
		func() {
			conn = connect()
		},
	)
	return conn
}

func connect() net.Conn {
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		panic(err)
	}
	return conn
}