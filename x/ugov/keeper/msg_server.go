package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"uagd/x/ugov/types"
)

type MsgServer struct{ Keeper }

func NewMsgServerImpl(k Keeper) *MsgServer { return &MsgServer{Keeper: k} }

var _ types.MsgServer = MsgServer{}

func (s MsgServer) CreatePlan(goCtx context.Context, msg *types.MsgCreatePlan) (*types.MsgCreatePlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	id := s.getNextPlanID(ctx)
	plan := types.Plan{
		Id:          id,
		Creator:     msg.Creator,
		FundAddress: msg.FundAddress,
		Title:       msg.Title,
		Description: msg.Description,
		Status:      types.PlanStatus_PLAN_STATUS_DRAFT,
		Plan:        msg.Plan,
		ProposalId:  0,
	}

	if err := s.SetPlan(ctx, plan); err != nil {
		return nil, err
	}
	if err := s.bumpNextPlanID(ctx, id+1); err != nil {
		return nil, err
	}

	return &types.MsgCreatePlanResponse{Id: id}, nil
}

func (s MsgServer) UpdatePlan(goCtx context.Context, msg *types.MsgUpdatePlan) (*types.MsgUpdatePlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	plan, found := s.GetPlan(ctx, msg.Id)
	if !found {
		return nil, fmt.Errorf("plan not found")
	}
	if plan.Creator != msg.Creator {
		return nil, fmt.Errorf("only creator can update plan")
	}
	if plan.Status != types.PlanStatus_PLAN_STATUS_DRAFT {
		return nil, fmt.Errorf("plan must be DRAFT to update")
	}

	plan.Title = msg.Title
	plan.Description = msg.Description
	plan.Plan = msg.Plan

	if err := s.SetPlan(ctx, plan); err != nil {
		return nil, err
	}

	return &types.MsgUpdatePlanResponse{}, nil
}

func (s MsgServer) SubmitPlan(goCtx context.Context, msg *types.MsgSubmitPlan) (*types.MsgSubmitPlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	plan, found := s.GetPlan(ctx, msg.Id)
	if !found {
		return nil, fmt.Errorf("plan not found")
	}
	if plan.Creator != msg.Creator {
		return nil, fmt.Errorf("only creator can submit plan")
	}
	if plan.Status != types.PlanStatus_PLAN_STATUS_DRAFT {
		return nil, fmt.Errorf("plan must be DRAFT to submit")
	}

	plan.Status = types.PlanStatus_PLAN_STATUS_SUBMITTED
	plan.ProposalId = msg.ProposalId

	if err := s.SetPlan(ctx, plan); err != nil {
		return nil, err
	}

	return &types.MsgSubmitPlanResponse{}, nil
}

func (s MsgServer) ExecuteFundPlan(goCtx context.Context, msg *types.MsgExecuteFundPlan) (*types.MsgExecuteFundPlanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	params, err := s.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	if params.Authority != "" && params.Authority != msg.Authority {
		return nil, fmt.Errorf("unauthorized: expected authority %s", params.Authority)
	}

	plan, found := s.GetPlan(ctx, msg.PlanId)
	if !found {
		return nil, fmt.Errorf("plan not found")
	}
	if plan.Status != types.PlanStatus_PLAN_STATUS_SUBMITTED {
		return nil, fmt.Errorf("plan must be SUBMITTED to execute")
	}

	authorityAddr := sdk.MustAccAddressFromBech32(msg.Authority)

	// Ensure fund address is carried (just in case)
	fp := plan.Plan
	fp.FundAddress = plan.FundAddress

	if err := s.fundKeeper.ExecuteFundPlan(ctx, fp, authorityAddr); err != nil {
		return nil, err
	}

	plan.Status = types.PlanStatus_PLAN_STATUS_EXECUTED
	if err := s.SetPlan(ctx, plan); err != nil {
		return nil, err
	}

	return &types.MsgExecuteFundPlanResponse{}, nil
}
