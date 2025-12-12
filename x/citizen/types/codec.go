package types

import (
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
)

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetCitizenRegion{},
		&MsgClearCitizenRegion{},
	)
	registry.RegisterImplementations((*tx.MsgResponse)(nil),
		&MsgSetCitizenRegionResponse{},
		&MsgClearCitizenRegionResponse{},
	)
}

var (
	ModuleCdc     codec.Codec
	moduleCdcOnce sync.Once
)

func initModuleCdc() {
	registry := codectypes.NewInterfaceRegistry()
	RegisterInterfaces(registry)
	ModuleCdc = codec.NewProtoCodec(registry)
}

func getModuleCodec() codec.Codec {
	moduleCdcOnce.Do(initModuleCdc)
	return ModuleCdc
}

func RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}
