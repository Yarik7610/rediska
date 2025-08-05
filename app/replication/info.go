package replication

import (
	"strconv"
	"strings"
)

type Info struct {
	Role             string
	MasterReplID     string
	MasterReplOffset int
}

func (i *Info) String() string {
	data := []string{
		"role:" + i.Role,
		"master_replid:" + i.MasterReplID,
		"master_repl_offset:" + strconv.Itoa(i.MasterReplOffset),
	}
	return strings.Join(data, "\r\n") + "\r\n"
}

func generateReplicationId() string {
	return "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
}
