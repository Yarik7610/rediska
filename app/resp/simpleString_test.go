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
			Name:        "Simple usual string",
			In:          simpleString{value: "PONG"},
			Expected:    []byte("+PONG\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Empty string",
			In:          simpleString{value: ""},
			Expected:    []byte("+\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with spaces",
			In:          simpleString{value: "HELLO WORLD"},
			Expected:    []byte("+HELLO WORLD\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with special characters",
			In:          simpleString{value: "OK!@#$%"},
			Expected:    []byte("+OK!@#$%\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Single character",
			In:          simpleString{value: "A"},
			Expected:    []byte("+A\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with numbers",
			In:          simpleString{value: "12345"},
			Expected:    []byte("+12345\r\n"),
			ShouldError: false,
		},
		{
			Name:        "String with tab",
			In:          simpleString{value: "HELLO\tWORLD"},
			Expected:    []byte("+HELLO\tWORLD\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Long string",
			In:          simpleString{value: "This is a very long string to test encoding with more than a few characters"},
			Expected:    []byte("+This is a very long string to test encoding with more than a few characters\r\n"),
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
			Name:        "Simple usual string",
			In:          []byte("+PONG\r\n"),
			Expected:    simpleString{value: "PONG"},
			ShouldError: false,
		},
		{
			Name:        "Empty string",
			In:          []byte("+\r\n"),
			Expected:    simpleString{value: ""},
			ShouldError: false,
		},
		{
			Name:        "String with spaces",
			In:          []byte("+HELLO WORLD\r\n"),
			Expected:    simpleString{value: "HELLO WORLD"},
			ShouldError: false,
		},
		{
			Name:        "String with special characters",
			In:          []byte("+OK!@#$%\r\n"),
			Expected:    simpleString{value: "OK!@#$%"},
			ShouldError: false,
		},
		{
			Name:        "Single character",
			In:          []byte("+A\r\n"),
			Expected:    simpleString{value: "A"},
			ShouldError: false,
		},
		{
			Name:        "String with numbers",
			In:          []byte("+12345\r\n"),
			Expected:    simpleString{value: "12345"},
			ShouldError: false,
		},
		{
			Name:        "String with tab",
			In:          []byte("+HELLO\tWORLD\r\n"),
			Expected:    simpleString{value: "HELLO\tWORLD"},
			ShouldError: false,
		},
		{
			Name:        "Long string",
			In:          []byte("+This is a very long string to test decoding with more than a few characters\r\n"),
			Expected:    simpleString{value: "This is a very long string to test decoding with more than a few characters"},
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
			Expected:    simpleString{value: "HELLO"},
			ShouldError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ss := simpleString{}
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
