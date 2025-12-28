package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func DefaultGenesis() *GenesisState {
	gs := &GenesisState{
		Params:         DefaultParams(),
		Metrics:        []RegionMetric{},
		Scores:         []GrowthScore{},
		OccupationList: []Occupation{},
	}

	// keep empty defaults for now (compile-first).
	// We’ll add “UA + regions” bootstrap once we confirm the real proto fields/types.

	return gs
}

func (gs GenesisState) Validate() error {
	// Minimal, safe validation only (no helpers assumed).
	for _, m := range gs.Metrics {
		if m.RegionId == "" {
			return sdkerrors.ErrInvalidRequest.Wrap("metrics: region_id cannot be empty")
		}
	}
	for _, s := range gs.Scores {
		if s.RegionId == "" {
			return sdkerrors.ErrInvalidRequest.Wrap("scores: region_id cannot be empty")
		}
	}
	for _, o := range gs.OccupationList {
		if o.RegionId == "" {
			return sdkerrors.ErrInvalidRequest.Wrap("occupation_list: region_id cannot be empty")
		}
	}
	return gs.Params.Validate()
}
