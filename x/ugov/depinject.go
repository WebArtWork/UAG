package ugov

import (
	"cosmossdk.io/core/appconfig"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"

	fundkeeper "uagd/x/fund/keeper"
	"uagd/x/ugov/keeper"
	"uagd/x/ugov/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType ensures only one instance of this module is wired by depinject.
func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.Register(
		&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec

	ParamsKeeper paramskeeper.Keeper
	FundKeeper   fundkeeper.Keeper
}

type ModuleOutputs struct {
	depinject.Out

	UgovKeeper keeper.Keeper
	Module     appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config != nil && in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	paramsSubspace := in.ParamsKeeper.Subspace(types.ModuleName)

	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		paramsSubspace,
		in.FundKeeper,
		authority,
	)
	m := NewAppModule(in.Cdc, k)

	return ModuleOutputs{UgovKeeper: k, Module: m}
}
