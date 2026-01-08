package citizen

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"uag/x/citizen/types"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "RegionByAddress", Use: "region [address]", Short: "Get region for an address"},
				{RpcMethod: "AddressesByRegion", Use: "addresses-by-region [region-id]", Short: "List addresses for a region"},
				{RpcMethod: "Params", Use: "params", Short: "Show citizen module params"},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: types.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "SetCitizenRegion", Use: "set-region [address] [region-id]", Short: "Set the region for a citizen"},
				{RpcMethod: "ClearCitizenRegion", Use: "clear-region [address]", Short: "Clear the region mapping for a citizen"},
			},
		},
	}
}
