package app

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fundtypes "uagd/x/fund/types"
	growthkeeper "uagd/x/growth/keeper"
)

// ProvideFundGrowthKeeper binds the concrete x/growth keeper to the x/fund expected keeper interface.
//
// Note: for now, limits are permissive (0 = no limit). We can tighten this once growth params + scoring
// are finalized.
func ProvideFundGrowthKeeper(gk growthkeeper.Keeper) fundtypes.GrowthKeeper {
	return fundGrowthAdapter{gk: gk}
}

type fundGrowthAdapter struct {
	gk growthkeeper.Keeper
}

func (a fundGrowthAdapter) GetEffectiveLimits(_ context.Context, _ fundtypes.Fund) (sdk.Coin, sdk.Coin) {
	zero := sdk.NewCoin(fundtypes.BaseDenom, sdkmath.ZeroInt())
	return zero, zero
}

func (a fundGrowthAdapter) GetRegionOccupation(ctx context.Context, regionID string) (sdkmath.LegacyDec, bool) {
	o, found, err := a.gk.GetOccupation(ctx, regionID)
	if err != nil || !found {
		return sdkmath.LegacyDec{}, false
	}
	return o.Occupation, true
}
