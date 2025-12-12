package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	fundtypes "uagd/x/fund/types"
)

// FundKeeper is the minimal surface ugov needs from x/fund.
// Replace plan interface{} with fundtypes.FundPlan once imported.
type FundKeeper interface {
	ExecuteFundPlan(ctx context.Context, plan fundtypes.FundPlan, authority sdk.AccAddress) error
}
