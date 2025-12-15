package keeper

import (
	"context"

	"uagd/x/fund/types"
)

func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Params is a value type in SDK v0.53 proto output (not *Params).
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	// Funds is []Fund (not []*Fund).
	for _, f := range genState.Funds {
		if err := k.SetFund(ctx, f); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	funds := k.GetAllFunds(ctx)

	return &types.GenesisState{
		Params: params,
		Funds:  funds,
	}, nil
}
