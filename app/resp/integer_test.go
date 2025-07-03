package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegerEncode(t *testing.T) {
	tests := []struct {
		Name        string
		In          Value
		Expected    []byte
		ShouldError bool
	}{
		{
			Name:        "Positive integer",
			In:          integer{value: 123},
			Expected:    []byte(":123\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Negative integer",
			In:          integer{value: -456},
			Expected:    []byte(":-456\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Zero",
			In:          integer{value: 0},
			Expected:    []byte(":0\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Large positive integer",
			In:          integer{value: 1234567890},
			Expected:    []byte(":1234567890\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Large negative integer",
			In:          integer{value: -987654321},
			Expected:    []byte(":-987654321\r\n"),
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

func TestIntegerDecode(t *testing.T) {
	tests := []struct {
		Name        string
		In          []byte
		Expected    Value
		ShouldError bool
	}{
		{
			Name:        "Positive integer",
			In:          []byte(":123\r\n"),
			Expected:    integer{value: 123},
			ShouldError: false,
		},
		{
			Name:        "Negative integer",
			In:          []byte(":-456\r\n"),
			Expected:    integer{value: -456},
			ShouldError: false,
		},
		{
			Name:        "Zero",
			In:          []byte(":0\r\n"),
			Expected:    integer{value: 0},
			ShouldError: false,
		},
		{
			Name:        "Large positive integer",
			In:          []byte(":1234567890\r\n"),
			Expected:    integer{value: 1234567890},
			ShouldError: false,
		},
		{
			Name:        "Large negative integer",
			In:          []byte(":-987654321\r\n"),
			Expected:    integer{value: -987654321},
			ShouldError: false,
		},
		{
			Name:        "Missing CRLF",
			In:          []byte(":123"),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Invalid prefix",
			In:          []byte("+123\r\n"),
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
			Name:        "Non-numeric value",
			In:          []byte(":abc\r\n"),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Value with forbidden CRLF",
			In:          []byte(":123\r\n456\r\n"),
			Expected:    integer{value: 123},
			ShouldError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			resp := integer{}
			out, err := resp.Decode(test.In)

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
