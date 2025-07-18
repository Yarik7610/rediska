package replication

type Replica interface {
	Base
	ConnectToMaster()
	ReadFromMaster() ([]byte, int, error)
}
