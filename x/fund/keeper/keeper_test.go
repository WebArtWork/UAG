package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"uagd/x/fund/keeper"
	fundmodule "uagd/x/fund/module"
	"uagd/x/fund/types"
)

type fixture struct {
	ctx          sdk.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
	growth       *mockGrowthKeeper
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(fundmodule.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_fund")).Ctx
	ctx = ctx.WithContext(sdk.WrapSDKContext(ctx))

	govAuthority := authtypes.NewModuleAddress(govtypes.ModuleName)
	growth := &mockGrowthKeeper{
		delegationLimit: sdk.NewInt64Coin(types.BaseDenom, 1_000_000),
		payrollLimit:    sdk.NewInt64Coin(types.BaseDenom, 1_000_000),
	}
	k := keeper.NewKeeper(storeService, encCfg.Codec, addressCodec, mockBankKeeper{}, mockStakingKeeper{validator: stakingValidator(addressCodec)}, growth, govAuthority)
	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{ctx: ctx, keeper: k, addressCodec: addressCodec, growth: growth}
}

func stakingValidator(accCodec address.Codec) stakingtypes.Validator {
	valCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	valBytes, _ := valCodec.StringToBytes("uagvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqyrlgyq")
	valAddrStr, _ := valCodec.BytesToString(valBytes)
	return stakingtypes.Validator{OperatorAddress: valAddrStr}
}

type mockBankKeeper struct{}

type mockStakingKeeper struct {
	validator stakingtypes.Validator
}

func (mockBankKeeper) SendCoins(context.Context, sdk.AccAddress, sdk.AccAddress, sdk.Coins) error {
	return nil
}

func (m mockStakingKeeper) ValidatorAddressCodec() address.Codec {
	return addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
}

func (m mockStakingKeeper) GetValidator(context.Context, sdk.ValAddress) (stakingtypes.Validator, error) {
	return m.validator, nil
}

func (mockStakingKeeper) Delegate(context.Context, sdk.AccAddress, math.Int, stakingtypes.BondStatus, stakingtypes.Validator, bool) (math.LegacyDec, error) {
	return math.LegacyZeroDec(), nil
}

type mockGrowthKeeper struct {
	delegationLimit sdk.Coin
	payrollLimit    sdk.Coin
	occupation      math.LegacyDec
	occupationFound bool
	lastFund        types.Fund
}

func (m *mockGrowthKeeper) GetEffectiveLimits(_ context.Context, fund types.Fund) (sdk.Coin, sdk.Coin) {
	m.lastFund = fund
	return m.delegationLimit, m.payrollLimit
}

func (m *mockGrowthKeeper) GetRegionOccupation(context.Context, string) (math.LegacyDec, bool) {
	return m.occupation, m.occupationFound
}

func TestSetAndGetFund(t *testing.T) {
	f := initFixture(t)
	addrStr, _ := f.addressCodec.BytesToString([]byte{1, 2, 3})
	fund := types.Fund{Address: addrStr, Type: types.FundType_FUND_TYPE_REGION, Active: true}
	if err := f.keeper.SetFund(f.ctx, fund); err != nil {
		t.Fatalf("set fund: %v", err)
	}
	got, found := f.keeper.GetFund(f.ctx, sdk.AccAddress([]byte{1, 2, 3}))
	if !found || got.Address != fund.Address {
		t.Fatalf("fund not retrieved")
	}
}

func TestValidateFundPlan(t *testing.T) {
	f := initFixture(t)
	addrStr, _ := f.addressCodec.BytesToString([]byte{9, 9, 9})
	valAddr := "uagvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqyrlgyq"
	fund := types.Fund{Address: addrStr, Type: types.FundType_FUND_TYPE_REGION, Active: true}
	_ = f.keeper.SetFund(f.ctx, fund)

	delegationAmount := sdk.NewCoin(types.BaseDenom, math.NewInt(60))
	payoutAmount := sdk.NewCoin(types.BaseDenom, math.NewInt(40))
	plan := types.FundPlan{
		FundAddress: addrStr,
		Delegations: []*types.FundDelegation{{ValidatorAddress: valAddr, Amount: &delegationAmount}},
		Payouts:     []*types.FundPayout{{RecipientAddress: addrStr, Amount: &payoutAmount}},
	}
	if err := f.keeper.ValidateFundPlan(f.ctx, plan); err != nil {
		t.Fatalf("expected plan valid: %v", err)
	}

}

func TestValidateFundPlanEffectiveLimits(t *testing.T) {
	f := initFixture(t)
	addrStr, _ := f.addressCodec.BytesToString([]byte{7, 7, 7})
	valAddr := "uagvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqyrlgyq"
	fund := types.Fund{Address: addrStr, Type: types.FundType_FUND_TYPE_REGION, Active: true}
	_ = f.keeper.SetFund(f.ctx, fund)

	delegationAmount := sdk.NewCoin(types.BaseDenom, math.NewInt(60))
	payoutAmount := sdk.NewCoin(types.BaseDenom, math.NewInt(40))
	plan := types.FundPlan{
		FundAddress: addrStr,
		Delegations: []*types.FundDelegation{{ValidatorAddress: valAddr, Amount: &delegationAmount}},
		Payouts:     []*types.FundPayout{{RecipientAddress: addrStr, Amount: &payoutAmount}},
	}

	// Tighten limits and expect failure.
	f.growth.delegationLimit = sdk.NewInt64Coin(types.BaseDenom, 50)
	if err := f.keeper.ValidateFundPlan(f.ctx, plan); err == nil {
		t.Fatalf("expected delegation limit enforcement")
	}

	// Expand limits and expect pass.
	f.growth.delegationLimit = sdk.NewInt64Coin(types.BaseDenom, 100)
	if err := f.keeper.ValidateFundPlan(f.ctx, plan); err != nil {
		t.Fatalf("expected plan valid with increased limit: %v", err)
	}
}

func TestValidateFundPlanOccupationLock(t *testing.T) {
	f := initFixture(t)
	addrStr, _ := f.addressCodec.BytesToString([]byte{8, 8, 8})
	valAddr := "uagvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqyrlgyq"
	fund := types.Fund{Address: addrStr, Type: types.FundType_FUND_TYPE_REGION, RegionId: "UA-TEST", Active: true}
	_ = f.keeper.SetFund(f.ctx, fund)

	delegationAmount := sdk.NewCoin(types.BaseDenom, math.NewInt(10))
	plan := types.FundPlan{FundAddress: addrStr, Delegations: []*types.FundDelegation{{ValidatorAddress: valAddr, Amount: &delegationAmount}}}

	f.growth.occupation = math.LegacyNewDec(60)
	f.growth.occupationFound = true
	if err := f.keeper.ValidateFundPlan(f.ctx, plan); err != types.ErrRegionLocked {
		t.Fatalf("expected occupation lock, got %v", err)
	}

	f.growth.occupation = math.LegacyNewDec(40)
	if err := f.keeper.ValidateFundPlan(f.ctx, plan); err != nil {
		t.Fatalf("expected unlocked region to pass: %v", err)
	}
}

func TestFundTypeGrowthScopes(t *testing.T) {
	f := initFixture(t)
	addrStr, _ := f.addressCodec.BytesToString([]byte{5, 5, 5})
	valAddr := "uagvaloper1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqyrlgyq"

	regional := types.Fund{Address: addrStr, Type: types.FundType_FUND_TYPE_REGION, RegionId: "UA-R1", Active: true}
	national := types.Fund{Address: addrStr + "n", Type: types.FundType_FUND_TYPE_UKRAINE, Active: true}
	projects := types.Fund{Address: addrStr + "p", Type: types.FundType_FUND_TYPE_PROJECTS, Active: true}
	_ = f.keeper.SetFund(f.ctx, regional)
	_ = f.keeper.SetFund(f.ctx, national)
	_ = f.keeper.SetFund(f.ctx, projects)

	delegationAmount := sdk.NewCoin(types.BaseDenom, math.NewInt(1))
	plan := types.FundPlan{Delegations: []*types.FundDelegation{{ValidatorAddress: valAddr, Amount: &delegationAmount}}}

	plan.FundAddress = regional.Address
	_ = f.keeper.ValidateFundPlan(f.ctx, plan)
	if f.growth.lastFund.RegionId != regional.RegionId || f.growth.lastFund.Type != regional.Type {
		t.Fatalf("expected regional limits lookup")
	}

	plan.FundAddress = national.Address
	_ = f.keeper.ValidateFundPlan(f.ctx, plan)
	if f.growth.lastFund.Type != national.Type {
		t.Fatalf("expected national scope for ukraine fund")
	}

	plan.FundAddress = projects.Address
	_ = f.keeper.ValidateFundPlan(f.ctx, plan)
	if f.growth.lastFund.Type != projects.Type {
		t.Fatalf("expected projects fund to use its own scope")
	}
}
