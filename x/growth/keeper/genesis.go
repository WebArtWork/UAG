package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"

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
	for _, o := range genState.OccupationList {
		if o == nil {
			continue
		}
		dec, err := sdkmath.LegacyNewDecFromStr(o.Occupation)
		if err != nil {
			return err
		}
		if err := k.SetRegionOccupation(ctx, o.RegionId, o.Period, dec); err != nil {
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

	occupations := []*types.Occupation{}
	occIter, err := k.Occupations.Iterate(ctx, nil)
	if err == nil {
		defer occIter.Close()
		for ; occIter.Valid(); occIter.Next() {
			key, kErr := occIter.Key()
			if kErr != nil {
				continue
			}
			val, vErr := occIter.Value()
			if vErr != nil {
				continue
			}
			occupations = append(occupations, &types.Occupation{
				RegionId:   key.K1(),
				Period:     key.K2(),
				Occupation: val.String(),
			})
		}
	}

	return &types.GenesisState{Params: &params, Metrics: metrics, Scores: scores, OccupationList: occupations}, nil
}
