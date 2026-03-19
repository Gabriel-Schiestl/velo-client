package client

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
)

type Data struct {
	Command string
	Key string
	Value any
	TTL *uint64
}

type valueType int

const (
	stringType valueType = iota
	byteSliceType
	otherType
)

type value struct {
	buf []byte
	valueType valueType
	len [2]byte
}

func getValueBytes(val any) (*value, error) {
	var buf bytes.Buffer
	var valType valueType

	switch v := val.(type) {
	case string:
		buf.WriteString(v)
		valType = stringType
	case []byte:
		buf.Write(v)
		valType = byteSliceType
	default:
		encoder := gob.NewEncoder(&buf)
		if err := encoder.Encode(val); err != nil {
			return nil, err
		}
		valType = otherType
	}

	len, err := getLengthUint16FromValue(buf.Bytes())
	if err != nil {
		return nil, err
	}

	value := new(value)
	value.buf = buf.Bytes()
	value.valueType = valType
	value.len = len

	return value, nil
}

func parseValue(data []byte) (any, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("data too short to contain value type and length")
	}

	valType := valueType(data[0])
	valLen := binary.BigEndian.Uint16(data[1:3])
	if len(data) < int(3+valLen) {
		return nil, fmt.Errorf("data too short to contain value of declared length")
	}

	valData := data[3 : 3+valLen]

	switch valType {
	case stringType:
		return string(valData), nil
	case byteSliceType:
		return valData, nil
	case otherType:
		var value any
		buf := bytes.NewBuffer(valData)
		decoder := gob.NewDecoder(buf)
		if err := decoder.Decode(&value); err != nil {
			return nil, err
		}
		return value, nil
	default:
		return nil, fmt.Errorf("unknown value type: %d", valType)
	}
}

func getLengthUint16FromValue(value []byte) ([2]byte, error) {
	if len(value) > 65535 {
		return [2]byte{}, fmt.Errorf("value length exceeds uint16 limit")
	}

	var lenBuf [2]byte
	binary.BigEndian.PutUint16(lenBuf[:], uint16(len(value)))
	return lenBuf, nil
}

func getBigEndianFromUint64(value uint64) ([8]byte, error) {
	if value > 4294967295 {
		return [8]byte{}, fmt.Errorf("value exceeds uint32 limit")
	}

	var lenBuf [8]byte
	binary.BigEndian.PutUint64(lenBuf[:], value)
	return lenBuf, nil
}