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

func (*Controller) Decode(b []byte) (rest []byte, value Value, err error) {
	l := len(b)
	if l == 0 {
		return nil, nil, fmt.Errorf("expected not fully empty string")
	}

	switch b[0] {
	case '*':
		return Array{}.Decode(b)
	case '$':
		return BulkString{}.Decode(b)
	case ':':
		return Integer{}.Decode(b)
	case '+':
		return SimpleString{}.Decode(b)
	case '-':
		return SimpleError{}.Decode(b)
	default:
		return nil, nil, fmt.Errorf("detected unknown RESP type")
	}
}

func StrPtr(s string) *string {
	return &s
}
