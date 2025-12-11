package keeper

import (
	"context"
	"fmt"

	"uagd/x/growth/types"
)

func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if genState.Params == nil {
		return fmt.Errorf("params are required")
	}
	if err := k.SetParams(ctx, *genState.Params); err != nil {
		return err
	}
	for _, m := range genState.Metrics {
		if m == nil {
			continue
		}
		if err := k.SetRegionMetric(ctx, *m); err != nil {
			return err
		}
	}
	for _, s := range genState.Scores {
		if s == nil {
			continue
		}
		if err := k.SetGrowthScore(ctx, *s); err != nil {
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
	metrics := []*types.RegionMetric{}
	metricIter, err := k.RegionMetrics.Iterate(ctx, nil)
	if err == nil {
		defer metricIter.Close()
		for ; metricIter.Valid(); metricIter.Next() {
			val, err := metricIter.Value()
			if err == nil {
				metrics = append(metrics, &val)
			}
		}
	}

	scores := []*types.GrowthScore{}
	scoreIter, err := k.GrowthScores.Iterate(ctx, nil)
	if err == nil {
		defer scoreIter.Close()
		for ; scoreIter.Valid(); scoreIter.Next() {
			val, err := scoreIter.Value()
			if err == nil {
				scores = append(scores, &val)
			}
		}
	}

	return &types.GenesisState{Params: &params, Metrics: metrics, Scores: scores}, nil
}
