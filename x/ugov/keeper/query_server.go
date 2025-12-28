package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"uagd/x/ugov/types"
)

type QueryServer struct{ Keeper }

func NewQueryServer(k Keeper) *QueryServer { return &QueryServer{Keeper: k} }

var _ types.QueryServer = QueryServer{}

func (q QueryServer) Plan(ctx context.Context, req *types.QueryPlanRequest) (*types.QueryPlanResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id required")
	}

	planVal, found := q.GetPlan(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, "plan not found")
	}

	plan := planVal // create addressable copy
	return &types.QueryPlanResponse{Plan: &plan}, nil
}

func (q QueryServer) Plans(ctx context.Context, _ *types.QueryPlansRequest) (*types.QueryPlansResponse, error) {
	vals := q.GetAllPlans(ctx)

	out := make([]*types.Plan, 0, len(vals))
	for i := range vals {
		p := vals[i] // addressable copy per iteration index
		out = append(out, &p)
	}

	return &types.QueryPlansResponse{Plans: out}, nil
}
