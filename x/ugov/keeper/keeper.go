package keeper

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	corestore "cosmossdk.io/core/store"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"uagd/x/ugov/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeService corestore.KVStoreService
	ps           paramtypes.Subspace

	fundKeeper   FundKeeper
	govAuthority sdk.AccAddress
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService corestore.KVStoreService,
	ps paramtypes.Subspace,
	fundKeeper FundKeeper,
	govAuthority sdk.AccAddress,
) Keeper {
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(paramtypes.NewKeyTable().RegisterParamSet(&paramsWrapper{}))
	}
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		ps:           ps,
		fundKeeper:   fundKeeper,
		govAuthority: govAuthority,
	}
}

type paramsWrapper struct{ Params types.Params }

func (p *paramsWrapper) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair([]byte("Admin"), &p.Params.Admin, validateAdmin),
	}
}

func validateAdmin(i interface{}) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid param type")
	}
	if s == "" {
		return nil
	}
	_, err := sdk.AccAddressFromBech32(s)
	return err
}

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var w paramsWrapper
	k.ps.GetParamSet(ctx, &w)
	return w.Params
}

func (k Keeper) SetParams(ctx sdk.Context, p types.Params) {
	k.ps.SetParamSet(ctx, &paramsWrapper{Params: p})
}

func (k Keeper) GovAuthority() sdk.AccAddress { return k.govAuthority }

func (k Keeper) store(ctx sdk.Context) storetypes.KVStore {
	return runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx.Context()))
}

func presidentKey(role types.PresidentRoleType, regionId string) []byte {
	b := []byte{byte(role), byte(len(regionId))}
	b = append(b, []byte(regionId)...)
	return append(types.PresidentKeyPrefix, b...)
}

func planKey(id uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	return append(types.PlanKeyPrefix, b...)
}

// --- Presidents ---

func (k Keeper) SetPresident(ctx sdk.Context, p types.President) {
	bz := mustMarshalJSON(p)
	k.store(ctx).Set(presidentKey(p.RoleType, p.RegionId), bz)
}

func (k Keeper) GetPresident(ctx sdk.Context, role types.PresidentRoleType, regionId string) (types.President, bool) {
	bz := k.store(ctx).Get(presidentKey(role, regionId))
	if bz == nil {
		return types.President{}, false
	}
	var p types.President
	mustUnmarshalJSON(bz, &p)
	return p, true
}

func (k Keeper) GetAllPresidents(ctx sdk.Context) []types.President {
	it := storetypes.KVStorePrefixIterator(k.store(ctx), types.PresidentKeyPrefix)
	defer it.Close()

	out := []types.President{}
	for ; it.Valid(); it.Next() {
		var p types.President
		mustUnmarshalJSON(it.Value(), &p)
		out = append(out, p)
	}
	return out
}

// --- Plans ---

func (k Keeper) nextPlanId(ctx sdk.Context) uint64 {
	bz := k.store(ctx).Get(types.PlanSeqKey)
	var cur uint64
	if bz != nil {
		cur = binary.BigEndian.Uint64(bz)
	}
	cur++
	nbz := make([]byte, 8)
	binary.BigEndian.PutUint64(nbz, cur)
	k.store(ctx).Set(types.PlanSeqKey, nbz)
	return cur
}

func (k Keeper) SetPlan(ctx sdk.Context, p types.StoredFundPlan) {
	bz := mustMarshalJSON(p)
	k.store(ctx).Set(planKey(p.Id), bz)
}

func (k Keeper) GetPlan(ctx sdk.Context, id uint64) (types.StoredFundPlan, bool) {
	bz := k.store(ctx).Get(planKey(id))
	if bz == nil {
		return types.StoredFundPlan{}, false
	}
	var p types.StoredFundPlan
	mustUnmarshalJSON(bz, &p)
	return p, true
}

func (k Keeper) GetPlansByStatus(ctx sdk.Context, status types.FundPlanStatus) []types.StoredFundPlan {
	it := storetypes.KVStorePrefixIterator(k.store(ctx), types.PlanKeyPrefix)
	defer it.Close()
	out := []types.StoredFundPlan{}
	for ; it.Valid(); it.Next() {
		var p types.StoredFundPlan
		mustUnmarshalJSON(it.Value(), &p)
		if p.Status == status {
			out = append(out, p)
		}
	}
	return out
}

// --- Auth helpers ---

func (k Keeper) MustBeAdmin(ctx sdk.Context, authority string) error {
	p := k.GetParams(ctx)
	if p.Admin == "" {
		return fmt.Errorf("ugov admin not set")
	}
	if authority != p.Admin {
		return fmt.Errorf("unauthorized: expected ugov admin")
	}
	return nil
}

func (k Keeper) MustBeGovAuthority(authority string) error {
	if k.govAuthority.Empty() {
		return fmt.Errorf("gov authority not configured")
	}
	if authority != k.govAuthority.String() {
		return fmt.Errorf("unauthorized: expected gov authority")
	}
	return nil
}

func (k Keeper) MustBePresident(ctx sdk.Context, creator string, role types.PresidentRoleType, regionId string) error {
	p, ok := k.GetPresident(ctx, role, regionId)
	if !ok || !p.Active {
		return fmt.Errorf("no active president for role/region")
	}
	if p.Address != creator {
		return fmt.Errorf("unauthorized: creator is not president")
	}
	return nil
}

// --- Business ---

func (k Keeper) CreatePlan(ctx sdk.Context, creator, fundAddr, title, desc string, role types.PresidentRoleType, regionId string, planJSON []byte) (uint64, error) {
	if err := k.MustBePresident(ctx, creator, role, regionId); err != nil {
		return 0, err
	}
	var tmp any
	if err := json.Unmarshal(planJSON, &tmp); err != nil {
		return 0, fmt.Errorf("invalid plan_json: %w", err)
	}

	id := k.nextPlanId(ctx)
	sp := types.StoredFundPlan{
		Id:              id,
		FundAddress:     fundAddr,
		Title:           title,
		Description:     desc,
		Creator:         creator,
		Status:          types.PLAN_STATUS_DRAFT,
		GovProposalId:   0,
		CreatedAtHeight: ctx.BlockHeight(),
		PlanJSON:        planJSON,
	}
	k.SetPlan(ctx, sp)
	return id, nil
}

func (k Keeper) MarkSubmitted(ctx sdk.Context, planId uint64, proposalId uint64) error {
	p, ok := k.GetPlan(ctx, planId)
	if !ok {
		return fmt.Errorf("plan not found")
	}
	if p.Status != types.PLAN_STATUS_DRAFT {
		return fmt.Errorf("plan status must be DRAFT")
	}
	p.Status = types.PLAN_STATUS_SUBMITTED
	p.GovProposalId = proposalId
	k.SetPlan(ctx, p)
	return nil
}

func (k Keeper) MarkRejectedByProposalId(ctx sdk.Context, proposalId uint64) {
	it := storetypes.KVStorePrefixIterator(k.store(ctx), types.PlanKeyPrefix)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		var p types.StoredFundPlan
		mustUnmarshalJSON(it.Value(), &p)
		if p.GovProposalId == proposalId && p.Status == types.PLAN_STATUS_SUBMITTED {
			p.Status = types.PLAN_STATUS_REJECTED
			k.SetPlan(ctx, p)
			return
		}
	}
}

func mustMarshalJSON[T any](v T) []byte {
	bz, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bz
}

func mustUnmarshalJSON[T any](bz []byte, v *T) {
	if err := json.Unmarshal(bz, v); err != nil {
		panic(err)
	}
}
