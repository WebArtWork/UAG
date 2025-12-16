package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	fundtypes "uagd/x/fund/types"
)

// FundKeeper is the minimal surface ugov needs from x/fund.
type FundKeeper interface {
	ExecuteFundPosition(ctx context.Context, position fundtypes.FundPosition, authority sdk.AccAddress) error
}
