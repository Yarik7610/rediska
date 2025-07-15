package resp

import "fmt"

type Value interface {
	Encode() ([]byte, error)
	Decode([]byte) ([]byte, Value, error)
}

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (rc *Controller) Decode(b []byte) (Value, error) {
	l := len(b)
	if l == 0 {
		return nil, fmt.Errorf("expected not fully empty string")
	}

	var res Value
	var err error

	switch b[0] {
	case '*':
		_, res, err = Array{}.Decode(b)
	case '$':
		_, res, err = BulkString{}.Decode(b)
	case ':':
		_, res, err = Integer{}.Decode(b)
	case '+':
		_, res, err = SimpleString{}.Decode(b)
	case '-':
		_, res, err = SimpleError{}.Decode(b)
	default:
		return nil, fmt.Errorf("detected unknown RESP type")
	}

	return res, err
}

func StrPtr(s string) *string {
	return &s
}
