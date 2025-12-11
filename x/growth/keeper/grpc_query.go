package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"uagd/x/growth/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) RegionMetric(ctx context.Context, req *types.QueryRegionMetricRequest) (*types.QueryRegionMetricResponse, error) {
	if req == nil || req.RegionId == "" || req.Period == "" {
		return nil, status.Error(codes.InvalidArgument, "region id and period required")
	}
	metric, found := k.GetRegionMetric(ctx, req.RegionId, req.Period)
	if !found {
		return nil, status.Error(codes.NotFound, "metric not found")
	}
	return &types.QueryRegionMetricResponse{Metric: &metric}, nil
}

func (k Keeper) GrowthScore(ctx context.Context, req *types.QueryGrowthScoreRequest) (*types.QueryGrowthScoreResponse, error) {
	if req == nil || req.RegionId == "" || req.Period == "" {
		return nil, status.Error(codes.InvalidArgument, "region id and period required")
	}
	score, found := k.GetGrowthScore(ctx, req.RegionId, req.Period)
	if !found {
		return nil, status.Error(codes.NotFound, "score not found")
	}
	return &types.QueryGrowthScoreResponse{Score: &score}, nil
}
