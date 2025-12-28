package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:  Params{},
		Entries: []CitizenRegion{},
	}
}

func (gs GenesisState) Validate() error {
	// Params is a VALUE, no nil checks, no Validate()
	for _, entry := range gs.Entries {
		if err := ValidateCitizenRegion(entry); err != nil {
			return err
		}
	}
	return nil
}

func ValidateCitizenRegion(entry CitizenRegion) error {
	if entry.RegionId == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("region_id cannot be empty")
	}
	if entry.Address == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("address cannot be empty")
	}
	return nil
}
