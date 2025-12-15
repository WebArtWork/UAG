package growth

import (
	"cosmossdk.io/core/appmodule"

	"uagd/x/growth/keeper"
	"uagd/x/growth/types"
)

var _ appmodule.AppModule = AppModule{}

type AppModule struct {
	keeper keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{keeper: k}
}

// Required by cosmossdk.io/core/appmodule in your SDK build.
func (AppModule) IsAppModule() {}

// Required by cosmossdk.io/core/appmodule in your SDK build.
func (AppModule) IsOnePerModuleType() {}

func (AppModule) Name() string { return types.ModuleName }
