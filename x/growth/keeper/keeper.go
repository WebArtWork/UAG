package keeper

import (
	"context"

	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"

	"uagd/x/growth/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	logger       log.Logger
}

func NewKeeper(storeService corestore.KVStoreService, cdc codec.Codec, addressCodec address.Codec) Keeper {
	return Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		logger:       log.NewNopLogger(),
	}
}

func (k Keeper) Logger() log.Logger { return k.logger }

func (k Keeper) SetLogger(l log.Logger) Keeper {
	k.logger = l
	return k
}

// --- keys (core KVStoreService) ---
var (
	keyRegionMetric = []byte{0x01}
	keyGrowthScore  = []byte{0x02}
	keyOccupation   = []byte{0x03}
)

func metricKey(regionID string) []byte {
	b := make([]byte, 0, 1+len(regionID))
	b = append(b, keyRegionMetric...)
	b = append(b, []byte(regionID)...)
	return b
}

func scoreKey(regionID string) []byte {
	b := make([]byte, 0, 1+len(regionID))
	b = append(b, keyGrowthScore...)
	b = append(b, []byte(regionID)...)
	return b
}

func occupationKey(regionID string) []byte {
	b := make([]byte, 0, 1+len(regionID))
	b = append(b, keyOccupation...)
	b = append(b, []byte(regionID)...)
	return b
}

// --- minimal KV storage used by msg_server ---

func (k Keeper) SetRegionMetric(ctx context.Context, m types.RegionMetric) error {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := k.cdc.Marshal(&m)
	if err != nil {
		return err
	}
	return store.Set(metricKey(m.RegionId), bz)
}

func (k Keeper) GetRegionMetric(ctx context.Context, regionID string) (types.RegionMetric, bool, error) {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := store.Get(metricKey(regionID))
	if err != nil {
		return types.RegionMetric{}, false, err
	}
	if bz == nil {
		return types.RegionMetric{}, false, nil
	}

	var out types.RegionMetric
	if err := k.cdc.Unmarshal(bz, &out); err != nil {
		return types.RegionMetric{}, false, err
	}
	return out, true, nil
}

func (k Keeper) SetGrowthScore(ctx context.Context, s types.GrowthScore) error {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := k.cdc.Marshal(&s)
	if err != nil {
		return err
	}
	return store.Set(scoreKey(s.RegionId), bz)
}

func (k Keeper) GetGrowthScore(ctx context.Context, regionID string) (types.GrowthScore, bool, error) {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := store.Get(scoreKey(regionID))
	if err != nil {
		return types.GrowthScore{}, false, err
	}
	if bz == nil {
		return types.GrowthScore{}, false, nil
	}

	var out types.GrowthScore
	if err := k.cdc.Unmarshal(bz, &out); err != nil {
		return types.GrowthScore{}, false, err
	}
	return out, true, nil
}

func (k Keeper) SetOccupation(ctx context.Context, o types.Occupation) error {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := k.cdc.Marshal(&o)
	if err != nil {
		return err
	}
	return store.Set(occupationKey(o.RegionId), bz)
}

func (k Keeper) GetOccupation(ctx context.Context, regionID string) (types.Occupation, bool, error) {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := store.Get(occupationKey(regionID))
	if err != nil {
		return types.Occupation{}, false, err
	}
	if bz == nil {
		return types.Occupation{}, false, nil
	}

	var out types.Occupation
	if err := k.cdc.Unmarshal(bz, &out); err != nil {
		return types.Occupation{}, false, err
	}
	return out, true, nil
}
