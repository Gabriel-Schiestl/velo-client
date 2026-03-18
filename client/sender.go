package client

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Gabriel-Schiestl/velo-client/internal/connection"
)

type Client struct{
	conn net.Conn
}

func NewClient(addr string) *Client {
	return &Client{conn: connection.GetConn(addr)}
}

func (c *Client) Set(key string, value any, ttl *time.Duration) error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(value); err != nil {
		return err
	}

	var intTTL *int64
	if(ttl != nil) {
		ttlMiliSeconds := ttl.Milliseconds()
		ttlSeconds := ttlMiliSeconds / 1000

		intTTL = &ttlSeconds
	}

	data := Data{
		Command: "SET",
		Key: key,
		Value: buf.Bytes(),
		TTL: intTTL,
	}

	if err := c.send(data); err != nil {
		return err
	}

	return nil
}

func (c *Client) send(data Data) error {
	var buf strings.Builder
	buf.WriteString(data.Command)
	buf.WriteString(" ")
	buf.WriteString(data.Key)
	
	if data.Value != nil {
		buf.WriteString(" ")
		buf.Write(data.Value)
	}
	if data.TTL != nil {
		buf.WriteString(" ")
		fmt.Fprintf(&buf, "%d", *data.TTL)
	}

	buf.WriteString("\n")

	_, err := c.conn.Write([]byte(buf.String()))
	if err != nil {
		return err
	}
	return nil
}