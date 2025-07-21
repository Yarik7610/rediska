package rdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeLength(t *testing.T) {
	tests := []struct {
		name              string
		buffer            []byte
		expectedLength    int
		expectedIsSpecial bool
		expectedErr       error
	}{
		{
			name:              "decodeLength: Case 0 - 6-bit length",
			buffer:            []byte{0x3A},
			expectedLength:    58,
			expectedIsSpecial: false,
			expectedErr:       nil,
		},
		{
			name:              "decodeLength: Case 1 - 14-bit length",
			buffer:            []byte{0x4A, 0xBC},
			expectedLength:    2748,
			expectedIsSpecial: false,
			expectedErr:       nil,
		},
		{
			name:              "decodeLength: Case 2 - 32-bit length",
			buffer:            []byte{0x80, 0x40, 0xE2, 0x01, 0x00},
			expectedLength:    123456,
			expectedIsSpecial: false,
			expectedErr:       nil,
		},
		{
			name:              "decodeLength: Case 3 - Special string",
			buffer:            []byte{0xC5},
			expectedLength:    0,
			expectedIsSpecial: true,
			expectedErr:       fmt.Errorf("unsupported integer string format"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dec := &decoder{b: test.buffer, pos: 0, len: len(test.buffer)}
			length, isSpecial, err := dec.decodeLength()

			assert.Equal(t, test.expectedLength, length)
			assert.Equal(t, test.expectedIsSpecial, isSpecial)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
