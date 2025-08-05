package replication

type BaseController interface {
	Info() *Info
	SetMasterReplID(replID string)
	SetMasterReplOfffset(replOffset int)
	IncrMasterReplOffset(replOffset int)
}

type baseController struct {
	info *Info
}

var _ BaseController = (*baseController)(nil)

func newBaseController(info *Info) *baseController {
	return &baseController{info: info}
}

func (bc *baseController) Info() *Info {
	return bc.info
}

func (bc *baseController) SetMasterReplID(replID string) {
	bc.info.MasterReplID = replID
}

func (bc *baseController) SetMasterReplOfffset(replOffset int) {
	bc.info.MasterReplOffset = replOffset
}

func (bc *baseController) IncrMasterReplOffset(replOffset int) {
	bc.info.MasterReplOffset += replOffset
}
