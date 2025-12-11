package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"uagd/x/growth/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
	types.UnimplementedMsgServer
}

func NewMsgServerImpl(k Keeper) types.MsgServer {
	return msgServer{Keeper: k}
}

func (m msgServer) SetRegionMetric(ctx context.Context, msg *types.MsgSetRegionMetric) (*types.MsgSetRegionMetricResponse, error) {
	params, err := m.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	if params.Oracle != msg.Authority {
		return nil, types.ErrUnauthorized
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	if err := m.Keeper.SetRegionMetric(ctx, types.RegionMetric{
		RegionId:     msg.RegionId,
		Period:       msg.Period,
		TaxIndex:     msg.TaxIndex,
		GdpIndex:     msg.GdpIndex,
		ExportsIndex: msg.ExportsIndex,
	}); err != nil {
		return nil, err
	}
	score := m.ComputeGrowthScore(types.RegionMetric{
		RegionId:     msg.RegionId,
		Period:       msg.Period,
		TaxIndex:     msg.TaxIndex,
		GdpIndex:     msg.GdpIndex,
		ExportsIndex: msg.ExportsIndex,
	})

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"region_metric_updated",
			sdk.NewAttribute("region_id", msg.RegionId),
			sdk.NewAttribute("period", msg.Period),
			sdk.NewAttribute("delegation_multiplier", score.DelegationMultiplier),
			sdk.NewAttribute("payroll_multiplier", score.PayrollMultiplier),
		),
	)

	return &types.MsgSetRegionMetricResponse{Score: &score}, nil
}
