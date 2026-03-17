package client

import (
	"fmt"
	"net"
	"strings"
)

type Sender struct{
	conn net.Conn
}

func NewSender(conn net.Conn) *Sender {
	return &Sender{conn: conn}
}

func (s *Sender) Send(data Data) error {
	var buf strings.Builder
	buf.WriteString(data.Command)
	buf.WriteString(" ")
	buf.WriteString(data.Key)
	
	if data.Value != nil {
		buf.WriteString(" ")
		buf.WriteString(*data.Value)
	}
	if data.TTL != nil {
		buf.WriteString(" ")
		fmt.Fprintf(&buf, "%d", *data.TTL)
	}

	buf.WriteString("\n")


	_, err := s.conn.Write([]byte(buf.String()))
	if err != nil {
		return err
	}
	return nil
}