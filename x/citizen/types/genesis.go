package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default GenesisState.
func DefaultGenesis() *GenesisState {
	p := DefaultParams()
	return &GenesisState{
		Params:  &p,
		Entries: []*CitizenRegion{},
	}
}

// Validate performs basic genesis state validation.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	seen := make(map[string]struct{})
	for _, entry := range gs.Entries {
		if entry == nil {
			continue
		}
		if entry.Address == "" {
			return fmt.Errorf("address required")
		}
		if _, err := bech32Validate(entry.Address); err != nil {
			return err
		}
		if entry.RegionId == "" {
			return fmt.Errorf("region id required")
		}
		if _, exists := seen[entry.Address]; exists {
			return fmt.Errorf("duplicate address: %s", entry.Address)
		}
		seen[entry.Address] = struct{}{}
	}
	return nil
}

func bech32Validate(addr string) (string, error) {
	if addr == "" {
		return "", fmt.Errorf("address required")
	}
	if _, err := sdk.AccAddressFromBech32(addr); err != nil {
		return "", err
	}
	return addr, nil
}
