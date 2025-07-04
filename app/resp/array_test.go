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
			Name:        "Simple array with mixed types",
			In:          array{value: []Value{bulkString{value: stringPtr("hello")}, integer{value: 123}, simpleString{value: "PONG"}}},
			Expected:    []byte("*3\r\n$5\r\nhello\r\n:123\r\n+PONG\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Empty array",
			In:          array{value: []Value{}},
			Expected:    []byte("*0\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Nil array",
			In:          array{value: nil},
			Expected:    []byte(NULL_ARRAY_RESP_2),
			ShouldError: false,
		},
		{
			Name:        "Array with single element",
			In:          array{value: []Value{bulkString{value: stringPtr("A")}}},
			Expected:    []byte("*1\r\n$1\r\nA\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Array with multiple strings",
			In:          array{value: []Value{bulkString{value: stringPtr("hello")}, bulkString{value: stringPtr("world")}}},
			Expected:    []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Array with integers",
			In:          array{value: []Value{integer{value: 1}, integer{value: -2}, integer{value: 0}}},
			Expected:    []byte("*3\r\n:1\r\n:-2\r\n:0\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Array with simple strings and errors",
			In:          array{value: []Value{simpleString{value: "OK"}, simpleError{value: "ERR invalid"}}},
			Expected:    []byte("*2\r\n+OK\r\n-ERR invalid\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Array with nested array",
			In:          array{value: []Value{bulkString{value: stringPtr("test")}, array{value: []Value{bulkString{value: stringPtr("inner")}, integer{value: 42}}}}},
			Expected:    []byte("*2\r\n$4\r\ntest\r\n*2\r\n$5\r\ninner\r\n:42\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Deep nested array",
			In:          array{value: []Value{array{value: []Value{array{value: []Value{integer{value: 1}}}}}}},
			Expected:    []byte("*1\r\n*1\r\n*1\r\n:1\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Array with special characters in bulk string",
			In:          array{value: []Value{bulkString{value: stringPtr("hello!@#$%")}, bulkString{value: stringPtr("test\t123")}}},
			Expected:    []byte("*2\r\n$10\r\nhello!@#$%\r\n$8\r\ntest\t123\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Array with bulk string containing CRLF",
			In:          array{value: []Value{bulkString{value: stringPtr("hello\r\nworld")}}},
			Expected:    []byte("*1\r\n$12\r\nhello\r\nworld\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Array with invalid simple string (contains CR)",
			In:          array{value: []Value{simpleString{value: "hello\rworld"}}},
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
			Name:         "Simple array with mixed types",
			In:           []byte("*3\r\n$5\r\nhello\r\n:123\r\n+PONG\r\n"),
			ExpectedRest: []byte{},
			Expected:     array{value: []Value{bulkString{value: stringPtr("hello")}, integer{value: 123}, simpleString{value: "PONG"}}},
			ShouldError:  false,
		},
		{
			Name:         "Empty array",
			In:           []byte("*0\r\n"),
			ExpectedRest: []byte{},
			Expected:     array{value: []Value{}},
			ShouldError:  false,
		},
		{
			Name:         "Nil array",
			In:           []byte(NULL_ARRAY_RESP_2),
			ExpectedRest: []byte{},
			Expected:     array{value: nil},
			ShouldError:  false,
		},
		{
			Name:         "Array with single element",
			In:           []byte("*1\r\n$1\r\nA\r\n"),
			ExpectedRest: []byte{},
			Expected:     array{value: []Value{bulkString{value: stringPtr("A")}}},
			ShouldError:  false,
		},
		{
			Name:         "Array with multiple strings",
			In:           []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"),
			ExpectedRest: []byte{},
			Expected:     array{value: []Value{bulkString{value: stringPtr("hello")}, bulkString{value: stringPtr("world")}}},
			ShouldError:  false,
		},
		{
			Name:         "Array with integers",
			In:           []byte("*3\r\n:1\r\n:-2\r\n:0\r\n"),
			ExpectedRest: []byte{},
			Expected:     array{value: []Value{integer{value: 1}, integer{value: -2}, integer{value: 0}}},
			ShouldError:  false,
		},
		{
			Name:         "Array with simple strings and errors",
			In:           []byte("*2\r\n+OK\r\n-ERR invalid\r\n"),
			ExpectedRest: []byte{},
			Expected:     array{value: []Value{simpleString{value: "OK"}, simpleError{value: "ERR invalid"}}},
			ShouldError:  false,
		},
		{
			Name:         "Array with nested array",
			In:           []byte("*2\r\n$4\r\ntest\r\n*2\r\n$5\r\ninner\r\n:42\r\n"),
			ExpectedRest: []byte{},
			Expected:     array{value: []Value{bulkString{value: stringPtr("test")}, array{value: []Value{bulkString{value: stringPtr("inner")}, integer{value: 42}}}}},
			ShouldError:  false,
		},
		{
			Name:         "Deep nested array",
			In:           []byte("*1\r\n*1\r\n*1\r\n:1\r\n"),
			ExpectedRest: []byte{},
			Expected:     array{value: []Value{array{value: []Value{array{value: []Value{integer{value: 1}}}}}}},
			ShouldError:  false,
		},
		{
			Name:         "Array with special characters",
			In:           []byte("*2\r\n$10\r\nhello!@#$%\r\n$8\r\ntest\t123\r\n"),
			ExpectedRest: []byte{},
			Expected:     array{value: []Value{bulkString{value: stringPtr("hello!@#$%")}, bulkString{value: stringPtr("test\t123")}}},
			ShouldError:  false,
		},
		{
			Name:         "Array with string containing CRLF",
			In:           []byte("*1\r\n$12\r\nhello\r\nworld\r\n"),
			ExpectedRest: []byte{},
			Expected:     array{value: []Value{bulkString{value: stringPtr("hello\r\nworld")}}},
			ShouldError:  false,
		},
		{
			Name:         "Array with remaining bytes",
			In:           []byte("*1\r\n$5\r\nhello\r\n:123\r\n"),
			ExpectedRest: []byte(":123\r\n"),
			Expected:     array{value: []Value{bulkString{value: stringPtr("hello")}}},
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
			a := array{}
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
