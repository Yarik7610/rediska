package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleErrorEncode(t *testing.T) {
	tests := []struct {
		Name        string
		In          Value
		Expected    []byte
		ShouldError bool
	}{
		{
			Name:        "Empty error",
			In:          SimpleError{Value: ""},
			Expected:    []byte("-\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Error with spaces",
			In:          SimpleError{Value: "ERR syntax error"},
			Expected:    []byte("-ERR syntax error\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Error with special characters",
			In:          SimpleError{Value: "ERR!@#$%"},
			Expected:    []byte("-ERR!@#$%\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Error with tab",
			In:          SimpleError{Value: "ERR\tinvalid"},
			Expected:    []byte("-ERR\tinvalid\r\n"),
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

func TestSimpleErrorDecode(t *testing.T) {
	tests := []struct {
		Name        string
		In          []byte
		Expected    Value
		ShouldError bool
	}{
		{
			Name:        "Empty error",
			In:          []byte("-\r\n"),
			Expected:    SimpleError{Value: ""},
			ShouldError: false,
		},
		{
			Name:        "Error with spaces",
			In:          []byte("-ERR syntax error\r\n"),
			Expected:    SimpleError{Value: "ERR syntax error"},
			ShouldError: false,
		},
		{
			Name:        "Error with special characters",
			In:          []byte("-ERR!@#$%\r\n"),
			Expected:    SimpleError{Value: "ERR!@#$%"},
			ShouldError: false,
		},
		{
			Name:        "Error with tab",
			In:          []byte("-ERR\tinvalid\r\n"),
			Expected:    SimpleError{Value: "ERR\tinvalid"},
			ShouldError: false,
		},
		{
			Name:        "Missing CRLF",
			In:          []byte("-ERR invalid"),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Invalid prefix",
			In:          []byte("+ERR invalid\r\n"),
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
			Name:        "Error with forbidden CRLF",
			In:          []byte("-ERR invalid\r\ncommand\r\n"),
			Expected:    SimpleError{Value: "ERR invalid"},
			ShouldError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			se := SimpleError{}
			_, out, err := se.Decode(test.In)

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
