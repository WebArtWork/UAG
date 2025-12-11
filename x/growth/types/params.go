package types

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p Params) Validate() error {
	if p.Oracle != "" {
		if _, err := sdk.AccAddressFromBech32(p.Oracle); err != nil {
			return err
		}
	}
	if p.NationalRegionId == "" {
		return fmt.Errorf("national region id required")
	}
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte("CurrentPeriod"), &p.CurrentPeriod, validateStringNonEmpty),
		paramtypes.NewParamSetPair([]byte("Oracle"), &p.Oracle, validateOracle),
		paramtypes.NewParamSetPair([]byte("NationalRegionId"), &p.NationalRegionId, validateStringNonEmpty),
	}
}

func validateStringNonEmpty(i interface{}) error {
	if v, ok := i.(string); ok {
		if v == "" {
			return fmt.Errorf("value cannot be empty")
		}
		return nil
	}
	return fmt.Errorf("invalid parameter type: %T", i)
}

func validateOracle(i interface{}) error {
	if v, ok := i.(string); ok {
		if v == "" {
			return nil
		}
		if _, err := sdk.AccAddressFromBech32(v); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("invalid parameter type: %T", i)
}

func DefaultParams() Params {
	return Params{
		CurrentPeriod:    "",
		Oracle:           "",
		NationalRegionId: NationalDefaultRegionID,
	}
}

func (m Module) Validate(c context.Context) error {
	return nil
}
