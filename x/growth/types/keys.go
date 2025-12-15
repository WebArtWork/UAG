package types

// ModuleName defines the module name
const ModuleName = "growth"

// StoreKey defines the primary module store key
const StoreKey = ModuleName

// RouterKey is the message route for this module
const RouterKey = ModuleName

var (
	// ParamsKey is the key for params in collections
	ParamsKey = []byte{0x00}

	// Collections key prefixes (used by keeper)
	RegionMetricKeyPrefix = []byte{0x01}
	GrowthScoreKeyPrefix  = []byte{0x02}
	OccupationKeyPrefix   = []byte{0x03}
)
