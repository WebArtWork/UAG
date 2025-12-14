package types

import (
	"context"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type BankKeeper interface {
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type StakingKeeper interface {
	ValidatorAddressCodec() address.Codec
	GetValidator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error)
	Delegate(ctx context.Context, delAddr sdk.AccAddress, bondAmt math.Int, tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (math.LegacyDec, error)
}

type GrowthKeeper interface {
	GetEffectiveLimits(ctx context.Context, fund Fund) (delegationLimit sdk.Coin, payrollLimit sdk.Coin)
	GetRegionOccupation(ctx context.Context, regionID string) (math.LegacyDec, bool)
}

type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
	AddressCodec() address.Codec
}
