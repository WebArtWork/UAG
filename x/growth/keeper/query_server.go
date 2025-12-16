package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"uagd/x/growth/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) RegionMetric(ctx context.Context, req *types.QueryRegionMetricRequest) (*types.QueryRegionMetricResponse, error) {
	if req == nil || req.RegionId == "" {
		return nil, status.Error(codes.InvalidArgument, "region_id is required")
	}

	metric, found, err := k.GetRegionMetric(ctx, req.RegionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get region metric: %v", err)
	}
	if !found {
		return &types.QueryRegionMetricResponse{Metric: nil}, nil
	}

	return &types.QueryRegionMetricResponse{Metric: &metric}, nil
}

func (k Keeper) GrowthScore(ctx context.Context, req *types.QueryGrowthScoreRequest) (*types.QueryGrowthScoreResponse, error) {
	if req == nil || req.RegionId == "" {
		return nil, status.Error(codes.InvalidArgument, "region_id is required")
	}

	score, found, err := k.GetGrowthScore(ctx, req.RegionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get growth score: %v", err)
	}
	if !found {
		return &types.QueryGrowthScoreResponse{Score: nil}, nil
	}

	return &types.QueryGrowthScoreResponse{Score: &score}, nil
}
