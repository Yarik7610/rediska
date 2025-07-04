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
			Name:        "Positive Integer",
			In:          Integer{Value: 123},
			Expected:    []byte(":123\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Negative Integer",
			In:          Integer{Value: -456},
			Expected:    []byte(":-456\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Zero",
			In:          Integer{Value: 0},
			Expected:    []byte(":0\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Large positive Integer",
			In:          Integer{Value: 1234567890},
			Expected:    []byte(":1234567890\r\n"),
			ShouldError: false,
		},
		{
			Name:        "Large negative Integer",
			In:          Integer{Value: -987654321},
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
			Name:        "Positive Integer",
			In:          []byte(":123\r\n"),
			Expected:    Integer{Value: 123},
			ShouldError: false,
		},
		{
			Name:        "Negative Integer",
			In:          []byte(":-456\r\n"),
			Expected:    Integer{Value: -456},
			ShouldError: false,
		},
		{
			Name:        "Zero",
			In:          []byte(":0\r\n"),
			Expected:    Integer{Value: 0},
			ShouldError: false,
		},
		{
			Name:        "Large positive Integer",
			In:          []byte(":1234567890\r\n"),
			Expected:    Integer{Value: 1234567890},
			ShouldError: false,
		},
		{
			Name:        "Large negative Integer",
			In:          []byte(":-987654321\r\n"),
			Expected:    Integer{Value: -987654321},
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
			Name:        "Non-numeric Value",
			In:          []byte(":abc\r\n"),
			Expected:    nil,
			ShouldError: true,
		},
		{
			Name:        "Value with forbidden CRLF",
			In:          []byte(":123\r\n456\r\n"),
			Expected:    Integer{Value: 123},
			ShouldError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			i := Integer{}
			_, out, err := i.Decode(test.In)

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
