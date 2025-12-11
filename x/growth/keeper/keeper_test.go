package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	fundtypes "uagd/x/fund/types"
	"uagd/x/growth/keeper"
	growthmodule "uagd/x/growth/module"
	"uagd/x/growth/types"
)

type fixture struct {
	ctx    context.Context
	keeper keeper.Keeper
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(growthmodule.AppModule{})
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_growth")).Ctx

	k := keeper.NewKeeper(storeService, encCfg.Codec)
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}
	return &fixture{ctx: ctx, keeper: k}
}

func TestSetParams(t *testing.T) {
	f := initFixture(t)
	params := types.Params{CurrentPeriod: "2025", Oracle: "", NationalRegionId: "UA"}
	if err := f.keeper.SetParams(f.ctx, params); err != nil {
		t.Fatalf("set params: %v", err)
	}
	got, err := f.keeper.GetParams(f.ctx)
	if err != nil || got.CurrentPeriod != "2025" {
		t.Fatalf("params not stored")
	}
}

func TestSetMetricComputesScore(t *testing.T) {
	f := initFixture(t)
	metric := types.RegionMetric{RegionId: "ua-test", Period: "p1", TaxIndex: "1.2", GdpIndex: "1.1", ExportsIndex: "1.0"}
	if err := f.keeper.SetRegionMetric(f.ctx, metric); err != nil {
		t.Fatalf("set metric: %v", err)
	}
	_, found := f.keeper.GetRegionMetric(f.ctx, "ua-test", "p1")
	if !found {
		t.Fatalf("metric not found")
	}
	score, found := f.keeper.GetGrowthScore(f.ctx, "ua-test", "p1")
	if !found {
		t.Fatalf("score not stored")
	}
	if score.Score == "" {
		t.Fatalf("score empty")
	}
}

func TestGetEffectiveLimits(t *testing.T) {
	f := initFixture(t)
	params := types.Params{CurrentPeriod: "2025", Oracle: "", NationalRegionId: "UA"}
	_ = f.keeper.SetParams(f.ctx, params)

	metric := types.RegionMetric{RegionId: "UA-05", Period: "2025", TaxIndex: "1.2", GdpIndex: "1.1", ExportsIndex: "1.0"}
	if err := f.keeper.SetRegionMetric(f.ctx, metric); err != nil {
		t.Fatalf("set metric: %v", err)
	}

	baseDelegation := sdk.NewCoin(fundtypes.BaseDenom, math.NewInt(100))
	basePayroll := sdk.NewCoin(fundtypes.BaseDenom, math.NewInt(50))
	fund := fundtypes.Fund{RegionId: "UA-05", Type: fundtypes.FundType_FUND_TYPE_REGION, BaseDelegationLimit: &baseDelegation, BasePayrollLimit: &basePayroll}

	del, pay := f.keeper.GetEffectiveLimits(f.ctx, fund)
	if del.Amount.Int64() <= baseDelegation.Amount.Int64() {
		t.Fatalf("delegation limit not increased")
	}
	if pay.Amount.Int64() <= basePayroll.Amount.Int64() {
		t.Fatalf("payroll limit not increased")
	}
}
