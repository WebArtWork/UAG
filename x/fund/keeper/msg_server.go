package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

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
	params, err := m.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	if params.Admin != msg.Authority {
		return nil, types.ErrUnauthorized
	}
	authAddr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}
	if err := m.Keeper.ExecuteFundPlan(ctx, *msg.Plan, authAddr); err != nil {
		return nil, err
	}
	return &types.MsgExecuteFundPlanResponse{}, nil
}
