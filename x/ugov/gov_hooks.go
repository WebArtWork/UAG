package ugov

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"uagd/x/ugov/keeper"
)

type planMarker interface {
	MarkRejectedByProposalId(ctx sdk.Context, proposalId uint64)
}

type proposalGetter interface {
	GetProposal(ctx context.Context, proposalID uint64) (v1.Proposal, error)
}

type GovHooks struct {
	k         planMarker
	proposals proposalGetter
}

func NewGovHooks(k keeper.Keeper, govKeeper *govkeeper.Keeper) GovHooks {
	return GovHooks{k: k, proposals: govProposalGetter{govKeeper}}
}

type govProposalGetter struct{ govKeeper *govkeeper.Keeper }

func (g govProposalGetter) GetProposal(ctx context.Context, proposalID uint64) (v1.Proposal, error) {
	if g.govKeeper == nil {
		return v1.Proposal{}, fmt.Errorf("gov keeper not configured")
	}
	return g.govKeeper.Proposals.Get(ctx, proposalID)
}

var _ govtypes.GovHooks = GovHooks{}

func (h GovHooks) AfterProposalSubmission(context.Context, uint64) error { return nil }

func (h GovHooks) AfterProposalDeposit(context.Context, uint64, sdk.AccAddress) error { return nil }

func (h GovHooks) AfterProposalFailedMinDeposit(context.Context, uint64) error { return nil }

func (h GovHooks) AfterProposalVote(context.Context, uint64, sdk.AccAddress) error { return nil }

func (h GovHooks) AfterProposalVotingPeriodEnded(ctx context.Context, proposalID uint64) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if h.proposals == nil {
		return fmt.Errorf("proposal getter not configured")
	}

	proposal, err := h.proposals.GetProposal(ctx, proposalID)
	if err != nil {
		return err
	}

	switch proposal.Status {
	case v1.StatusRejected, v1.StatusFailed:
		h.k.MarkRejectedByProposalId(sdkCtx, proposalID)
	}
	return nil
}
