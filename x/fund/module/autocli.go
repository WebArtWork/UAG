package fund

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"uagd/x/fund/types"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "Fund", Use: "fund [address]", Short: "Get a fund by address"},
				{RpcMethod: "Funds", Use: "funds", Short: "List all funds"},
				{RpcMethod: "FundsByType", Use: "funds-by-type", Short: "List funds by type"},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:           types.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{{RpcMethod: "ExecuteFundPlan", Use: "execute-plan", Short: "Execute a fund plan"}},
		},
	}
}
