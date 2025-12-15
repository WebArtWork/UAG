package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DefaultParams() Params {
	return Params{Authority: ""}
}

func (p Params) Validate() error {
	if p.Authority == "" {
		return nil
	}
	if _, err := sdk.AccAddressFromBech32(p.Authority); err != nil {
		return fmt.Errorf("invalid authority address: %w", err)
	}
	return nil
}
