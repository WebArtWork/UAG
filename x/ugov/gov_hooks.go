package ugov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"uagd/x/ugov/keeper"
)

type GovHooks struct{ k keeper.Keeper }

func NewGovHooks(k keeper.Keeper) GovHooks { return GovHooks{k: k} }

var _ govtypes.GovHooks = GovHooks{}

func (h GovHooks) AfterProposalSubmission(_ sdk.Context, _ uint64)                 {}
func (h GovHooks) AfterProposalDeposit(_ sdk.Context, _ uint64, _ sdk.AccAddress) {}
func (h GovHooks) AfterProposalVote(_ sdk.Context, _ uint64, _ sdk.AccAddress)    {}

func (h GovHooks) AfterProposalVotingPeriodEnded(ctx sdk.Context, proposalID uint64) {
	h.k.MarkRejectedByProposalId(ctx, proposalID)
}
