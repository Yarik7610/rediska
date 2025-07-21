package rdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringFunctions(t *testing.T) {
	tests := []struct {
		name        string
		buffer      []byte
		runTest     func(*decoder) (any, error)
		expected    any
		expectedErr error
	}{
		{
			name:   "decodeString: Regular string",
			buffer: []byte{0x5, 'h', 'e', 'l', 'l', 'o'},
			runTest: func(dec *decoder) (any, error) {
				return dec.decodeString()
			},
			expected:    "hello",
			expectedErr: nil,
		},
		{
			name:   "traverseSpecialString: Case 0 - 8-bit integer string",
			buffer: []byte{0x42},
			runTest: func(dec *decoder) (any, error) {
				return dec.traverseSpecialString(0)
			},
			expected:    66,
			expectedErr: nil,
		},
		{
			name:   "traverseSpecialString: Case 1 - 16-bit integer string",
			buffer: []byte{0x34, 0x12}, // 0x1234 in little-endian
			runTest: func(dec *decoder) (any, error) {
				return dec.traverseSpecialString(1)
			},
			expected:    4660,
			expectedErr: nil,
		},
		{
			name:   "traverseSpecialString: Case 2 - 32-bit integer string",
			buffer: []byte{0x78, 0x56, 0x34, 0x12}, // 0x12345678 in little-endian
			runTest: func(dec *decoder) (any, error) {
				return dec.traverseSpecialString(2)
			},
			expected:    305419896,
			expectedErr: nil,
		},
		{
			name:   "traverseSpecialString: Case 3 - Unsupported compressed string",
			buffer: []byte{0xFF}, // Arbitrary byte
			runTest: func(dec *decoder) (any, error) {
				return dec.traverseSpecialString(3)
			},
			expected:    0,
			expectedErr: fmt.Errorf("unsupported compressed string format"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dec := &decoder{b: test.buffer, pos: 0, len: len(test.buffer)}
			result, err := test.runTest(dec)

			assert.Equal(t, test.expected, result)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
