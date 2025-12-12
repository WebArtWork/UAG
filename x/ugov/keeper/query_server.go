package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"uagd/x/ugov/types"
)

type QueryServer struct{ Keeper }

func NewQueryServer(k Keeper) QueryServer { return QueryServer{Keeper: k} }

type QueryPresidentRequest struct {
	RoleType types.PresidentRoleType `json:"role_type"`
	RegionId string                  `json:"region_id"`
}
type QueryPresidentResponse struct {
	President *types.President `json:"president"`
}

func (q QueryServer) President(goCtx context.Context, req *QueryPresidentRequest) (*QueryPresidentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	p, ok := q.Keeper.GetPresident(ctx, req.RoleType, req.RegionId)
	if !ok {
		return &QueryPresidentResponse{President: nil}, nil
	}
	return &QueryPresidentResponse{President: &p}, nil
}

type QueryPresidentsRequest struct{}
type QueryPresidentsResponse struct {
	Presidents []types.President `json:"presidents"`
}

func (q QueryServer) Presidents(goCtx context.Context, _ *QueryPresidentsRequest) (*QueryPresidentsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return &QueryPresidentsResponse{Presidents: q.Keeper.GetAllPresidents(ctx)}, nil
}

type QueryFundPlanRequest struct{ PlanId uint64 `json:"plan_id"` }
type QueryFundPlanResponse struct{ Plan *types.StoredFundPlan `json:"plan"` }

func (q QueryServer) FundPlan(goCtx context.Context, req *QueryFundPlanRequest) (*QueryFundPlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	p, ok := q.Keeper.GetPlan(ctx, req.PlanId)
	if !ok {
		return nil, fmt.Errorf("plan not found")
	}
	return &QueryFundPlanResponse{Plan: &p}, nil
}

type QueryFundPlansByStatusRequest struct{ Status types.FundPlanStatus `json:"status"` }
type QueryFundPlansByStatusResponse struct{ Plans []types.StoredFundPlan `json:"plans"` }

func (q QueryServer) FundPlansByStatus(goCtx context.Context, req *QueryFundPlansByStatusRequest) (*QueryFundPlansByStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return &QueryFundPlansByStatusResponse{Plans: q.Keeper.GetPlansByStatus(ctx, req.Status)}, nil
}
