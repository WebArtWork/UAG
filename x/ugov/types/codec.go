package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterInterfaces(_ codectypes.InterfaceRegistry) {}

func NewCodec() *codec.ProtoCodec {
	ir := codectypes.NewInterfaceRegistry()
	RegisterInterfaces(ir)
	return codec.NewProtoCodec(ir)
}
