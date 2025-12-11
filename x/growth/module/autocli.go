package growth

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"uagd/x/growth/types"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "RegionMetric", Use: "region-metric [region-id] [period]", Short: "Get a region metric"},
				{RpcMethod: "GrowthScore", Use: "growth-score [region-id] [period]", Short: "Get a growth score"},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: types.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "SetRegionMetric", Use: "set-region-metric [region-id] [period] [tax] [gdp] [exports]", Short: "Set growth metrics for a region"},
			},
		},
	}
}
