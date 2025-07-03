package resp

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestIntegerDecode(t *testing.T) {
// 	tests := []struct {
// 		Name        string
// 		In          []byte
// 		Expected    int
// 		ShouldError bool
// 	}{
// 		{
// 			Name:        "Positive number",
// 			In:          []byte(":123\r\n"),
// 			Expected:    123,
// 			ShouldError: false,
// 		},
// 		{
// 			Name:        "Negative number",
// 			In:          []byte(":-456\r\n"),
// 			Expected:    -456,
// 			ShouldError: false,
// 		},
// 		{
// 			Name:        "Zero",
// 			In:          []byte(":0\r\n"),
// 			Expected:    0,
// 			ShouldError: false,
// 		},
// 		{
// 			Name:        "Large positive number",
// 			In:          []byte(":123456789\r\n"),
// 			Expected:    123456789,
// 			ShouldError: false,
// 		},
// 		{
// 			Name:        "Large negative number",
// 			In:          []byte(":-987654321\r\n"),
// 			Expected:    -987654321,
// 			ShouldError: false,
// 		},
// 		{
// 			Name:        "Missing CRLF",
// 			In:          []byte(":123"),
// 			Expected:    0,
// 			ShouldError: true,
// 		},
// 		{
// 			Name:        "Invalid prefix",
// 			In:          []byte("+123\r\n"),
// 			Expected:    0,
// 			ShouldError: true,
// 		},
// 		{
// 			Name:        "Empty input",
// 			In:          []byte(""),
// 			Expected:    0,
// 			ShouldError: true,
// 		},
// 		{
// 			Name:        "Non-numeric input",
// 			In:          []byte(":abc\r\n"),
// 			Expected:    0,
// 			ShouldError: true,
// 		},
// 		{
// 			Name:        "Spaces in number",
// 			In:          []byte(":123 456\r\n"),
// 			Expected:    0,
// 			ShouldError: true,
// 		},
// 		{
// 			Name:        "Extra CRLF",
// 			In:          []byte(":123\r\n456\r\n"),
// 			Expected:    123,
// 			ShouldError: false,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.Name, func(t *testing.T) {
// 			resp := RESPController{}
// 			out, err := resp.Integer.Decode(test.In)

// 			if test.ShouldError {
// 				assert.NotNil(t, err)
// 				assert.Equal(t, test.Expected, out)
// 			} else {
// 				assert.Nil(t, err)
// 				assert.Equal(t, test.Expected, out)
// 			}
// 		})
// 	}
// }
