package ugov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"uagd/x/ugov/keeper"
	"uagd/x/ugov/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, gen types.GenesisState) {
	_ = k.SetParams(ctx, gen.Params)

	// Store plans
	for _, p := range gen.Plans {
		_ = k.SetPlan(ctx, p)
	}

	// NextPlanId
	next := gen.NextPlanId
	if next == 0 {
		next = 1
	}
	_ = k.NextPlanID.Set(ctx, next)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params, _ := k.GetParams(ctx)

	nextID, err := k.NextPlanID.Get(ctx)
	if err != nil || nextID == 0 {
		nextID = 1
	}

	return &types.GenesisState{
		Params:     params,
		Plans:      k.GetAllPlans(ctx),
		NextPlanId: nextID,
	}
}
