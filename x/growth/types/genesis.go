package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:  Params{},
		Metrics: []RegionMetric{},
		Scores:  []GrowthScore{},
	}
}

func (gs GenesisState) Validate() error {
	for _, m := range gs.Metrics {
		if err := ValidateRegionMetric(m); err != nil {
			return err
		}
	}
	for _, s := range gs.Scores {
		if err := ValidateGrowthScore(s); err != nil {
			return err
		}
	}
	return nil
}

func ValidateRegionMetric(m RegionMetric) error {
	if m.RegionId == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("region_id cannot be empty")
	}
	return nil
}

func ValidateGrowthScore(s GrowthScore) error {
	if s.RegionId == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("region_id cannot be empty")
	}
	return nil
}
