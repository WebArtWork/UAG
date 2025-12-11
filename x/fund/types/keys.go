package types

import "cosmossdk.io/collections"

const (
	ModuleName    = "fund"
	StoreKey      = ModuleName
	RouterKey     = ModuleName
	GovModuleName = "gov"
)

var (
	FundKeyPrefix = collections.NewPrefix("fund")
	ParamsKey     = collections.NewPrefix("p_fund")
)

const BaseDenom = "uuag"
