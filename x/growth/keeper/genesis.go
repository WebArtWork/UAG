package keeper

import (
	"context"

	"uag/x/growth/types"
)

func (k Keeper) InitGenesis(ctx context.Context, gs types.GenesisState) error {
	// keep empty for now (compile-first).
	// later we’ll store gs.Metrics / gs.Scores / gs.OccupationList into KV collections.
	return nil
}

func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	// keep default for now (compile-first).
	return types.DefaultGenesis(), nil
}
