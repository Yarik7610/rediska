package replication

type Replica interface {
	Main
	ReadValueFromMaster() ([]byte, error)
}
