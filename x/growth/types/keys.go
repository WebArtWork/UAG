package types

import "cosmossdk.io/collections"

const (
	ModuleName              = "growth"
	StoreKey                = ModuleName
	RouterKey               = ModuleName
	NationalDefaultRegionID = "UA"
)

var (
	ParamsKey           = collections.NewPrefix("p_growth")
	MetricKeyPrefix     = collections.NewPrefix("metric")
	ScoreKeyPrefix      = collections.NewPrefix("score")
	OccupationKeyPrefix = collections.NewPrefix("occup")
)
