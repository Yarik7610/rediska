package replication

type Replica interface {
	Base
	ReadFromMaster() ([]byte, int, error)
}
