package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"uagd/x/ugov/types"
)

type MsgServer struct{ Keeper }

func NewMsgServerImpl(k Keeper) *MsgServer { return &MsgServer{Keeper: k} }

func (s MsgServer) SetPresident(goCtx context.Context, msg *types.MsgSetPresident) (*types.President, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	if err := s.Keeper.MustBeAdmin(ctx, msg.Authority); err != nil {
		return nil, err
	}

	p := types.President{
		RoleType:    msg.RoleType,
		RegionId:    msg.RegionId,
		Address:     msg.Address,
		Active:      true,
		SinceHeight: ctx.BlockHeight(),
	}
	s.Keeper.SetPresident(ctx, p)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("ugov.president_set",
			sdk.NewAttribute("role_type", fmt.Sprintf("%d", msg.RoleType)),
			sdk.NewAttribute("region_id", msg.RegionId),
			sdk.NewAttribute("address", msg.Address),
		),
	)
	return &p, nil
}

func (s MsgServer) CreateFundPlan(goCtx context.Context, msg *types.MsgCreateFundPlan) (*types.StoredFundPlan, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// TODO: replace with fund->region mapping from x/fund.
	role := types.PRESIDENT_TYPE_NATIONAL
	regionId := ""

	id, err := s.Keeper.CreatePlan(ctx, msg.Creator, msg.FundAddress, msg.Title, msg.Description, role, regionId, msg.PlanJSON)
	if err != nil {
		return nil, err
	}

	p, _ := s.Keeper.GetPlan(ctx, id)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent("ugov.plan_created",
			sdk.NewAttribute("plan_id", fmt.Sprintf("%d", id)),
			sdk.NewAttribute("fund_address", msg.FundAddress),
			sdk.NewAttribute("creator", msg.Creator),
		),
	)
	return &p, nil
}

// Placeholder: in gov v1, you normally submit proposals via gov CLI with messages.
func (s MsgServer) SubmitFundPlanAsProposal(goCtx context.Context, msg *types.MsgSubmitFundPlanAsProposal) (*types.StoredFundPlan, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	if err := s.Keeper.MustBePresident(ctx, msg.Creator, types.PRESIDENT_TYPE_NATIONAL, ""); err != nil {
		return nil, err
	}

	if err := s.Keeper.MarkSubmitted(ctx, msg.PlanId, 0); err != nil {
		return nil, err
	}
	p, _ := s.Keeper.GetPlan(ctx, msg.PlanId)
	return &p, nil
}

func (s MsgServer) ExecuteFundPlan(goCtx context.Context, msg *types.MsgExecuteFundPlan) (*types.StoredFundPlan, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	if err := s.Keeper.MustBeGovAuthority(msg.Authority); err != nil {
		return nil, err
	}

	plan, ok := s.Keeper.GetPlan(ctx, msg.PlanId)
	if !ok {
		return nil, fmt.Errorf("plan not found")
	}
	if plan.Status != types.PLAN_STATUS_SUBMITTED && plan.Status != types.PLAN_STATUS_DRAFT {
		return nil, fmt.Errorf("plan must be SUBMITTED (or DRAFT for dev) to execute")
	}

	// TODO: decode plan.PlanJSON -> fundtypes.FundPlan and call:
	// err := s.fundKeeper.ExecuteFundPlan(ctx, decodedPlan, sdk.MustAccAddressFromBech32(msg.Authority))

	plan.Status = types.PLAN_STATUS_EXECUTED
	plan.ExecutedAtHeight = ctx.BlockHeight()
	s.Keeper.SetPlan(ctx, plan)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("ugov.plan_executed",
			sdk.NewAttribute("plan_id", fmt.Sprintf("%d", plan.Id)),
			sdk.NewAttribute("fund_address", plan.FundAddress),
		),
	)

	return &plan, nil
}
