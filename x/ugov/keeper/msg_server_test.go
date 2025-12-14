package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"

	fundtypes "uagd/x/fund/types"
	ugovkeeper "uagd/x/ugov/keeper"
	"uagd/x/ugov/types"
)

type stubFundKeeper struct {
	called   int
	lastPlan fundtypes.FundPlan
	lastAuth sdk.AccAddress
	err      error
}

func (s *stubFundKeeper) ExecuteFundPlan(ctx context.Context, plan fundtypes.FundPlan, authority sdk.AccAddress) error {
	s.called++
	s.lastPlan = plan
	s.lastAuth = authority
	return s.err
}

func setupKeeper(t *testing.T) (ugovkeeper.Keeper, sdk.Context, *stubFundKeeper) {
	t.Helper()
	encCfg := moduletestutil.MakeTestEncodingConfig()
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	tkey := storetypes.NewTransientStoreKey("transient_ugov_test")
	ctx := testutil.DefaultContextWithDB(t, storeKey, tkey).Ctx.WithLogger(log.NewNopLogger())
	ctx = ctx.WithContext(sdk.WrapSDKContext(ctx))

	pk := paramskeeper.NewKeeper(encCfg.Codec, encCfg.Amino, storeKey, tkey)
	subspace := pk.Subspace(types.ModuleName)

	fund := &stubFundKeeper{}
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	k := ugovkeeper.NewKeeper(encCfg.Codec, runtime.NewKVStoreService(storeKey), subspace, fund, authority)

	// set national president to satisfy plan creation
	president := types.President{RoleType: types.PRESIDENT_TYPE_NATIONAL, Address: sampleAcc(), Active: true}
	k.SetPresident(ctx, president)

	return k, ctx, fund
}

func sampleAcc() string {
	pk := authtypes.NewModuleAddress("sampler")
	return sdk.AccAddress(pk).String()
}

func TestExecuteFundPlanRequiresSubmittedStatus(t *testing.T) {
	k, ctx, _ := setupKeeper(t)
	msgServer := ugovkeeper.NewMsgServerImpl(k)

	planJSON, _ := fundtypes.GetFundPlanCodec().MarshalJSON(&fundtypes.FundPlan{FundAddress: sampleAcc()})
	id, err := k.CreatePlan(ctx, sampleAcc(), sampleAcc(), "title", "desc", types.PRESIDENT_TYPE_NATIONAL, "", planJSON)
	if err != nil {
		t.Fatalf("create plan: %v", err)
	}

	_, err = msgServer.ExecuteFundPlan(ctx, &types.MsgExecuteFundPlan{Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(), PlanId: id})
	if err == nil {
		t.Fatalf("expected error for non-submitted plan")
	}
}

func TestExecuteFundPlanHappyPath(t *testing.T) {
	k, ctx, fund := setupKeeper(t)
	ctx = ctx.WithBlockHeight(10)
	msgServer := ugovkeeper.NewMsgServerImpl(k)

	planJSON, _ := fundtypes.GetFundPlanCodec().MarshalJSON(&fundtypes.FundPlan{FundAddress: sampleAcc()})
	id, err := k.CreatePlan(ctx, sampleAcc(), sampleAcc(), "title", "desc", types.PRESIDENT_TYPE_NATIONAL, "", planJSON)
	if err != nil {
		t.Fatalf("create plan: %v", err)
	}
	if err := k.MarkSubmitted(ctx, id, 1); err != nil {
		t.Fatalf("mark submitted: %v", err)
	}

	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	res, err := msgServer.ExecuteFundPlan(ctx, &types.MsgExecuteFundPlan{Authority: authority.String(), PlanId: id})
	if err != nil {
		t.Fatalf("execute plan: %v", err)
	}

	if fund.called != 1 {
		t.Fatalf("expected fund keeper to be called")
	}
	if !fund.lastAuth.Equals(authority) {
		t.Fatalf("unexpected authority used for execution")
	}
	if res == nil {
		t.Fatalf("expected response")
	}

	stored, ok := k.GetPlan(ctx, id)
	if !ok {
		t.Fatalf("plan missing after execution")
	}
	if stored.Status != types.PLAN_STATUS_EXECUTED {
		t.Fatalf("expected executed status, got %d", stored.Status)
	}
	if stored.ExecutedAtHeight != ctx.BlockHeight() {
		t.Fatalf("expected executed height %d got %d", ctx.BlockHeight(), stored.ExecutedAtHeight)
	}
}
