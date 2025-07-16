package resp

import (
	"fmt"
	"log"
)

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

func CreateArray(args ...string) Array {
	var values []Value
	for _, arg := range args {
		values = append(values, BulkString{Value: StrPtr(arg)})
	}
	return Array{Value: values}
}

func AssertEqualSimpleString(value Value, raw string) {
	v, ok := value.(SimpleString)
	if !ok {
		log.Fatalf("assertion failed: expected SimpleString, got: %T", value)
	}
	if v.Value != raw {
		log.Fatalf("assertion failed: expected %s, got: %v", raw, v.Value)
	}
}

func StrPtr(s string) *string {
	return &s
}
