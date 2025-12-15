package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fundtypes "uagd/x/fund/types"
	"uagd/x/growth/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	types.UnimplementedQueryServer

	Schema collections.Schema

	Params collections.Item[types.Params]

	// Metrics keyed by (region_id, period)
	RegionMetrics collections.Map[collections.Pair[string, string], types.RegionMetric]

	// Scores keyed by region_id (period-specific logic can be added later)
	GrowthScores collections.Map[string, types.GrowthScore]

	// Occupations keyed by (region_id, period) -> stored as STRING (decimal) to avoid missing LegacyDec codec.
	Occupations collections.Map[collections.Pair[string, string], string]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,

		Params: collections.NewItem(
			sb,
			types.ParamsKey,
			"params",
			codec.CollValue[types.Params](cdc),
		),

		RegionMetrics: collections.NewMap(
			sb,
			types.RegionMetricKeyPrefix,
			"region_metrics",
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
			codec.CollValue[types.RegionMetric](cdc),
		),

		GrowthScores: collections.NewMap(
			sb,
			types.GrowthScoreKeyPrefix,
			"growth_scores",
			collections.StringKey,
			codec.CollValue[types.GrowthScore](cdc),
		),

		Occupations: collections.NewMap(
			sb,
			types.OccupationKeyPrefix,
			"occupations",
			collections.PairKeyCodec(collections.StringKey, collections.StringKey),
			collections.StringValue,
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

func (k Keeper) SetParams(ctx context.Context, p types.Params) error {
	return k.Params.Set(ctx, p)
}

func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	p, err := k.Params.Get(ctx)
	if err != nil {
		return types.Params{}, err
	}
	return p, nil
}

func (k Keeper) SetRegionMetric(ctx context.Context, metric types.RegionMetric) error {
	// FIX: validator expects VALUE not *VALUE
	if err := types.ValidateRegionMetric(metric); err != nil {
		return err
	}
	key := collections.Join(metric.RegionId, metric.Period)
	return k.RegionMetrics.Set(ctx, key, metric)
}

func (k Keeper) GetRegionMetric(ctx context.Context, regionID, period string) (types.RegionMetric, bool) {
	key := collections.Join(regionID, period)
	v, err := k.RegionMetrics.Get(ctx, key)
	if err != nil {
		return types.RegionMetric{}, false
	}
	return v, true
}

func (k Keeper) SetGrowthScore(ctx context.Context, score types.GrowthScore) error {
	// FIX: validator expects VALUE not *VALUE
	if err := types.ValidateGrowthScore(score); err != nil {
		return err
	}
	return k.GrowthScores.Set(ctx, score.RegionId, score)
}

func (k Keeper) GetGrowthScore(ctx context.Context, regionID string) (types.GrowthScore, bool) {
	v, err := k.GrowthScores.Get(ctx, regionID)
	if err != nil {
		return types.GrowthScore{}, false
	}
	return v, true
}

func (k Keeper) SetRegionOccupation(ctx context.Context, regionID, period string, occ sdkmath.LegacyDec) error {
	if regionID == "" || period == "" {
		return fmt.Errorf("region_id and period required")
	}
	key := collections.Join(regionID, period)
	// store as string
	return k.Occupations.Set(ctx, key, occ.String())
}

func (k Keeper) GetRegionOccupation(ctx context.Context, regionID, period string) (sdkmath.LegacyDec, bool) {
	key := collections.Join(regionID, period)
	s, err := k.Occupations.Get(ctx, key)
	if err != nil {
		return sdkmath.LegacyDec{}, false
	}
	v, err := sdkmath.LegacyNewDecFromStr(s)
	if err != nil {
		return sdkmath.LegacyDec{}, false
	}
	return v, true
}

// Used by fund module: keep stable signature.
func (k Keeper) GetEffectiveLimits(ctx context.Context, fund fundtypes.Fund) (delegationLimit sdk.Coin, payrollLimit sdk.Coin) {
	// Defaults: unlimited (0). You can wire real base caps from Params later.
	baseDelegation := sdk.NewCoin(fundtypes.BaseDenom, sdkmath.NewInt(0))
	basePayroll := sdk.NewCoin(fundtypes.BaseDenom, sdkmath.NewInt(0))

	score, found := k.GetGrowthScore(ctx, fund.RegionId)
	if !found {
		return baseDelegation, basePayroll
	}

	// Multipliers are LegacyDec VALUE already (don’t parse strings).
	if baseDelegation.Amount.IsPositive() {
		baseDelegation = sdk.NewCoin(baseDelegation.Denom, score.DelegationMultiplier.MulInt(baseDelegation.Amount).TruncateInt())
	}
	if basePayroll.Amount.IsPositive() {
		basePayroll = sdk.NewCoin(basePayroll.Denom, score.PayrollMultiplier.MulInt(basePayroll.Amount).TruncateInt())
	}

	return baseDelegation, basePayroll
}
