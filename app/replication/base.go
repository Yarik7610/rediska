package replication

type Base interface {
	Info() *Info
	SetMasterReplID(replID string)
	SetMasterReplOfffset(replOffset int)
}
