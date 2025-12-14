package keeper

import (
	"context"

	"uagd/x/fund/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
	types.UnimplementedMsgServer
}

func NewMsgServerImpl(k Keeper) types.MsgServer {
	return msgServer{Keeper: k}
}

func (m msgServer) ExecuteFundPlan(ctx context.Context, msg *types.MsgExecuteFundPlan) (*types.MsgExecuteFundPlanResponse, error) {
	return nil, types.ErrDirectExecDisabled
}
