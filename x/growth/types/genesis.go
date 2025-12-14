package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
)

// DefaultGenesis returns the default GenesisState.
func DefaultGenesis() *GenesisState {
	p := DefaultParams()
	return &GenesisState{
		Params:         &p,
		Metrics:        []*RegionMetric{},
		Scores:         []*GrowthScore{},
		OccupationList: []*Occupation{},
	}
}

// Validate performs basic genesis state validation.
func (gs GenesisState) Validate() error {
	if gs.Params == nil {
		return fmt.Errorf("params are required")
	}
	if err := gs.Params.Validate(); err != nil {
		return err
	}
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
	for _, o := range gs.OccupationList {
		if err := ValidateOccupation(o); err != nil {
			return err
		}
	}
	return nil
}

// ValidateRegionMetric ensures a metric is well-formed.
func ValidateRegionMetric(m *RegionMetric) error {
	if m == nil {
		return fmt.Errorf("metric is required")
	}
	if m.RegionId == "" {
		return fmt.Errorf("region id required")
	}
	if m.Period == "" {
		return fmt.Errorf("period required")
	}
	tax, err := sdkmath.LegacyNewDecFromStr(m.TaxIndex)
	if err != nil {
		return err
	}
	gdp, err := sdkmath.LegacyNewDecFromStr(m.GdpIndex)
	if err != nil {
		return err
	}
	exports, err := sdkmath.LegacyNewDecFromStr(m.ExportsIndex)
	if err != nil {
		return err
	}
	if tax.IsNegative() || gdp.IsNegative() || exports.IsNegative() {
		return ErrInvalidMetric
	}
	return nil
}

// ValidateGrowthScore ensures a derived score is well-formed.
func ValidateGrowthScore(s *GrowthScore) error {
	if s == nil {
		return fmt.Errorf("score is required")
	}
	if s.RegionId == "" {
		return fmt.Errorf("region id required")
	}
	if s.Period == "" {
		return fmt.Errorf("period required")
	}
	if _, err := sdkmath.LegacyNewDecFromStr(s.Score); err != nil {
		return err
	}
	if _, err := sdkmath.LegacyNewDecFromStr(s.DelegationMultiplier); err != nil {
		return err
	}
	if _, err := sdkmath.LegacyNewDecFromStr(s.PayrollMultiplier); err != nil {
		return err
	}
	return nil
}

// ValidateOccupation ensures an occupation entry is well-formed.
func ValidateOccupation(o *Occupation) error {
	if o == nil {
		return fmt.Errorf("occupation is required")
	}
	if o.RegionId == "" {
		return fmt.Errorf("region id required")
	}
	if o.Period == "" {
		return fmt.Errorf("period required")
	}
	dec, err := sdkmath.LegacyNewDecFromStr(o.Occupation)
	if err != nil {
		return err
	}
	if dec.IsNegative() {
		return ErrInvalidMetric
	}
	return nil
}
