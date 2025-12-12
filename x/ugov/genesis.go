package ugov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"uagd/x/ugov/keeper"
	"uagd/x/ugov/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, gen types.GenesisState) {
	k.SetParams(ctx, gen.Params)
	for _, p := range gen.Presidents {
		k.SetPresident(ctx, p)
	}
	for _, pl := range gen.Plans {
		k.SetPlan(ctx, pl)
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	gs := types.DefaultGenesis()
	gs.Params = k.GetParams(ctx)
	gs.Presidents = k.GetAllPresidents(ctx)
	gs.Plans = append(gs.Plans, k.GetPlansByStatus(ctx, types.PLAN_STATUS_DRAFT)...)
	gs.Plans = append(gs.Plans, k.GetPlansByStatus(ctx, types.PLAN_STATUS_SUBMITTED)...)
	gs.Plans = append(gs.Plans, k.GetPlansByStatus(ctx, types.PLAN_STATUS_EXECUTED)...)
	gs.Plans = append(gs.Plans, k.GetPlansByStatus(ctx, types.PLAN_STATUS_REJECTED)...)
	return &gs
}
