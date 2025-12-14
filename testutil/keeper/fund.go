package keeper

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

	_ "uagd/app" // ensure bech32 prefix init runs in tests
	fundkeeper "uagd/x/fund/keeper"
	fundmodule "uagd/x/fund/module"
	"uagd/x/fund/types"
)

type noOpBankKeeper struct{}

func (noOpBankKeeper) SendCoins(context.Context, sdk.AccAddress, sdk.AccAddress, sdk.Coins) error {
	return nil
}

type noOpStakingKeeper struct{ addrCodec address.Codec }

func (k noOpStakingKeeper) ValidatorAddressCodec() address.Codec { return k.addrCodec }
func (noOpStakingKeeper) GetValidator(context.Context, sdk.ValAddress) (stakingtypes.Validator, error) {
	return stakingtypes.Validator{}, nil
}
func (noOpStakingKeeper) Delegate(context.Context, sdk.AccAddress, math.Int, stakingtypes.BondStatus, stakingtypes.Validator, bool) (math.LegacyDec, error) {
	return math.LegacyZeroDec(), nil
}

type noOpGrowthKeeper struct{}

func (noOpGrowthKeeper) GetEffectiveLimits(context.Context, types.Fund) (sdk.Coin, sdk.Coin) {
	return sdk.NewCoin(types.BaseDenom, math.ZeroInt()), sdk.NewCoin(types.BaseDenom, math.ZeroInt())
}

func (noOpGrowthKeeper) GetRegionOccupation(context.Context, string) (math.LegacyDec, bool) {
	return math.LegacyZeroDec(), false
}

// FundKeeper creates a fund Keeper and an in-memory context for tests.
func FundKeeper(t testing.TB) (fundkeeper.Keeper, context.Context) {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(fundmodule.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	bankKeeper := noOpBankKeeper{}
	stakingKeeper := noOpStakingKeeper{addrCodec: addressCodec}
	growthKeeper := noOpGrowthKeeper{}
	govAuthority := authtypes.NewModuleAddress(govtypes.ModuleName)

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_fund_test")).Ctx

	k := fundkeeper.NewKeeper(
		storeService,
		encCfg.Codec,
		addressCodec,
		bankKeeper,
		stakingKeeper,
		growthKeeper,
		govAuthority,
	)

	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return k, ctx
}
