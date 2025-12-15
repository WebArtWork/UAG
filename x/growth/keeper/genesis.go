package keeper

import (
	"context"

	"uagd/x/growth/types"
)

func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Params is a VALUE (not *Params) in SDK v0.53 proto output.
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	// RegionMetrics / GrowthScores / Occupations should be VALUE slices (not []*T).
	for _, m := range genState.RegionMetrics {
		// m is value, never nil; store as-is
		if err := k.SetRegionMetric(ctx, m); err != nil {
			return err
		}
	}

	for _, s := range genState.GrowthScores {
		if err := k.SetGrowthScore(ctx, s); err != nil {
			return err
		}
	}

	for _, o := range genState.Occupations {
		// IMPORTANT: o.Occupation is already sdkmath.LegacyDec (value) in v0.53.
		// Do not parse from string.
		if err := k.SetRegionOccupation(ctx, o.RegionId, o.Occupation); err != nil {
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

	metrics := k.GetAllRegionMetrics(ctx)
	scores := k.GetAllGrowthScores(ctx)

	// Occupations might be stored as a map in keeper; export as []types.Occupation.
	occupations := k.GetAllOccupations(ctx)

	return &types.GenesisState{
		Params:        params,
		RegionMetrics: metrics,
		GrowthScores:  scores,
		Occupations:   occupations,
	}, nil
}
