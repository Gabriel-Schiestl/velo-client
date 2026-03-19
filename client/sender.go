package client

import (
	"bytes"
	"encoding/gob"
	"net"
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

	var intTTL *uint64
	if(ttl != nil) {
		ttlMiliSeconds := ttl.Milliseconds()
		ttlSeconds := ttlMiliSeconds / 1000

		intTTL = new(uint64)
		*intTTL = uint64(ttlSeconds)
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
	var buf bytes.Buffer

	buf.WriteByte(byte(len(data.Command)))
	buf.WriteString(data.Command)

	buf.WriteByte(byte(len(data.Key)))
	buf.WriteString(data.Key)

	if data.Value != nil {
		buf.WriteByte(1)

		value, err := getValueBytes(data.Value)
		if err != nil {
			return err
		}

		buf.WriteByte(byte(value.valueType))
		buf.Write(value.len[:])
		buf.Write(value.buf)
	} else {
		buf.WriteByte(0)
	}

	if data.TTL != nil {
		buf.WriteByte(1)

		ttlBuf, err := getBigEndianFromUint64(*data.TTL)
		if err != nil {
			return err
		}
		buf.Write(ttlBuf[:])
	} else {
		buf.WriteByte(0)
	}

	_, err := c.conn.Write(buf.Bytes())
	return err
}