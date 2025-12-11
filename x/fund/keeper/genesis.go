package keeper

import (
	"context"
	"fmt"

	"uagd/x/fund/types"
)

func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if genState.Params == nil {
		return fmt.Errorf("params are required")
	}
	if err := k.Params.Set(ctx, *genState.Params); err != nil {
		return err
	}
	for _, f := range genState.Funds {
		if f == nil {
			return fmt.Errorf("fund entry cannot be nil")
		}
		if err := k.SetFund(ctx, *f); err != nil {
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
	fundPtrs := make([]*types.Fund, len(funds))
	for i := range funds {
		fundPtrs[i] = &funds[i]
	}
	return &types.GenesisState{Params: &params, Funds: fundPtrs}, nil
}
