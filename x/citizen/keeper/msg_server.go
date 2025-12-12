package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"uagd/x/citizen/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
	types.UnimplementedMsgServer
}

func NewMsgServerImpl(k Keeper) types.MsgServer {
	return msgServer{Keeper: k}
}

func (m msgServer) SetCitizenRegion(ctx context.Context, msg *types.MsgSetCitizenRegion) (*types.MsgSetCitizenRegionResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	registrar, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}
	if !m.IsRegistrar(ctx, registrar) {
		return nil, types.ErrUnauthorized
	}
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}
	if err := m.Keeper.SetRegion(ctx, addr, msg.RegionId); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCitizenRegionSet,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address),
			sdk.NewAttribute(types.AttributeKeyRegionID, msg.RegionId),
		),
	)
	return &types.MsgSetCitizenRegionResponse{}, nil
}

func (m msgServer) ClearCitizenRegion(ctx context.Context, msg *types.MsgClearCitizenRegion) (*types.MsgClearCitizenRegionResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	registrar, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return nil, err
	}
	if !m.IsRegistrar(ctx, registrar) {
		return nil, types.ErrUnauthorized
	}
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}
	if err := m.Keeper.DeleteRegion(ctx, addr); err != nil {
		return nil, err
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCitizenRegionCleared,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address),
		),
	)
	return &types.MsgClearCitizenRegionResponse{}, nil
}
