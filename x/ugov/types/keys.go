package types

const (
	ModuleName = "ugov"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	PresidentKeyPrefix = []byte{0x01}
	PlanKeyPrefix      = []byte{0x02}
	PlanSeqKey         = []byte{0x03}
)
