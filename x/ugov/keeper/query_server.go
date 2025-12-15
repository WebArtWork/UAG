package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"uagd/x/ugov/types"
)

// NOTE:
// Your current ugov keeper code references proto types that don't exist in your generated ugov `types`
// (e.g. types.President, types.StoredFundPlan). Until the ugov protos are aligned, we keep the
// query server compiling by returning Unimplemented.
//
// Once you confirm the real message names in proto, we’ll wire storage + real responses.

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryParamsResponse{Params: params}, nil
}

// If your proto defines other queries (presidents/plans/etc), keep them compiling as stubs.
// Replace these method signatures to match exactly what `types.QueryServer` requires in your repo.

func (k Keeper) President(ctx context.Context, req *types.QueryPresidentRequest) (*types.QueryPresidentResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ugov: President query not implemented yet (proto/types mismatch)")
}

func (k Keeper) Presidents(ctx context.Context, req *types.QueryPresidentsRequest) (*types.QueryPresidentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ugov: Presidents query not implemented yet (proto/types mismatch)")
}

func (k Keeper) FundPlan(ctx context.Context, req *types.QueryFundPlanRequest) (*types.QueryFundPlanResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ugov: FundPlan query not implemented yet (proto/types mismatch)")
}

func (k Keeper) FundPlans(ctx context.Context, req *types.QueryFundPlansRequest) (*types.QueryFundPlansResponse, error) {
	return nil, status.Error(codes.Unimplemented, "ugov: FundPlans query not implemented yet (proto/types mismatch)")
}
