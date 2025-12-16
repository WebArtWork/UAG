package types

const (
	ModuleName = "ugov"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	ParamsKey     = []byte{0x00}
	PlanKeyPrefix = []byte{0x01}
	NextPlanIDKey = []byte{0x02}
)
