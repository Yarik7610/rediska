package replication

import (
	"fmt"
	"strings"
)

type Info struct {
	role             string
	masterReplID     string
	masterReplOffset int
}

func NewMasterInfo() *Info {
	return &Info{
		role:             "master",
		masterReplID:     generateReplicationId(),
		masterReplOffset: 0,
	}
}

func NewReplicaInfo() *Info {
	return &Info{
		role:             "slave",
		masterReplID:     generateReplicationId(),
		masterReplOffset: 0,
	}
}

func (i *Info) String() string {
	data := []string{
		fmt.Sprintf("role:%s", i.role),
		fmt.Sprintf("master_replid:%s", i.masterReplID),
		fmt.Sprintf("master_repl_offset:%d", i.masterReplOffset),
	}
	return strings.Join(data, "\r\n") + "\r\n"
}

func generateReplicationId() string {
	return "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
}
