package types

import "cosmossdk.io/collections"

const (
	ModuleName = "citizen"
	StoreKey   = ModuleName
)

var (
	ParamsKey             = collections.NewPrefix("p_citizen")
	RegionByAddressPrefix = []byte{0x01}
	AddressByRegionPrefix = []byte{0x02}
)
