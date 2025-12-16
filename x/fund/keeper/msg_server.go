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

func (m msgServer) ExecuteFundPosition(ctx context.Context, msg *types.MsgExecuteFundPosition) (*types.MsgExecuteFundPositionResponse, error) {
	return nil, types.ErrDirectExecDisabled
}
