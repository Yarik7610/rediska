package resp

import "fmt"

type Value interface {
	Encode() ([]byte, error)
	Decode([]byte) ([]byte, Value, error)
}

type RESPController struct {
	SimpleString SimpleString
	SimpleError  SimpleError
	Integer      Integer
	BulkString   BulkString
	Array        Array
}

var Controller = RESPController{}

func (rc *RESPController) Decode(b []byte) (Value, error) {
	l := len(b)
	if l == 0 {
		return nil, fmt.Errorf("expected not fully empty string")
	}

	var res Value
	var err error

	switch b[0] {
	case '*':
		_, res, err = rc.Array.Decode(b)
	case '$':
		_, res, err = rc.BulkString.Decode(b)
	case ':':
		_, res, err = rc.Integer.Decode(b)
	case '+':
		_, res, err = rc.SimpleString.Decode(b)
	case '-':
		_, res, err = rc.SimpleError.Decode(b)
	default:
		return nil, fmt.Errorf("detected unknown RESP type")
	}

	return res, err
}
