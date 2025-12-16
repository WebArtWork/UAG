package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/ugov interfaces and concrete types on the provided LegacyAmino codec.
// NOTE: We keep this for CLI / legacy compatibility, even if most flows use protobuf.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// no sdk.Msg concrete registrations needed with ADR-031, kept intentionally minimal
}

// RegisterInterfaces registers the x/ugov interfaces and implementations with the interface registry.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// This is REQUIRED so baseapp can resolve type_url -> concrete Msg for Msg service routes.
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
