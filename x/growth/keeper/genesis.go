package keeper

import (
	"context"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"

	"uagd/x/growth/types"
)

func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.SetParams(ctx, genState.Params); err != nil {
		return err
	}

	// Store metrics
	for _, m := range genState.Metrics {
		if err := k.SetRegionMetric(ctx, m); err != nil {
			return err
		}
	}

	// Store scores
	for _, s := range genState.Scores {
		if err := k.SetGrowthScore(ctx, s); err != nil {
			return err
		}
	}

	// Store occupations
	for _, o := range genState.OccupationList {
		if err := k.SetRegionOccupation(ctx, o.RegionId, o.Period, o.Occupation); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	metrics := make([]types.RegionMetric, 0)
	{
		it, err := k.RegionMetrics.Iterate(ctx, nil)
		if err == nil {
			defer it.Close()
			for ; it.Valid(); it.Next() {
				v, err := it.Value()
				if err == nil {
					metrics = append(metrics, v)
				}
			}
		}
	}

	scores := make([]types.GrowthScore, 0)
	{
		it, err := k.GrowthScores.Iterate(ctx, nil)
		if err == nil {
			defer it.Close()
			for ; it.Valid(); it.Next() {
				v, err := it.Value()
				if err == nil {
					scores = append(scores, v)
				}
			}
		}
	}

	occupationList := make([]types.Occupation, 0)
	{
		it, err := k.Occupations.Iterate(ctx, nil)
		if err == nil {
			defer it.Close()
			for ; it.Valid(); it.Next() {
				key, err := it.Key()
				if err != nil {
					continue
				}
				valStr, err := it.Value()
				if err != nil {
					continue
				}
				dec, err := sdkmath.LegacyNewDecFromStr(valStr)
				if err != nil {
					continue
				}

				occupationList = append(occupationList, types.Occupation{
					RegionId:   key.K1(),
					Period:     key.K2(),
					Occupation: dec,
				})
			}
		}
	}

	_ = collections.Pair[string, string]{}

	return &types.GenesisState{
		Params:         params,
		Metrics:        metrics,
		Scores:         scores,
		OccupationList: occupationList,
	}, nil
}
