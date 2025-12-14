package fund

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
"github.com/cosmos/cosmos-sdk/codec"
authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"uagd/x/fund/keeper"
	"uagd/x/fund/types"
)

var _ depinject.OnePerModuleType = AppModule{}

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
	AddressCodec address.Codec

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
	StakingKeeper types.StakingKeeper
}

type ModuleOutputs struct {
	depinject.Out

	FundKeeper keeper.Keeper
	Module     appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
authority := authtypes.NewModuleAddress(govtypes.ModuleName)
if in.Config != nil && in.Config.Authority != "" {
authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
}

k := keeper.NewKeeper(in.StoreService, in.Cdc, in.AddressCodec, in.BankKeeper, in.StakingKeeper, authority)
m := NewAppModule(in.Cdc, k, in.AccountKeeper, in.BankKeeper)
return ModuleOutputs{FundKeeper: k, Module: m}
}
