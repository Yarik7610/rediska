package resp

import (
	"os"
	"testing"
)

func stringPtr(s string) *string {
	return &s
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
