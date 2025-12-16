package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m *MsgCreatePlan) ValidateBasic() error {
	if m == nil {
		return fmt.Errorf("empty message")
	}
	if m.Creator == "" {
		return fmt.Errorf("creator required")
	}
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return fmt.Errorf("invalid creator: %w", err)
	}
	if m.FundAddress == "" {
		return fmt.Errorf("fund_address required")
	}
	// fund_address is AddressString, so check bech32
	if _, err := sdk.AccAddressFromBech32(m.FundAddress); err != nil {
		return fmt.Errorf("invalid fund_address: %w", err)
	}
	if m.Title == "" {
		return fmt.Errorf("title required")
	}
	// m.Position is non-nullable in proto, but still validate basics here if needed.
	return nil
}

func (m *MsgUpdatePlan) ValidateBasic() error {
	if m == nil {
		return fmt.Errorf("empty message")
	}
	if m.Creator == "" {
		return fmt.Errorf("creator required")
	}
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return fmt.Errorf("invalid creator: %w", err)
	}
	if m.Id == 0 {
		return fmt.Errorf("id required")
	}
	if m.Title == "" {
		return fmt.Errorf("title required")
	}
	return nil
}

func (m *MsgSubmitPlan) ValidateBasic() error {
	if m == nil {
		return fmt.Errorf("empty message")
	}
	if m.Creator == "" {
		return fmt.Errorf("creator required")
	}
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return fmt.Errorf("invalid creator: %w", err)
	}
	if m.Id == 0 {
		return fmt.Errorf("id required")
	}
	return nil
}

func (m *MsgExecuteFundPosition) ValidateBasic() error {
	if m == nil {
		return fmt.Errorf("empty message")
	}
	if m.Authority == "" {
		return fmt.Errorf("authority required")
	}
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return fmt.Errorf("invalid authority: %w", err)
	}
	if m.PlanId == 0 {
		return fmt.Errorf("plan_id required")
	}
	return nil
}
