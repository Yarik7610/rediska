package replication

type Replica interface {
	Main
	ReadFromMaster() ([]byte, error)
}
