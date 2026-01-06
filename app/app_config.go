package app

import (
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	circuitmodulev1 "cosmossdk.io/api/cosmos/circuit/module/v1"
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

	citizenmodulev1 "uagd/x/citizen/types"
	fundmodulev1 "uagd/x/fund/types"
	growthmodulev1 "uagd/x/growth/types"
	ugovmodulev1 "uagd/x/ugov/types"

	"cosmossdk.io/core/appconfig"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/runtime"

	"google.golang.org/protobuf/types/known/durationpb"
)

var AppConfig = appconfig.Compose(&runtime.AppConfig{
	Modules: []*runtime.Module{
		runtime.NewModule(
			"auth",
			appconfig.WrapAny(&authmodulev1.Module{}),
		),
		runtime.NewModule(
			"bank",
			appconfig.WrapAny(&bankmodulev1.Module{}),
		),
		runtime.NewModule(
			"circuit",
			appconfig.WrapAny(&circuitmodulev1.Module{}),
		),
		runtime.NewModule(
			"distribution",
			appconfig.WrapAny(&distributionmodulev1.Module{}),
		),
		runtime.NewModule(
			"epochs",
			appconfig.WrapAny(&epochsmodulev1.Module{}),
		),
		runtime.NewModule(
			"evidence",
			appconfig.WrapAny(&evidencemodulev1.Module{}),
		),
		runtime.NewModule(
			"feegrant",
			appconfig.WrapAny(&feegrantmodulev1.Module{}),
		),
		runtime.NewModule(
			"genutil",
			appconfig.WrapAny(&genutilmodulev1.Module{}),
		),
		runtime.NewModule(
			"gov",
			appconfig.WrapAny(&govmodulev1.Module{
				MaxMetadataLen: 255,
				Authority:      "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
				ExpeditedVotingPeriod: &durationpb.Duration{
					Seconds: 172800,
				},
				ExpeditedThreshold:         "0.667000000000000000",
				ExpeditedMinDeposit:        "50000000",
				ExpeditedMinDepositRatio:   "0.500000000000000000",
				MinInitialDepositRatio:     "0.000000000000000000",
				BurnVoteQuorum:             false,
				BurnProposalDepositPrevote: false,
				BurnVoteVeto:               true,
				VotingPeriod:               &durationpb.Duration{Seconds: 172800},
				MaxDepositPeriod:           &durationpb.Duration{Seconds: 172800},
				MinDeposit:                 "100000000",
				Quorum:                     "0.334000000000000000",
				Threshold:                  "0.500000000000000000",
				VetoThreshold:              "0.334000000000000000",
			}),
		),
		runtime.NewModule(
			"group",
			appconfig.WrapAny(&groupmodulev1.Module{
				MaxExecutionPeriod: &durationpb.Duration{
					Seconds: 1209600,
				},
				MaxMetadataLen: 255,
			}),
		),
		runtime.NewModule(
			"nft",
			appconfig.WrapAny(&nftmodulev1.Module{}),
		),
		runtime.NewModule(
			"params",
			appconfig.WrapAny(&paramsmodulev1.Module{}),
		),
		runtime.NewModule(
			"slashing",
			appconfig.WrapAny(&slashingmodulev1.Module{
				SignedBlocksWindow:      10000,
				MinSignedPerWindow:      "0.500000000000000000",
				DowntimeJailDuration:    &durationpb.Duration{Seconds: 600},
				SlashFractionDoubleSign: "0.050000000000000000",
				SlashFractionDowntime:   "0.010000000000000000",
			}),
		),
		runtime.NewModule(
			"staking",
			appconfig.WrapAny(&stakingmodulev1.Module{
				UnbondingTime: &durationpb.Duration{
					Seconds: 1814400,
				},
				MaxValidators:     100,
				MaxEntries:        7,
				HistoricalEntries: 10000,
				BondDenom:         "uuag",
				MinCommissionRate: "0.000000000000000000",
			}),
		),
		runtime.NewModule(
			"tx",
			appconfig.WrapAny(&txconfigv1.Config{}),
		),
		runtime.NewModule(
			"upgrade",
			appconfig.WrapAny(&upgrademodulev1.Module{}),
		),

		// uagd modules
		runtime.NewModule(
			"citizen",
			appconfig.WrapAny(&citizenmodulev1.Module{}),
		),
		runtime.NewModule(
			"fund",
			appconfig.WrapAny(&fundmodulev1.Module{}),
		),
		runtime.NewModule(
			"growth",
			appconfig.WrapAny(&growthmodulev1.Module{}),
		),
		runtime.NewModule(
			"ugov",
			appconfig.WrapAny(&ugovmodulev1.Module{}),
		),
	},
}, log.NewNopLogger())
