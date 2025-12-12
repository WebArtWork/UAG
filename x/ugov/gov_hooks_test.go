package ugov

import (
	"context"
	"testing"

	slog "cosmossdk.io/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

type stubPlanMarker struct{ rejected []uint64 }

type stubProposalGetter struct {
	proposal v1.Proposal
	err      error
}

func (s *stubPlanMarker) MarkRejectedByProposalId(_ sdk.Context, proposalID uint64) {
	s.rejected = append(s.rejected, proposalID)
}

func (s stubProposalGetter) GetProposal(context.Context, uint64) (v1.Proposal, error) {
	return s.proposal, s.err
}

func TestAfterProposalVotingPeriodEndedMarksRejectionOnlyOnFailure(t *testing.T) {
	sdkCtx := sdk.NewContext(nil, tmproto.Header{}, false, slog.NewNopLogger()).WithContext(context.Background())
	ctx := sdk.WrapSDKContext(sdkCtx)

	failingHook := GovHooks{
		k:         &stubPlanMarker{},
		proposals: stubProposalGetter{proposal: v1.Proposal{Status: v1.StatusRejected}},
	}

	if err := failingHook.AfterProposalVotingPeriodEnded(ctx, 1); err != nil {
		t.Fatalf("hook returned error for failed proposal: %v", err)
	}

	if len(failingHook.k.(*stubPlanMarker).rejected) != 1 {
		t.Fatalf("expected plan marker to be called for failed proposal")
	}

	passingHook := GovHooks{
		k:         &stubPlanMarker{},
		proposals: stubProposalGetter{proposal: v1.Proposal{Status: v1.StatusPassed}},
	}

	if err := passingHook.AfterProposalVotingPeriodEnded(ctx, 2); err != nil {
		t.Fatalf("hook returned error for passed proposal: %v", err)
	}

	if len(passingHook.k.(*stubPlanMarker).rejected) != 0 {
		t.Fatalf("expected plan marker not to be called for passed proposal")
	}
}
