package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayEncode(t *testing.T) {
	tests := []struct {
		Name        string
		In          Value
		Expected    []byte
		ShouldError bool
	}{
		{
			Name:        "Simple Array with mixed types",
			In:          Array{Value: []Value{BulkString{Value: strPtr("hello")}, Integer{Value: 123}, SimpleString{Value: "PONG"}}},
			Expected:    []byte("*3\r\n$5\r\nhello\r\n:123\r\n+PONG\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Nil Array",
			In:          CreateBulkStringArray(),
			Expected:    []byte(NULL_ARRAY_RESP_2),
			ShouldError: false,
		},
		{
			Name:        "Array with nested Array",
			In:          Array{Value: []Value{BulkString{Value: strPtr("test")}, Array{Value: []Value{BulkString{Value: strPtr("inner")}, Integer{Value: 42}}}}},
			Expected:    []byte("*2\r\n$4\r\ntest\r\n*2\r\n$5\r\ninner\r\n:42\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Deep nested Array",
			In:          Array{Value: []Value{Array{Value: []Value{Array{Value: []Value{Integer{Value: 1}}}}}}},
			Expected:    []byte("*1\r\n*1\r\n*1\r\n:1\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Array with special characters in bulk string",
			In:          CreateBulkStringArray("hello!@#$%", "test\t123"),
			Expected:    []byte("*2\r\n$10\r\nhello!@#$%\r\n$8\r\ntest\t123\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Array with empty bulk string",
			In:          CreateBulkStringArray(""),
			Expected:    []byte("*1\r\n" + NULL_BULK_STRING_RESP_2),
			ShouldError: false,
		},
		{
			Name:        "Array with bulk string containing CRLF",
			In:          CreateBulkStringArray("hello\r\nworld"),
			Expected:    []byte("*1\r\n$12\r\nhello\r\nworld\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Array with invalid simple string (contains CR)",
			In:          Array{Value: []Value{SimpleString{Value: "hello\rworld"}}},
			Expected:    nil,
			ShouldError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			out, err := test.In.Encode()

			if test.ShouldError {
				assert.NotNil(t, err)
				assert.Nil(t, out)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.Expected, out)
			}
		})
	}
}

func TestArrayDecode(t *testing.T) {
	tests := []struct {
		Name         string
		In           []byte
		ExpectedRest []byte
		Expected     Value
		ShouldError  bool
	}{
		{
			Name:         "Simple Array with mixed types",
			In:           []byte("*3\r\n$5\r\nhello\r\n:123\r\n+PONG\r\n"),
			ExpectedRest: []byte{},
			Expected:     Array{Value: []Value{BulkString{Value: strPtr("hello")}, Integer{Value: 123}, SimpleString{Value: "PONG"}}},
			ShouldError:  false,
		},
		{
			Name:         "Empty Array",
			In:           []byte("*0\r\n"),
			ExpectedRest: []byte{},
			Expected:     Array{Value: []Value{}},
			ShouldError:  false,
		},
		{
			Name:         "Nil Array",
			In:           []byte(NULL_ARRAY_RESP_2),
			ExpectedRest: []byte{},
			Expected:     Array{Value: nil},
			ShouldError:  false,
		},
		{
			Name:         "Array with integers",
			In:           []byte("*3\r\n:1\r\n:-2\r\n:0\r\n"),
			ExpectedRest: []byte{},
			Expected:     Array{Value: []Value{Integer{Value: 1}, Integer{Value: -2}, Integer{Value: 0}}},
			ShouldError:  false,
		},
		{
			Name:         "Array with simple strings and errors",
			In:           []byte("*2\r\n+OK\r\n-ERR invalid\r\n"),
			ExpectedRest: []byte{},
			Expected:     Array{Value: []Value{SimpleString{Value: "OK"}, SimpleError{Value: "ERR invalid"}}},
			ShouldError:  false,
		},
		{
			Name:         "Array with nested Array",
			In:           []byte("*2\r\n$4\r\ntest\r\n*2\r\n$5\r\ninner\r\n:42\r\n"),
			ExpectedRest: []byte{},
			Expected:     Array{Value: []Value{BulkString{Value: strPtr("test")}, Array{Value: []Value{BulkString{Value: strPtr("inner")}, Integer{Value: 42}}}}},
			ShouldError:  false,
		},
		{
			Name:         "Deep nested Array",
			In:           []byte("*1\r\n*1\r\n*1\r\n:1\r\n"),
			ExpectedRest: []byte{},
			Expected:     Array{Value: []Value{Array{Value: []Value{Array{Value: []Value{Integer{Value: 1}}}}}}},
			ShouldError:  false,
		},
		{
			Name:         "Array with special characters",
			In:           []byte("*2\r\n$10\r\nhello!@#$%\r\n$8\r\ntest\t123\r\n"),
			ExpectedRest: []byte{},
			Expected:     CreateBulkStringArray("hello!@#$%", "test\t123"),
			ShouldError:  false,
		},
		{
			Name:         "Array with string containing CRLF",
			In:           []byte("*1\r\n$12\r\nhello\r\nworld\r\n"),
			ExpectedRest: []byte{},
			Expected:     CreateBulkStringArray("hello\r\nworld"),
			ShouldError:  false,
		},
		{
			Name:         "Array with remaining bytes",
			In:           []byte("*1\r\n$5\r\nhello\r\n:123\r\n"),
			ExpectedRest: []byte(":123\r\n"),
			Expected:     CreateBulkStringArray("hello"),
			ShouldError:  false,
		},
		{
			Name:         "Empty input",
			In:           []byte(""),
			ExpectedRest: nil,
			Expected:     nil,
			ShouldError:  true,
		},
		{
			Name:         "Invalid prefix",
			In:           []byte("$5\r\nhello\r\n"),
			ExpectedRest: nil,
			Expected:     nil,
			ShouldError:  true,
		},
		{
			Name:         "Missing CRLF after count",
			In:           []byte("*2\r\n$5\r\nhello"),
			ExpectedRest: nil,
			Expected:     nil,
			ShouldError:  true,
		},
		{
			Name:         "Invalid count",
			In:           []byte("*abc\r\n$5\r\nhello\r\n"),
			ExpectedRest: nil,
			Expected:     nil,
			ShouldError:  true,
		},
		{
			Name:         "Invalid element",
			In:           []byte("*2\r\n$5\r\nhello\r\n#invalid\r\n"),
			ExpectedRest: nil,
			Expected:     nil,
			ShouldError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			a := Array{}
			rest, out, err := a.Decode(test.In)

			if test.ShouldError {
				assert.NotNil(t, err)
				assert.Nil(t, out)
				assert.Nil(t, rest)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.Expected, out)
				assert.Equal(t, test.ExpectedRest, rest)
			}
		})
	}
}
