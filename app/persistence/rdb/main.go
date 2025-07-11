package rdb

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/redis-starter-go/app/memory"
)

type Header struct {
	name    string
	version int
}

func (h *Header) String() string {
	return fmt.Sprintf("HEADER\nName: %s, version: %d\n", h.name, h.version)
}

type Metadata struct {
	data map[string]string
}

func (m *Metadata) String() string {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintln("METADATA"))
	for key, value := range m.data {
		b.WriteString(fmt.Sprintf("Key: %s, Value: %s\n", key, value))
	}

	return b.String()
}

type Database struct {
	dbSelector      int
	keysSize        int
	expiredKeysSize int
	items           map[string]memory.Item
}

const (
	OP_EOF          = 0xFF
	OP_SELECTDB     = 0xFE
	OP_EXPIRETIME   = 0xFD
	OP_EXPIRETIMEMS = 0xFC
	OP_RESIZEDB     = 0xFB
	OP_AUX          = 0xFA
)

func IsFileExists(dir, filename string) bool {
	path := filepath.Join(dir, filename)
	fmt.Println(path)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func ReadRDB(dir, filename string) ([]byte, error) {
	path := filepath.Join(dir, filename)
	data, err := os.ReadFile(path)
	return data, err
}
