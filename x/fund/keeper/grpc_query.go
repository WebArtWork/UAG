package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"uagd/x/fund/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Fund(ctx context.Context, req *types.QueryFundRequest) (*types.QueryFundResponse, error) {
	if req == nil || req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address required")
	}
	addr, err := k.addressCodec.StringToBytes(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	fund, found := k.GetFund(ctx, addr)
	if !found {
		return nil, status.Error(codes.NotFound, "fund not found")
	}
	return &types.QueryFundResponse{Fund: &fund}, nil
}

func (k Keeper) Funds(ctx context.Context, _ *types.QueryFundsRequest) (*types.QueryFundsResponse, error) {
	funds := k.GetAllFunds(ctx)
	fundPtrs := make([]*types.Fund, len(funds))
	for i := range funds {
		fundPtrs[i] = &funds[i]
	}
	return &types.QueryFundsResponse{Funds: fundPtrs}, nil
}

func (k Keeper) FundsByType(ctx context.Context, req *types.QueryFundsByTypeRequest) (*types.QueryFundsByTypeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request required")
	}
	funds := k.GetFundsByType(ctx, req.Type)
	var filtered []*types.Fund
	if req.RegionId != "" {
		for _, f := range funds {
			if f.RegionId == req.RegionId {
				filtered = append(filtered, &f)
			}
		}
	} else {
		filtered = make([]*types.Fund, len(funds))
		for i := range funds {
			filtered[i] = &funds[i]
		}
	}
	return &types.QueryFundsByTypeResponse{Funds: filtered}, nil
}
