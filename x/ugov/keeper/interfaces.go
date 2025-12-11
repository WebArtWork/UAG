package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FundKeeper is the minimal surface ugov needs from x/fund.
// Replace plan interface{} with fundtypes.FundPlan once imported.
type FundKeeper interface {
	ExecuteFundPlan(ctx sdk.Context, plan interface{}, authority sdk.AccAddress) error
}
