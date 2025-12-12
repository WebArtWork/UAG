package ugov

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"uagd/x/ugov/keeper"
)

type GovHooks struct{ k keeper.Keeper }

func NewGovHooks(k keeper.Keeper) GovHooks { return GovHooks{k: k} }

var _ govtypes.GovHooks = GovHooks{}

func (h GovHooks) AfterProposalSubmission(context.Context, uint64) error { return nil }

func (h GovHooks) AfterProposalDeposit(context.Context, uint64, sdk.AccAddress) error { return nil }

func (h GovHooks) AfterProposalFailedMinDeposit(context.Context, uint64) error { return nil }

func (h GovHooks) AfterProposalVote(context.Context, uint64, sdk.AccAddress) error { return nil }

func (h GovHooks) AfterProposalVotingPeriodEnded(ctx context.Context, proposalID uint64) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	h.k.MarkRejectedByProposalId(sdkCtx, proposalID)
	return nil
}
