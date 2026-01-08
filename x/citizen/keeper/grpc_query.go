package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"uag/x/citizen/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) RegionByAddress(ctx context.Context, req *types.QueryRegionByAddressRequest) (*types.QueryRegionByAddressResponse, error) {
	if req == nil || req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address required")
	}
	if _, err := k.addressCodec.StringToBytes(req.Address); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	region, found := k.GetRegionByString(ctx, req.Address)
	if !found {
		return nil, status.Error(codes.NotFound, "region not found")
	}

	return &types.QueryRegionByAddressResponse{RegionId: region}, nil
}

func (k Keeper) AddressesByRegion(ctx context.Context, req *types.QueryAddressesByRegionRequest) (*types.QueryAddressesByRegionResponse, error) {
	if req == nil || req.RegionId == "" {
		return nil, status.Error(codes.InvalidArgument, "region id required")
	}
	addresses := k.GetAddressesByRegion(ctx, req.RegionId)
	return &types.QueryAddressesByRegionResponse{Addresses: addresses}, nil
}

func (k Keeper) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	// QueryParamsResponse.Params is a VALUE (not *Params).
	return &types.QueryParamsResponse{Params: params}, nil
}
