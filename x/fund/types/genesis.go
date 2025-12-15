package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: Params{},
		Funds:  []Fund{},
	}
}

func (gs GenesisState) Validate() error {
	for _, f := range gs.Funds {
		if err := ValidateFund(f); err != nil {
			return err
		}
	}
	return nil
}

func ValidateFund(f Fund) error {
	if f.Address == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("fund address cannot be empty")
	}
	return nil
}
