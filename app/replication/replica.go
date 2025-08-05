package replication

import (
	"net"
)

type ReplicaController interface {
	BaseController
	GetMasterConn() net.Conn
	SetMasterConn(conn net.Conn)
}

type replicaController struct {
	*baseController
	masterConn net.Conn
}

func NewReplicaController() ReplicaController {
	replicaInfo := initReplicaInfo()
	return &replicaController{
		baseController: newBaseController(replicaInfo),
	}
}

func (rc *replicaController) GetMasterConn() net.Conn {
	return rc.masterConn
}

func (rc *replicaController) SetMasterConn(conn net.Conn) {
	rc.masterConn = conn
}

func initReplicaInfo() *Info {
	return &Info{
		Role:             "slave",
		MasterReplID:     "?",
		MasterReplOffset: -1,
	}
}
