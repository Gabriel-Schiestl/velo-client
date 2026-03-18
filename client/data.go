package client

type Data struct {
	Command string
	Key string
	Value []byte
	TTL *int64
}