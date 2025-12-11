package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Params struct {
	Admin string `json:"admin" yaml:"admin"`
}

func DefaultParams() Params {
	return Params{Admin: ""}
}

func (p Params) Validate() error {
	if p.Admin == "" {
		return nil
	}
	if _, err := sdk.AccAddressFromBech32(p.Admin); err != nil {
		return fmt.Errorf("invalid admin address: %w", err)
	}
	return nil
}
