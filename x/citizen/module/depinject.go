package citizen

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"

	"uag/x/citizen/keeper"
	"uag/x/citizen/types"
)

var _ depinject.OnePerModuleType = AppModule{}

func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.RegisterModule(
		&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec
	AddressCodec address.Codec
}

type ModuleOutputs struct {
	depinject.Out

	CitizenKeeper keeper.Keeper
	Module        appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(in.StoreService, in.Cdc, in.AddressCodec)
	m := NewAppModule(in.Cdc, k)
	return ModuleOutputs{CitizenKeeper: k, Module: m}
}
