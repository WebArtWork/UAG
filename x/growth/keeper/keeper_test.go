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

func TestRegionOccupation(t *testing.T) {
	f := initFixture(t)
	params := types.Params{CurrentPeriod: "2025", Oracle: "", NationalRegionId: "UA"}
	_ = f.keeper.SetParams(f.ctx, params)

	occ := math.LegacyNewDec(55)
	if err := f.keeper.SetRegionOccupation(f.ctx, "UA-05", "2025", occ); err != nil {
		t.Fatalf("set occupation: %v", err)
	}

	got, found := f.keeper.GetRegionOccupation(f.ctx, "UA-05")
	if !found || !got.Equal(occ) {
		t.Fatalf("expected occupation to be stored")
	}
}

func TestGenesisOccupationsRoundTrip(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(growthmodule.AppModule{})

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_growth_genesis_1")).Ctx

	k := keeper.NewKeeper(storeService, encCfg.Codec)
	params := types.Params{CurrentPeriod: "2025", Oracle: "", NationalRegionId: "UA"}
	if err := k.SetParams(ctx, params); err != nil {
		t.Fatalf("set params: %v", err)
	}

	regionOcc := math.LegacyMustNewDecFromStr("0.55")
	nationalOcc := math.LegacyMustNewDecFromStr("0.75")
	if err := k.SetRegionOccupation(ctx, "UA-05", "2025", regionOcc); err != nil {
		t.Fatalf("set region occupation: %v", err)
	}
	if err := k.SetRegionOccupation(ctx, params.NationalRegionId, "2025", nationalOcc); err != nil {
		t.Fatalf("set national occupation: %v", err)
	}

	gen, err := k.ExportGenesis(ctx)
	if err != nil {
		t.Fatalf("export genesis: %v", err)
	}

	storeKey2 := storetypes.NewKVStoreKey(types.StoreKey)
	storeService2 := runtime.NewKVStoreService(storeKey2)
	ctx2 := testutil.DefaultContextWithDB(t, storeKey2, storetypes.NewTransientStoreKey("transient_growth_genesis_2")).Ctx
	k2 := keeper.NewKeeper(storeService2, encCfg.Codec)

	if err := k2.InitGenesis(ctx2, *gen); err != nil {
		t.Fatalf("init genesis: %v", err)
	}

	gotRegion, foundRegion := k2.GetRegionOccupation(ctx2, "UA-05")
	if !foundRegion || !gotRegion.Equal(regionOcc) {
		t.Fatalf("expected region occupation to round-trip")
	}
	gotNational, foundNational := k2.GetRegionOccupation(ctx2, params.NationalRegionId)
	if !foundNational || !gotNational.Equal(nationalOcc) {
		t.Fatalf("expected national occupation to round-trip")
	}
}
