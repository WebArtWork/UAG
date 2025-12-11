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
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"uagd/x/fund/keeper"
	fundmodule "uagd/x/fund/module"
	"uagd/x/fund/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(fundmodule.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_fund")).Ctx

	k := keeper.NewKeeper(storeService, encCfg.Codec, addressCodec, mockBankKeeper{}, mockStakingKeeper{validator: stakingValidator(addressCodec)})
	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{ctx: ctx, keeper: k, addressCodec: addressCodec}
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
