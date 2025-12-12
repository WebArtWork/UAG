package keeper_test

import (
	"context"
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	_ "uagd/app"
	"uagd/x/citizen/keeper"
	citizenmodule "uagd/x/citizen/module"
	"uagd/x/citizen/types"
)

type fixture struct {
	ctx       context.Context
	keeper    keeper.Keeper
	msgServer types.MsgServer
	registrar sdk.AccAddress
	nonMember sdk.AccAddress
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(citizenmodule.AppModule{})
	addressCodec := address.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_citizen")).Ctx

	k := keeper.NewKeeper(storeService, encCfg.Codec, addressCodec)

	registrarKey := secp256k1.GenPrivKey()
	registrar := sdk.AccAddress(registrarKey.PubKey().Address())
	nonMember := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	if err := k.SetParams(ctx, types.Params{Registrars: []string{registrar.String()}}); err != nil {
		t.Fatalf("set params: %v", err)
	}

	return &fixture{ctx: ctx, keeper: k, msgServer: keeper.NewMsgServerImpl(k), registrar: registrar, nonMember: nonMember}
}

func TestSetCitizenRegionUpdatesIndex(t *testing.T) {
	f := initFixture(t)
	target := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	if _, err := f.msgServer.SetCitizenRegion(f.ctx, types.NewMsgSetCitizenRegion(f.registrar.String(), target.String(), "ua-1")); err != nil {
		t.Fatalf("set region: %v", err)
	}
	region, found := f.keeper.GetRegion(f.ctx, target)
	if !found || region != "ua-1" {
		t.Fatalf("region not stored")
	}
	addresses := f.keeper.GetAddressesByRegion(f.ctx, "ua-1")
	if len(addresses) != 1 || addresses[0] != target.String() {
		t.Fatalf("address not indexed")
	}

	if _, err := f.msgServer.SetCitizenRegion(f.ctx, types.NewMsgSetCitizenRegion(f.registrar.String(), target.String(), "ua-2")); err != nil {
		t.Fatalf("update region: %v", err)
	}
	region, found = f.keeper.GetRegion(f.ctx, target)
	if !found || region != "ua-2" {
		t.Fatalf("region not updated")
	}
	addresses = f.keeper.GetAddressesByRegion(f.ctx, "ua-2")
	if len(addresses) != 1 || addresses[0] != target.String() {
		t.Fatalf("new index missing")
	}
	for _, addr := range f.keeper.GetAddressesByRegion(f.ctx, "ua-1") {
		if addr == target.String() {
			t.Fatalf("old index not cleared")
		}
	}
}

func TestSetCitizenRegionUnauthorized(t *testing.T) {
	f := initFixture(t)
	target := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	if _, err := f.msgServer.SetCitizenRegion(f.ctx, types.NewMsgSetCitizenRegion(f.nonMember.String(), target.String(), "ua-3")); err == nil {
		t.Fatalf("expected unauthorized error")
	}
}
