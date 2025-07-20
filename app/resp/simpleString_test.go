package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleStringEncode(t *testing.T) {
	tests := []struct {
		Name        string
		In          Value
		Expected    []byte
		ShouldError bool
	}{
		{
			Name:        "Empty string",
			In:          SimpleString{Value: ""},
			Expected:    []byte("+\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with spaces",
			In:          SimpleString{Value: "HELLO WORLD"},
			Expected:    []byte("+HELLO WORLD\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with special characters",
			In:          SimpleString{Value: "OK!@#$%"},
			Expected:    []byte("+OK!@#$%\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with tab",
			In:          SimpleString{Value: "HELLO\tWORLD"},
			Expected:    []byte("+HELLO\tWORLD\r\n"),
			ShouldError: false,
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

func TestSimpleStringDecode(t *testing.T) {
	tests := []struct {
		Name        string
		In          []byte
		Expected    Value
		ShouldError bool
	}{
		{
			Name:        "Empty string",
			In:          []byte("+\r\n"),
			Expected:    SimpleString{Value: ""},
			ShouldError: false,
		},
		{
			Name:        "String with spaces",
			In:          []byte("+HELLO WORLD\r\n"),
			Expected:    SimpleString{Value: "HELLO WORLD"},
			ShouldError: false,
		},
		{
			Name:        "String with special characters",
			In:          []byte("+OK!@#$%\r\n"),
			Expected:    SimpleString{Value: "OK!@#$%"},
			ShouldError: false,
		},
		{
			Name:        "String with tab",
			In:          []byte("+HELLO\tWORLD\r\n"),
			Expected:    SimpleString{Value: "HELLO\tWORLD"},
			ShouldError: false,
		},
		{
			Name:        "Missing CRLF",
			In:          []byte("+PONG"),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Invalid prefix",
			In:          []byte("-PONG\r\n"),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Empty input",
			In:          []byte(""),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "String with forbidden CRLF",
			In:          []byte("+HELLO\r\nWORLD\r\n"),
			Expected:    SimpleString{Value: "HELLO"},
			ShouldError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ss := SimpleString{}
			_, out, err := ss.Decode(test.In)

			if test.ShouldError {
				assert.NotNil(t, err)
				assert.Equal(t, test.Expected, out)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.Expected, out)
			}
		})
	}
}
