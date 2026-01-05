package app

import (
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	authzmodulev1 "cosmossdk.io/api/cosmos/authz/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	circuitmodulev1 "cosmossdk.io/api/cosmos/circuit/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	distributionmodulev1 "cosmossdk.io/api/cosmos/distribution/module/v1"
	evidencemodulev1 "cosmossdk.io/api/cosmos/evidence/module/v1"
	feegrantmodulev1 "cosmossdk.io/api/cosmos/feegrant/module/v1"
	genutilmodulev1 "cosmossdk.io/api/cosmos/genutil/module/v1"
	govmodulev1 "cosmossdk.io/api/cosmos/gov/module/v1"
	groupmodulev1 "cosmossdk.io/api/cosmos/group/module/v1"
	mintmodulev1 "cosmossdk.io/api/cosmos/mint/module/v1"
	nftmodulev1 "cosmossdk.io/api/cosmos/nft/module/v1"
	paramsmodulev1 "cosmossdk.io/api/cosmos/params/module/v1"
	slashingmodulev1 "cosmossdk.io/api/cosmos/slashing/module/v1"
	stakingmodulev1 "cosmossdk.io/api/cosmos/staking/module/v1"
	upgrademodulev1 "cosmossdk.io/api/cosmos/upgrade/module/v1"
	"cosmossdk.io/core/appconfig"
	depinject "cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/runtime"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"google.golang.org/protobuf/types/known/durationpb"

	wasmmodulev1 "github.com/CosmWasm/wasmd/api/cosmwasm/wasm/module/v1"
	_ "github.com/CosmWasm/wasmd/x/wasm" // registers depinject wiring (init)
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	citizenmodulev1 "uagd/x/citizen/types"
	fundmodulev1 "uagd/x/fund/types"
	growthmodulev1 "uagd/x/growth/types"
	uagdmodulev1 "uagd/x/uagd/types"
	ugovmodulev1 "uagd/x/ugov/types"
)

func baseAppConfig() depinject.Config {
	govAuthority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	return appconfig.Compose(
		&runtime.AppConfig{
			Modules: []*runtime.Module{
				runtime.NewModule("auth", appconfig.WrapAny(&authmodulev1.Module{})),
				runtime.NewModule("authz", appconfig.WrapAny(&authzmodulev1.Module{Authority: govAuthority})),
				runtime.NewModule("bank", appconfig.WrapAny(&bankmodulev1.Module{})),
				runtime.NewModule("circuit", appconfig.WrapAny(&circuitmodulev1.Module{})),
				runtime.NewModule(
					"consensus",
					appconfig.WrapAny(&consensusmodulev1.Module{
						Authority: govAuthority,
					}),
				),
				runtime.NewModule("distribution", appconfig.WrapAny(&distributionmodulev1.Module{})),
				runtime.NewModule("evidence", appconfig.WrapAny(&evidencemodulev1.Module{})),
				runtime.NewModule("feegrant", appconfig.WrapAny(&feegrantmodulev1.Module{})),
				runtime.NewModule(
					"genutil",
					appconfig.WrapAny(&genutilmodulev1.Module{
						BypassGenesisInvariants: false,
					}),
				),
				runtime.NewModule(
					"gov",
					appconfig.WrapAny(&govmodulev1.Module{
						MaxMetadataLen:                        512,
						Authority:                             govAuthority,
						ExpeditedProposalDeposit:              "5000000000",
						MinInitialDepositRatio:                "0.25",
						BurnProposalDepositPrevote:            false,
						BurnVoteQuorum:                        false,
						BurnVoteVeto:                          true,
						MinDeposit:                            "100000000",
						ExpeditedMinDeposit:                   "1000000000",
						MinDepositRatio:                       "0.01",
						MaxDepositPeriod:                      durationpb.New(172800000000000),
						VotingPeriod:                          durationpb.New(172800000000000),
						ExpeditedVotingPeriod:                 durationpb.New(86400000000000),
						Quorum:                                "0.40",
						Threshold:                             "0.50",
						VetoThreshold:                         "0.334",
						ExpeditedThreshold:                    "0.67",
						ProposalCancelRatio:                   "0.50",
						ProposalCancelDest:                    govAuthority,
						MinProposalDepositTokens:              "100000000",
						MinExpeditedProposalDepositTokens:     "1000000000",
						MinProposalDepositRatio:               "0.01",
						MinExpeditedProposalDepositRatio:      "0.10",
						BurnProposalDeposit:                   false,
						BurnVote:                              false,
						BurnVoteVetoDeposit:                   true,
						MinVotingPeriod:                       durationpb.New(86400000000000),
						MinExpeditedVotingPeriod:              durationpb.New(3600000000000),
						BurnProposalCancelDeposit:             false,
						MaxVotingPeriod:                       durationpb.New(259200000000000),
						MaxExpeditedVotingPeriod:              durationpb.New(129600000000000),
						MinDepositPeriod:                      durationpb.New(86400000000000),
						MinExpeditedDepositPeriod:             durationpb.New(43200000000000),
						MaxDepositPeriodTokens:                "10000000000",
						MaxExpeditedDepositPeriodTokens:       "50000000000",
						MaxDepositPeriodRatio:                 "0.10",
						MaxExpeditedDepositPeriodRatio:        "0.20",
						MinDepositTokens:                      "100000000",
						MinExpeditedDepositTokens:             "1000000000",
						ProposalDepositMinimumRatio:           "0.01",
						ProposalDepositMinimumTokens:          "100000000",
						ProposalDepositMinimumExpeditedRatio:  "0.10",
						ProposalDepositMinimumExpeditedTokens: "1000000000",
					}),
				),
				runtime.NewModule(
					"group",
					appconfig.WrapAny(&groupmodulev1.Module{
						MaxMetadataLen:           255,
						MaxExecutionPeriod:       durationpb.New(1209600000000000),
						MaxProposalTitleLen:      255,
						MaxProposalSummaryLen:    10200,
						MaxProposalMetadataLen:   255,
						MaxProposalMsgs:          100,
						MaxDecisionPolicyWindows: durationpb.New(1209600000000000),
					}),
				),
				runtime.NewModule(
					"mint",
					appconfig.WrapAny(&mintmodulev1.Module{
						MintDenom:           "uuag",
						InflationRateChange: "0.130000000000000000",
						InflationMax:        "0.200000000000000000",
						InflationMin:        "0.070000000000000000",
						GoalBonded:          "0.670000000000000000",
						BlocksPerYear:       4360000,
					}),
				),
				runtime.NewModule(
					"nft",
					appconfig.WrapAny(&nftmodulev1.Module{
						ClassCreationFee: "5000000",
					}),
				),
				runtime.NewModule("params", appconfig.WrapAny(&paramsmodulev1.Module{})),
				runtime.NewModule("slashing", appconfig.WrapAny(&slashingmodulev1.Module{})),
				runtime.NewModule("staking", appconfig.WrapAny(&stakingmodulev1.Module{})),
				runtime.NewModule("upgrade", appconfig.WrapAny(&upgrademodulev1.Module{})),

				// UAGD modules (use module config types from x/<module>/types, not api/*/module/v1)
				runtime.NewModule("uagd", appconfig.WrapAny(&uagdmodulev1.Module{Authority: govAuthority})),
				runtime.NewModule("citizen", appconfig.WrapAny(&citizenmodulev1.Module{Authority: govAuthority})),
				runtime.NewModule("fund", appconfig.WrapAny(&fundmodulev1.Module{Authority: govAuthority})),
				runtime.NewModule("growth", appconfig.WrapAny(&growthmodulev1.Module{Authority: govAuthority})),
				runtime.NewModule("ugov", appconfig.WrapAny(&ugovmodulev1.Module{Authority: govAuthority})),

				// wasm (CosmWasm) module defaults
				runtime.NewModule(
					wasmtypes.ModuleName,
					appconfig.WrapAny(&wasmmodulev1.Module{
						UploadPermission: wasmtypes.DefaultUploadAccess,
					}),
				),
			},
		},
		log.NewNopLogger(),
	)
}
