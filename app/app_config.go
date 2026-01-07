package app

import (
	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	authzmodulev1 "cosmossdk.io/api/cosmos/authz/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	circuitmodulev1 "cosmossdk.io/api/cosmos/circuit/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	crisismodulev1 "cosmossdk.io/api/cosmos/crisis/module/v1"
	distributionmodulev1 "cosmossdk.io/api/cosmos/distribution/module/v1"
	epochsmodulev1 "cosmossdk.io/api/cosmos/epochs/module/v1"
	evidencemodulev1 "cosmossdk.io/api/cosmos/evidence/module/v1"
	feegrantmodulev1 "cosmossdk.io/api/cosmos/feegrant/module/v1"
	genutilmodulev1 "cosmossdk.io/api/cosmos/genutil/module/v1"
	govmodulev1 "cosmossdk.io/api/cosmos/gov/module/v1"
	groupmodulev1 "cosmossdk.io/api/cosmos/group/module/v1"
	nftmodulev1 "cosmossdk.io/api/cosmos/nft/module/v1"
	paramsmodulev1 "cosmossdk.io/api/cosmos/params/module/v1"
	slashingmodulev1 "cosmossdk.io/api/cosmos/slashing/module/v1"
	stakingmodulev1 "cosmossdk.io/api/cosmos/staking/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
	upgrademodulev1 "cosmossdk.io/api/cosmos/upgrade/module/v1"

	"google.golang.org/protobuf/types/known/durationpb"

	citizenmodulev1 "uagd/x/citizen/types"
	growthmodulev1 "uagd/x/growth/types"

	"cosmossdk.io/depinject/appconfig"
)

var (
	// Module execution ordering.
	beginBlockers = []string{
		"upgrade",
		"staking",
		"slashing",
		"distribution",
		"gov",
		"auth",
		"bank",
		"params",
		"feegrant",
		"authz",
		"circuit",
		"citizen",
		"growth",
	}
	endBlockers = []string{
		"gov",
		"staking",
		"auth",
		"bank",
		"distribution",
		"slashing",
		"params",
		"circuit",
		"citizen",
		"growth",
	}
	initGenesis = []string{
		"auth",
		"bank",
		"distribution",
		"staking",
		"slashing",
		"gov",
		"genutil",
		"params",
		"feegrant",
		"authz",
		"circuit",
		"epochs",
		"evidence",
		"citizen",
		"growth",
		"wasm",
	}
)

// AppConfig is the depinject wiring configuration for the application.
var AppConfig = appconfig.Compose(&appv1alpha1.Config{
	Modules: []*appv1alpha1.ModuleConfig{
		{
			Name: "runtime",
			Config: appconfig.WrapAny(&runtimev1alpha1.Module{
				AppName:           Name,
				BeginBlockers:     beginBlockers,
				EndBlockers:       endBlockers,
				InitGenesis:       initGenesis,
				OverrideStoreKeys: []*runtimev1alpha1.StoreKeyConfig{{ModuleName: "auth", KvStoreKey: "acc"}},
			}),
		},
		{
			Name:   "auth",
			Config: appconfig.WrapAny(&authmodulev1.Module{Bech32Prefix: AccountAddressPrefix}),
		},
		{
			Name:   "authz",
			Config: appconfig.WrapAny(&authzmodulev1.Module{}),
		},
		{
			Name:   "bank",
			Config: appconfig.WrapAny(&bankmodulev1.Module{}),
		},
		{
			Name:   "circuit",
			Config: appconfig.WrapAny(&circuitmodulev1.Module{}),
		},
		{
			Name:   "consensus",
			Config: appconfig.WrapAny(&consensusmodulev1.Module{}),
		},
		{
			Name:   "crisis",
			Config: appconfig.WrapAny(&crisismodulev1.Module{}),
		},
		{
			Name:   "distribution",
			Config: appconfig.WrapAny(&distributionmodulev1.Module{}),
		},
		{
			Name:   "epochs",
			Config: appconfig.WrapAny(&epochsmodulev1.Module{}),
		},
		{
			Name:   "evidence",
			Config: appconfig.WrapAny(&evidencemodulev1.Module{}),
		},
		{
			Name:   "feegrant",
			Config: appconfig.WrapAny(&feegrantmodulev1.Module{}),
		},
		{
			Name:   "genutil",
			Config: appconfig.WrapAny(&genutilmodulev1.Module{}),
		},
		{
			Name: "gov",
			Config: appconfig.WrapAny(&govmodulev1.Module{
				MaxMetadataLen: 255,
				Authority:      "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
			}),
		},
		{
			Name: "group",
			Config: appconfig.WrapAny(&groupmodulev1.Module{
				MaxExecutionPeriod: &durationpb.Duration{
					Seconds: 1209600,
				},
				MaxMetadataLen: 255,
			}),
		},
		{
			Name:   "nft",
			Config: appconfig.WrapAny(&nftmodulev1.Module{}),
		},
		{
			Name:   "params",
			Config: appconfig.WrapAny(&paramsmodulev1.Module{}),
		},
		{
			Name:   "slashing",
			Config: appconfig.WrapAny(&slashingmodulev1.Module{}),
		},
		{
			Name: "staking",
			Config: appconfig.WrapAny(&stakingmodulev1.Module{
				Bech32PrefixValidator: AccountAddressPrefix + "valoper",
				Bech32PrefixConsensus: AccountAddressPrefix + "valcons",
			}),
		},
		{
			Name:   "tx",
			Config: appconfig.WrapAny(&txconfigv1.Config{}),
		},
		{
			Name:   "upgrade",
			Config: appconfig.WrapAny(&upgrademodulev1.Module{}),
		},
		// uagd modules
		{
			Name:   "citizen",
			Config: appconfig.WrapAny(&citizenmodulev1.Module{}),
		},
		{
			Name:   "growth",
			Config: appconfig.WrapAny(&growthmodulev1.Module{}),
		},
	},
})
