package types

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns param key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p Params) Validate() error {
	if p.Admin != "" {
		if _, err := sdk.AccAddressFromBech32(p.Admin); err != nil {
			return err
		}
	}
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte("Admin"), &p.Admin, validateAdmin),
	}
}

func validateAdmin(i interface{}) error {
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
	return Params{Admin: ""}
}

func (m Module) Validate(c context.Context) error {
	// no custom validation
	return nil
}
