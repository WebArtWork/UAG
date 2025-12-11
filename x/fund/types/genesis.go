package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default GenesisState.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: func() *Params { p := DefaultParams(); return &p }(),
		Funds:  []*Fund{},
	}
}

// Validate performs basic genesis state validation.
func (gs GenesisState) Validate() error {
	addrSet := map[string]struct{}{}
	for _, f := range gs.Funds {
		if err := ValidateFund(f); err != nil {
			return err
		}
		if f == nil {
			return fmt.Errorf("fund entry cannot be nil")
		}
		if _, exists := addrSet[f.Address]; exists {
			return fmt.Errorf("duplicate fund address %s", f.Address)
		}
		addrSet[f.Address] = struct{}{}
	}
	if gs.Params == nil {
		return fmt.Errorf("params are required")
	}
	return gs.Params.Validate()
}

// ValidateFund performs basic validation for a fund.
func ValidateFund(f *Fund) error {
	if f == nil {
		return fmt.Errorf("fund cannot be nil")
	}
	if f.Address == "" {
		return fmt.Errorf("fund address required")
	}
	if _, err := sdk.AccAddressFromBech32(f.Address); err != nil {
		return err
	}
	if f.President != "" {
		if _, err := sdk.AccAddressFromBech32(f.President); err != nil {
			return err
		}
	}
	switch f.Type {
	case FundType_FUND_TYPE_REGION, FundType_FUND_TYPE_UKRAINE, FundType_FUND_TYPE_PROJECTS:
	default:
		return fmt.Errorf("invalid fund type")
	}
	if f.BaseDelegationLimit != nil {
		if f.BaseDelegationLimit.Denom != "" && f.BaseDelegationLimit.Denom != BaseDenom {
			return fmt.Errorf("invalid delegation denom: %s", f.BaseDelegationLimit.Denom)
		}
	}
	if f.BasePayrollLimit != nil {
		if f.BasePayrollLimit.Denom != "" && f.BasePayrollLimit.Denom != BaseDenom {
			return fmt.Errorf("invalid payroll denom: %s", f.BasePayrollLimit.Denom)
		}
	}
	return nil
}
