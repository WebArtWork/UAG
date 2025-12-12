package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p Params) Validate() error {
	for _, addr := range p.Registrars {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return err
		}
	}
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte("Registrars"), &p.Registrars, validateRegistrars),
	}
}

func validateRegistrars(i interface{}) error {
	addrs, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for _, addr := range addrs {
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return err
		}
	}
	return nil
}

func DefaultParams() Params {
	return Params{Registrars: []string{}}
}

func (m Module) Validate(genesis interface{}) error { return nil }
