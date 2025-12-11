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
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		storeService:  storeService,
		cdc:           cdc,
		addressCodec:  addressCodec,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
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

func (k Keeper) ValidateFundPlan(ctx context.Context, plan types.FundPlan) error {
	fundAddr, err := k.addressCodec.StringToBytes(plan.FundAddress)
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
	totalDelegations := sdkmath.ZeroInt()
	totalPayouts := sdkmath.ZeroInt()
	for _, d := range plan.Delegations {
		if d == nil || d.Amount == nil {
			return fmt.Errorf("delegation entry invalid")
		}
		if d.Amount.Denom != types.BaseDenom {
			return types.ErrInvalidDenom
		}
		totalDelegations = totalDelegations.Add(d.Amount.Amount)
	}
	for _, p := range plan.Payouts {
		if p == nil || p.Amount == nil {
			return fmt.Errorf("payout entry invalid")
		}
		if p.Amount.Denom != types.BaseDenom {
			return types.ErrInvalidDenom
		}
		totalPayouts = totalPayouts.Add(p.Amount.Amount)
	}
	if fund.BaseDelegationLimit != nil && fund.BaseDelegationLimit.Amount.IsPositive() && totalDelegations.GT(fund.BaseDelegationLimit.Amount) {
		return types.ErrDelegationLimit
	}
	if fund.BasePayrollLimit != nil && fund.BasePayrollLimit.Amount.IsPositive() && totalPayouts.GT(fund.BasePayrollLimit.Amount) {
		return types.ErrPayrollLimit
	}
	return nil
}

func (k Keeper) ExecuteFundPlan(ctx context.Context, plan types.FundPlan, authority sdk.AccAddress) error {
	if err := k.ValidateFundPlan(ctx, plan); err != nil {
		return err
	}
	fundAddr, err := k.addressCodec.StringToBytes(plan.FundAddress)
	if err != nil {
		return err
	}
	fundAcc := sdk.AccAddress(fundAddr)

	for _, del := range plan.Delegations {
		if del == nil || del.Amount == nil {
			return fmt.Errorf("delegation entry invalid")
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

	for _, payout := range plan.Payouts {
		if payout == nil || payout.Amount == nil {
			return fmt.Errorf("payout entry invalid")
		}
		toAddr, err := k.addressCodec.StringToBytes(payout.RecipientAddress)
		if err != nil {
			return err
		}
		if err := k.bankKeeper.SendCoins(ctx, fundAcc, sdk.AccAddress(toAddr), sdk.NewCoins(*payout.Amount)); err != nil {
			return err
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"fund_plan_executed",
			sdk.NewAttribute("fund", plan.FundAddress),
			sdk.NewAttribute("plan_id", fmt.Sprintf("%d", plan.Id)),
		),
	)
	return nil
}
