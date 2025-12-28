package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"uagd/x/fund/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	types.UnimplementedQueryServer

	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	growthKeeper  types.GrowthKeeper
	govAuthority  sdk.AccAddress

	Schema    collections.Schema
	FundStore collections.Map[string, types.Fund]
	Params    collections.Item[types.Params]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	growthKeeper types.GrowthKeeper,
	govAuthority sdk.AccAddress,
) Keeper {
	// sdk.AccAddress is []byte, so nil can happen. Default to gov module authority.
	if govAuthority == nil {
		govAuthority = authtypes.NewModuleAddress(govtypes.ModuleName)
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		storeService:  storeService,
		cdc:           cdc,
		addressCodec:  addressCodec,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		growthKeeper:  growthKeeper,
		govAuthority:  govAuthority,
		FundStore:     collections.NewMap(sb, types.FundKeyPrefix, "funds", collections.StringKey, codec.CollValue[types.Fund](cdc)),
		Params:        collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

func (k Keeper) SetFund(ctx context.Context, fund types.Fund) error {
	return k.FundStore.Set(ctx, fund.Address, fund)
}

func (k Keeper) GetFund(ctx context.Context, addr sdk.AccAddress) (types.Fund, bool) {
	fund, err := k.FundStore.Get(ctx, addr.String())
	if err != nil {
		return types.Fund{}, false
	}
	return fund, true
}

func (k Keeper) MustGetFund(ctx context.Context, addr sdk.AccAddress) types.Fund {
	fund, found := k.GetFund(ctx, addr)
	if !found {
		panic(types.ErrFundNotFound)
	}
	return fund
}

func (k Keeper) GetAllFunds(ctx context.Context) []types.Fund {
	var funds []types.Fund
	iterator, err := k.FundStore.Iterate(ctx, nil)
	if err != nil {
		return funds
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		fund, err := iterator.Value()
		if err == nil {
			funds = append(funds, fund)
		}
	}
	return funds
}

func (k Keeper) GetFundsByType(ctx context.Context, fundType types.FundType) []types.Fund {
	funds := k.GetAllFunds(ctx)
	var filtered []types.Fund
	for _, f := range funds {
		if f.Type == fundType {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return types.Params{}, err
	}
	return params, nil
}

func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	return k.Params.Set(ctx, params)
}

func (k Keeper) assertGovAuthority(authority sdk.AccAddress) error {
	if k.govAuthority == nil || k.govAuthority.Empty() {
		return fmt.Errorf("gov authority not configured")
	}
	if !authority.Equals(k.govAuthority) {
		return types.ErrUnauthorized
	}
	return nil
}

func (k Keeper) ValidateFundPosition(ctx context.Context, position types.FundPosition) error {
	fundAddr, err := k.addressCodec.StringToBytes(position.FundAddress)
	if err != nil {
		return err
	}
	fund, found := k.GetFund(ctx, fundAddr)
	if !found {
		return types.ErrFundNotFound
	}
	if !fund.Active {
		return types.ErrFundInactive
	}

	// Regional occupation lock: if a region is more than 50% occupied, no outflows are allowed.
	if fund.Type == types.FundType_FUND_TYPE_REGION {
		occupation, found := k.growthKeeper.GetRegionOccupation(ctx, fund.RegionId)
		if found && occupation.GT(sdkmath.LegacyNewDec(50)) {
			return types.ErrRegionLocked
		}
	}

	delegationLimit, payrollLimit := k.growthKeeper.GetEffectiveLimits(ctx, fund)

	totalDelegations := sdkmath.ZeroInt()
	totalPayouts := sdkmath.ZeroInt()

	for _, d := range position.Delegations {
		// FundDelegation is a VALUE type (not *FundDelegation).
		if d.ValidatorAddress == "" {
			return fmt.Errorf("delegation validator_address is required")
		}
		if !d.Amount.IsValid() || d.Amount.IsZero() {
			return fmt.Errorf("delegation amount invalid")
		}
		if d.Amount.Denom != types.BaseDenom {
			return types.ErrInvalidDenom
		}
		totalDelegations = totalDelegations.Add(d.Amount.Amount)
	}

	for _, p := range position.Payouts {
		// FundPayout is a VALUE type (not *FundPayout).
		if p.RecipientAddress == "" {
			return fmt.Errorf("payout recipient_address is required")
		}
		if !p.Amount.IsValid() || p.Amount.IsZero() {
			return fmt.Errorf("payout amount invalid")
		}
		if p.Amount.Denom != types.BaseDenom {
			return types.ErrInvalidDenom
		}
		totalPayouts = totalPayouts.Add(p.Amount.Amount)
	}

	if delegationLimit.Amount.IsPositive() && totalDelegations.GT(delegationLimit.Amount) {
		return types.ErrDelegationLimit
	}
	if payrollLimit.Amount.IsPositive() && totalPayouts.GT(payrollLimit.Amount) {
		return types.ErrPayrollLimit
	}
	return nil
}

func (k Keeper) ExecuteFundPosition(ctx context.Context, position types.FundPosition, authority sdk.AccAddress) error {
	if err := k.assertGovAuthority(authority); err != nil {
		return err
	}
	if err := k.ValidateFundPosition(ctx, position); err != nil {
		return err
	}

	fundAddr, err := k.addressCodec.StringToBytes(position.FundAddress)
	if err != nil {
		return err
	}
	fundAcc := sdk.AccAddress(fundAddr)

	for _, del := range position.Delegations {
		if del.ValidatorAddress == "" {
			return fmt.Errorf("delegation validator_address is required")
		}
		if !del.Amount.IsValid() || del.Amount.IsZero() {
			return fmt.Errorf("delegation amount invalid")
		}

		valAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(del.ValidatorAddress)
		if err != nil {
			return err
		}

		validator, err := k.stakingKeeper.GetValidator(ctx, sdk.ValAddress(valAddr))
		if err != nil {
			return err
		}

		if _, err := k.stakingKeeper.Delegate(ctx, fundAcc, del.Amount.Amount, stakingtypes.Unbonded, validator, true); err != nil {
			return err
		}
	}

	for _, payout := range position.Payouts {
		if payout.RecipientAddress == "" {
			return fmt.Errorf("payout recipient_address is required")
		}
		if !payout.Amount.IsValid() || payout.Amount.IsZero() {
			return fmt.Errorf("payout amount invalid")
		}

		toAddr, err := k.addressCodec.StringToBytes(payout.RecipientAddress)
		if err != nil {
			return err
		}

		// payout.Amount is a VALUE type (not *Coin).
		if err := k.bankKeeper.SendCoins(ctx, fundAcc, sdk.AccAddress(toAddr), sdk.NewCoins(payout.Amount)); err != nil {
			return err
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"fund_position_executed",
			sdk.NewAttribute("fund", position.FundAddress),
			sdk.NewAttribute("position_id", fmt.Sprintf("%d", position.Id)),
		),
	)

	return nil
}
