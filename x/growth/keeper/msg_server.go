package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"uagd/x/growth/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

func (m msgServer) SetRegionMetric(goCtx context.Context, msg *types.MsgSetRegionMetric) (*types.MsgSetRegionMetricResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg == nil {
		return nil, fmt.Errorf("empty message")
	}
	if msg.RegionId == "" {
		return nil, fmt.Errorf("region_id required")
	}
	if msg.Period == "" {
		return nil, fmt.Errorf("period required")
	}

	metric := types.RegionMetric{
		RegionId:     msg.RegionId,
		Period:       msg.Period,
		TaxIndex:     msg.TaxIndex,
		GdpIndex:     msg.GdpIndex,
		ExportsIndex: msg.ExportsIndex,
	}

	if err := m.Keeper.SetRegionMetric(ctx, metric); err != nil {
		return nil, err
	}

	// Best-effort: compute and store a score for the region.
	_ = m.computeAndStoreScore(ctx, msg.RegionId)

	return &types.MsgSetRegionMetricResponse{}, nil
}

// computeAndStoreScore is intentionally minimal for now.
// Replace later with your real growth formula.
func (m msgServer) computeAndStoreScore(ctx context.Context, regionID string) error {
	one := sdkmath.LegacyNewDec(1)

	score := types.GrowthScore{
		RegionId:             regionID,
		Score:                one,
		DelegationMultiplier: one,
		PayrollMultiplier:    one,
	}

	return m.Keeper.SetGrowthScore(ctx, score)
}
