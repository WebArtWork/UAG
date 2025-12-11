package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = (*MsgExecuteFundPlan)(nil)

func NewMsgExecuteFundPlan(authority string, plan *FundPlan) *MsgExecuteFundPlan {
	return &MsgExecuteFundPlan{Authority: authority, Plan: plan}
}

func (msg *MsgExecuteFundPlan) Route() string { return RouterKey }

func (msg *MsgExecuteFundPlan) Type() string { return "ExecuteFundPlan" }

func (msg *MsgExecuteFundPlan) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg *MsgExecuteFundPlan) GetSignBytes() []byte {
	bz := getModuleCodec().MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgExecuteFundPlan) ValidateBasic() error {
	if msg.Plan == nil {
		return fmt.Errorf("plan is required")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return fmt.Errorf("invalid authority: %w", err)
	}
	if msg.Plan.FundAddress == "" {
		return fmt.Errorf("fund address required")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Plan.FundAddress); err != nil {
		return fmt.Errorf("invalid fund address: %w", err)
	}
	for _, d := range msg.Plan.Delegations {
		if d == nil || d.Amount == nil {
			return fmt.Errorf("delegation entry invalid")
		}
		if d.Amount.Denom != BaseDenom {
			return fmt.Errorf("invalid delegation denom: %s", d.Amount.Denom)
		}
	}
	for _, p := range msg.Plan.Payouts {
		if p == nil || p.Amount == nil {
			return fmt.Errorf("payout entry invalid")
		}
		if p.Amount.Denom != BaseDenom {
			return fmt.Errorf("invalid payout denom: %s", p.Amount.Denom)
		}
		if _, err := sdk.AccAddressFromBech32(p.RecipientAddress); err != nil {
			return fmt.Errorf("invalid recipient: %w", err)
		}
	}
	return nil
}
