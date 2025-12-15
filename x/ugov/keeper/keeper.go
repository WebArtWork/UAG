package keeper

import (
	"context"

	"cosmossdk.io/collections"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	fundkeeper "uagd/x/fund/keeper"
	"uagd/x/ugov/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	types.UnimplementedQueryServer

	Schema collections.Schema

	Params     collections.Item[types.Params]
	Plans      collections.Map[uint64, types.Plan]
	NextPlanID collections.Item[uint64]

	fundKeeper fundkeeper.Keeper
	authority  string // bech32
}

func NewKeeper(
	cdc codec.Codec,
	storeService corestore.KVStoreService,
	_ any, // old templates: params subspace; not needed here
	fundKeeper fundkeeper.Keeper,
	authority sdk.AccAddress,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,

		Params:     collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Plans:      collections.NewMap(sb, types.PlanKeyPrefix, "plans", collections.Uint64Key, codec.CollValue[types.Plan](cdc)),
		NextPlanID: collections.NewItem(sb, types.NextPlanIDKey, "next_plan_id", collections.Uint64Value),

		fundKeeper: fundKeeper,
		authority:  authority.String(),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	p, err := k.Params.Get(ctx)
	if err != nil {
		dp := types.DefaultParams()
		dp.Authority = k.authority
		return dp, nil
	}
	if p.Authority == "" {
		p.Authority = k.authority
	}
	return p, nil
}

func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	if params.Authority == "" {
		params.Authority = k.authority
	}
	if err := params.Validate(); err != nil {
		return err
	}
	return k.Params.Set(ctx, params)
}

func (k Keeper) getNextPlanID(ctx context.Context) uint64 {
	id, err := k.NextPlanID.Get(ctx)
	if err != nil || id == 0 {
		return 1
	}
	return id
}

func (k Keeper) bumpNextPlanID(ctx context.Context, next uint64) error {
	if next == 0 {
		next = 1
	}
	return k.NextPlanID.Set(ctx, next)
}

func (k Keeper) GetPlan(ctx context.Context, id uint64) (types.Plan, bool) {
	p, err := k.Plans.Get(ctx, id)
	if err != nil {
		return types.Plan{}, false
	}
	return p, true
}

func (k Keeper) SetPlan(ctx context.Context, p types.Plan) error {
	return k.Plans.Set(ctx, p.Id, p)
}

func (k Keeper) GetAllPlans(ctx context.Context) []types.Plan {
	out := make([]types.Plan, 0)
	it, err := k.Plans.Iterate(ctx, nil)
	if err != nil {
		return out
	}
	defer it.Close()
	for ; it.Valid(); it.Next() {
		v, err := it.Value()
		if err == nil {
			out = append(out, v)
		}
	}
	return out
}

// FIX: gov_hooks.go expects sdk.Context and NO return value.
func (k Keeper) MarkRejectedByProposalId(_ sdk.Context, _ uint64) {
	// no-op for now
}
