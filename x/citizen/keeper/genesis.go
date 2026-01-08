package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"uag/x/citizen/types"
)

func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Params is a VALUE (not *Params) in SDK v0.53 proto output.
	if err := k.SetParams(ctx, genState.Params); err != nil {
		return err
	}

	// Entries is []CitizenRegion (not []*CitizenRegion).
	for _, entry := range genState.Entries {
		addr, err := sdk.AccAddressFromBech32(entry.Address)
		if err != nil {
			return err
		}
		if err := k.SetRegion(ctx, addr, entry.RegionId); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	params, err := k.ParamsStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	entries := make([]types.CitizenRegion, 0)

	iterator, err := k.RegionByAddressStore.Iterate(ctx, nil)
	if err == nil {
		defer iterator.Close()

		for ; iterator.Valid(); iterator.Next() {
			addrBytes, err := iterator.Key()
			if err != nil {
				continue
			}
			region, err := iterator.Value()
			if err != nil {
				continue
			}
			addrStr, err := k.addressCodec.BytesToString(addrBytes)
			if err != nil {
				continue
			}
			entries = append(entries, types.CitizenRegion{Address: addrStr, RegionId: region})
		}
	}

	return &types.GenesisState{
		Params:  params,
		Entries: entries,
	}, nil
}
