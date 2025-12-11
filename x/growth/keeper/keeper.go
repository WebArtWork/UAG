package keeper

import (
	"context"

	"cosmossdk.io/collections"
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
	types.UnimplementedQueryServer

	Schema        collections.Schema
	Params        collections.Item[types.Params]
	RegionMetrics collections.Map[collections.Pair[string, string], types.RegionMetric]
	GrowthScores  collections.Map[collections.Pair[string, string], types.GrowthScore]
}

func NewKeeper(storeService corestore.KVStoreService, cdc codec.Codec) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		storeService:  storeService,
		cdc:           cdc,
		Params:        collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		RegionMetrics: collections.NewMap(sb, types.MetricKeyPrefix, "metrics", collections.PairKeyCodec(collections.StringKey, collections.StringKey), codec.CollValue[types.RegionMetric](cdc)),
		GrowthScores:  collections.NewMap(sb, types.ScoreKeyPrefix, "scores", collections.PairKeyCodec(collections.StringKey, collections.StringKey), codec.CollValue[types.GrowthScore](cdc)),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}
	return k.Params.Set(ctx, params)
}

func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	return k.Params.Get(ctx)
}

func (k Keeper) SetRegionMetric(ctx context.Context, metric types.RegionMetric) error {
	if err := types.ValidateRegionMetric(&metric); err != nil {
		return err
	}
	if err := k.RegionMetrics.Set(ctx, collections.Join(metric.RegionId, metric.Period), metric); err != nil {
		return err
	}
	score := k.ComputeGrowthScore(metric)
	return k.SetGrowthScore(ctx, score)
}

func (k Keeper) GetRegionMetric(ctx context.Context, regionID, period string) (types.RegionMetric, bool) {
	metric, err := k.RegionMetrics.Get(ctx, collections.Join(regionID, period))
	if err != nil {
		return types.RegionMetric{}, false
	}
	return metric, true
}

func (k Keeper) GetRegionMetricsForRegion(ctx context.Context, regionID string) []types.RegionMetric {
	var metrics []types.RegionMetric
	iterator, err := k.RegionMetrics.Iterate(ctx, nil)
	if err != nil {
		return metrics
	}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key, err := iterator.Key()
		if err != nil {
			continue
		}
		if key.K1() != regionID {
			continue
		}
		value, err := iterator.Value()
		if err == nil {
			metrics = append(metrics, value)
		}
	}
	return metrics
}

func (k Keeper) GetRegionMetricsForPeriod(ctx context.Context, period string) []types.RegionMetric {
	var metrics []types.RegionMetric
	iterator, err := k.RegionMetrics.Iterate(ctx, nil)
	if err != nil {
		return metrics
	}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key, err := iterator.Key()
		if err != nil {
			continue
		}
		if key.K2() != period {
			continue
		}
		value, err := iterator.Value()
		if err == nil {
			metrics = append(metrics, value)
		}
	}
	return metrics
}

func (k Keeper) ComputeGrowthScore(metric types.RegionMetric) types.GrowthScore {
	tax := mustDec(metric.TaxIndex)
	gdp := mustDec(metric.GdpIndex)
	exports := mustDec(metric.ExportsIndex)

	score := tax.Add(gdp).Add(exports).QuoInt64(3)
	score = clampDec(score, sdkmath.LegacyZeroDec(), sdkmath.LegacyNewDec(2))
	base := sdkmath.LegacyMustNewDecFromStr("0.5")
	delegationMultiplier := clampDec(base.Add(score), sdkmath.LegacyZeroDec(), sdkmath.LegacyNewDec(3))
	payrollMultiplier := clampDec(base.Add(score), sdkmath.LegacyZeroDec(), sdkmath.LegacyNewDec(3))

	return types.GrowthScore{
		RegionId:             metric.RegionId,
		Period:               metric.Period,
		Score:                score.String(),
		DelegationMultiplier: delegationMultiplier.String(),
		PayrollMultiplier:    payrollMultiplier.String(),
	}
}

func (k Keeper) SetGrowthScore(ctx context.Context, score types.GrowthScore) error {
	if err := types.ValidateGrowthScore(&score); err != nil {
		return err
	}
	return k.GrowthScores.Set(ctx, collections.Join(score.RegionId, score.Period), score)
}

func (k Keeper) GetGrowthScore(ctx context.Context, regionID, period string) (types.GrowthScore, bool) {
	score, err := k.GrowthScores.Get(ctx, collections.Join(regionID, period))
	if err != nil {
		return types.GrowthScore{}, false
	}
	return score, true
}

// GetEffectiveLimits returns base fund limits multiplied by growth multipliers for the current period.
// Intended integration: x/fund keeper should call this helper when validating fund plans instead of using
// static base limits directly.
func (k Keeper) GetEffectiveLimits(ctx context.Context, fund fundtypes.Fund) (delegationLimit sdk.Coin, payrollLimit sdk.Coin) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return baseDelegationLimitOrZero(fund), basePayrollLimitOrZero(fund)
	}

	regionID := params.NationalRegionId
	if fund.Type == fundtypes.FundType_FUND_TYPE_REGION && fund.RegionId != "" {
		regionID = fund.RegionId
	}

	delegationMultiplier := sdkmath.LegacyOneDec()
	payrollMultiplier := sdkmath.LegacyOneDec()
	if score, found := k.GetGrowthScore(ctx, regionID, params.CurrentPeriod); found {
		if dm, err := sdkmath.LegacyNewDecFromStr(score.DelegationMultiplier); err == nil {
			delegationMultiplier = dm
		}
		if pm, err := sdkmath.LegacyNewDecFromStr(score.PayrollMultiplier); err == nil {
			payrollMultiplier = pm
		}
	}

	baseDelegation := baseDelegationLimitOrZero(fund)
	if baseDelegation.Denom == "" {
		baseDelegation.Denom = fundtypes.BaseDenom
	}
	basePayroll := basePayrollLimitOrZero(fund)
	if basePayroll.Denom == "" {
		basePayroll.Denom = fundtypes.BaseDenom
	}

	delegationAmount := sdkmath.LegacyNewDecFromInt(baseDelegation.Amount).Mul(delegationMultiplier).TruncateInt()
	payrollAmount := sdkmath.LegacyNewDecFromInt(basePayroll.Amount).Mul(payrollMultiplier).TruncateInt()

	return sdk.NewCoin(baseDelegation.Denom, delegationAmount), sdk.NewCoin(basePayroll.Denom, payrollAmount)
}

func mustDec(v string) sdkmath.LegacyDec {
	dec, err := sdkmath.LegacyNewDecFromStr(v)
	if err != nil {
		panic(err)
	}
	return dec
}

func clampDec(value, min, max sdkmath.LegacyDec) sdkmath.LegacyDec {
	if value.LT(min) {
		return min
	}
	if value.GT(max) {
		return max
	}
	return value
}

// Base limit helpers when fund module has nil coins.
func baseDelegationLimitOrZero(f fundtypes.Fund) sdk.Coin {
	if f.BaseDelegationLimit == nil {
		return sdk.NewCoin(fundtypes.BaseDenom, sdkmath.ZeroInt())
	}
	return *f.BaseDelegationLimit
}

func basePayrollLimitOrZero(f fundtypes.Fund) sdk.Coin {
	if f.BasePayrollLimit == nil {
		return sdk.NewCoin(fundtypes.BaseDenom, sdkmath.ZeroInt())
	}
	return *f.BasePayrollLimit
}
