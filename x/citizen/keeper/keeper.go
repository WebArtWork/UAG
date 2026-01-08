package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"uag/x/citizen/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	types.UnimplementedQueryServer

	Schema                 collections.Schema
	ParamsStore            collections.Item[types.Params]
	RegionByAddressStore   collections.Map[[]byte, string]
	AddressByRegionIndexes collections.KeySet[collections.Pair[string, []byte]]
}

func NewKeeper(storeService corestore.KVStoreService, cdc codec.Codec, addressCodec address.Codec) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		ParamsStore:  collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		RegionByAddressStore: collections.NewMap(
			sb,
			collections.NewPrefix(types.RegionByAddressPrefix),
			"region_by_address",
			collections.BytesKey,
			collections.StringValue,
		),
		AddressByRegionIndexes: collections.NewKeySet(
			sb,
			collections.NewPrefix(types.AddressByRegionPrefix),
			"address_by_region",
			collections.PairKeyCodec(collections.StringKey, collections.BytesKey),
		),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}
	return k.ParamsStore.Set(ctx, params)
}

func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	return k.ParamsStore.Get(ctx)
}

func (k Keeper) IsRegistrar(ctx context.Context, addr sdk.AccAddress) bool {
	params, err := k.GetParams(ctx)
	if err != nil {
		return false
	}
	addrStr := addr.String()
	for _, r := range params.Registrars {
		if r == addrStr {
			return true
		}
	}
	return false
}

func (k Keeper) SetRegion(ctx context.Context, addr sdk.AccAddress, regionID string) error {
	if regionID == "" {
		return fmt.Errorf("region id required")
	}
	addrBytes, err := k.addressCodec.StringToBytes(addr.String())
	if err != nil {
		return err
	}
	oldRegion, err := k.RegionByAddressStore.Get(ctx, addrBytes)
	if err == nil {
		_ = k.AddressByRegionIndexes.Remove(ctx, collections.Join(oldRegion, addrBytes))
	} else if !errors.Is(err, collections.ErrNotFound) {
		return err
	}
	if err := k.RegionByAddressStore.Set(ctx, addrBytes, regionID); err != nil {
		return err
	}
	return k.AddressByRegionIndexes.Set(ctx, collections.Join(regionID, addrBytes))
}

func (k Keeper) GetRegion(ctx context.Context, addr sdk.AccAddress) (string, bool) {
	addrBytes, err := k.addressCodec.StringToBytes(addr.String())
	if err != nil {
		return "", false
	}
	region, err := k.RegionByAddressStore.Get(ctx, addrBytes)
	if err != nil {
		return "", false
	}
	return region, true
}

func (k Keeper) GetRegionByString(ctx context.Context, addrStr string) (string, bool) {
	addr, err := k.addressCodec.StringToBytes(addrStr)
	if err != nil {
		return "", false
	}
	region, err := k.RegionByAddressStore.Get(ctx, addr)
	if err != nil {
		return "", false
	}
	return region, true
}

func (k Keeper) DeleteRegion(ctx context.Context, addr sdk.AccAddress) error {
	addrBytes, err := k.addressCodec.StringToBytes(addr.String())
	if err != nil {
		return err
	}
	region, err := k.RegionByAddressStore.Get(ctx, addrBytes)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil
		}
		return err
	}
	if err := k.RegionByAddressStore.Remove(ctx, addrBytes); err != nil {
		return err
	}
	_ = k.AddressByRegionIndexes.Remove(ctx, collections.Join(region, addrBytes))
	return nil
}

func (k Keeper) GetAddressesByRegion(ctx context.Context, regionID string) []string {
	var addresses []string
	rng := collections.NewPrefixedPairRange[string, []byte](regionID)
	iterator, err := k.AddressByRegionIndexes.Iterate(ctx, rng)
	if err != nil {
		return addresses
	}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key, err := iterator.Key()
		if err != nil {
			continue
		}
		if key.K1() != regionID {
			continue
		}
		addrStr, err := k.addressCodec.BytesToString(key.K2())
		if err != nil {
			continue
		}
		addresses = append(addresses, addrStr)
	}
	return addresses
}
