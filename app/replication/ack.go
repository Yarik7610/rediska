package replication

type Ack struct {
	Addr   string
	Offset int
}
