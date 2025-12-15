package growth

import (
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"uagd/x/growth/keeper"
	"uagd/x/growth/types"
)

type Inputs struct {
	StoreService corestore.KVStoreService
	Cdc          codec.Codec
	AddressCodec address.Codec
}

type Outputs struct {
	GrowthKeeper keeper.Keeper
}

func ProvideModule(in Inputs) Outputs {
	k := keeper.NewKeeper(in.StoreService, in.Cdc, in.AddressCodec)
	return Outputs{GrowthKeeper: k}
}

func ProvideStoreKey() string {
	return types.StoreKey
}
