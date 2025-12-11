package keeper

import (
	"context"
	"testing"

	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	fundkeeper "uagd/x/fund/keeper"
	fundmodule "uagd/x/fund/module"
	"uagd/x/fund/types"
)

// FundKeeper creates a fund Keeper and an in-memory context for tests.
func FundKeeper(t testing.TB) (fundkeeper.Keeper, context.Context) {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(fundmodule.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_fund_test")).Ctx

	k := fundkeeper.NewKeeper(
		storeService,
		encCfg.Codec,
		addressCodec,
		nil, // bankKeeper
		nil, // stakingKeeper
	)

	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return k, ctx
}
