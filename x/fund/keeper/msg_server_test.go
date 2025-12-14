package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	testkeeper "uagd/testutil/keeper"
	"uagd/x/fund/keeper"
	"uagd/x/fund/types"
)

func TestExecuteFundPlanDirectDisabled(t *testing.T) {
	k, ctx := testkeeper.FundKeeper(t)
	msgServer := keeper.NewMsgServerImpl(k)

	zeroAddr := sdk.AccAddress(make([]byte, 20))
	msg := &types.MsgExecuteFundPlan{Authority: zeroAddr.String(), Plan: &types.FundPlan{FundAddress: zeroAddr.String()}}
	_, err := msgServer.ExecuteFundPlan(ctx, msg)
	if err == nil || err != types.ErrDirectExecDisabled {
		t.Fatalf("expected direct execution disabled error, got %v", err)
	}
}
