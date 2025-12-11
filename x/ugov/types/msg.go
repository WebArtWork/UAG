package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgSetPresident           = "set_president"
	TypeMsgCreateFundPlan         = "create_fund_plan"
	TypeMsgSubmitFundPlanProposal = "submit_fund_plan_as_proposal"
	TypeMsgExecuteFundPlan        = "execute_fund_plan"
)

type MsgSetPresident struct {
	Authority string            `json:"authority" yaml:"authority"`
	RoleType  PresidentRoleType `json:"role_type" yaml:"role_type"`
	RegionId  string            `json:"region_id" yaml:"region_id"`
	Address   string            `json:"address" yaml:"address"`
}

func (m MsgSetPresident) Route() string { return RouterKey }
func (m MsgSetPresident) Type() string  { return TypeMsgSetPresident }
func (m MsgSetPresident) GetSigners() []sdk.AccAddress {
	a, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{a}
}
func (m MsgSetPresident) ValidateBasic() error {
	if m.Authority == "" {
		return fmt.Errorf("authority required")
	}
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return err
	}
	if m.Address == "" {
		return fmt.Errorf("address required")
	}
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return err
	}
	if m.RoleType == PRESIDENT_TYPE_REGION && m.RegionId == "" {
		return fmt.Errorf("region_id required for region president")
	}
	if m.RoleType == PRESIDENT_TYPE_NATIONAL && m.RegionId != "" {
		return fmt.Errorf("region_id must be empty for national president")
	}
	return nil
}

type MsgCreateFundPlan struct {
	Creator     string `json:"creator" yaml:"creator"`
	FundAddress string `json:"fund_address" yaml:"fund_address"`
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	PlanJSON    []byte `json:"plan_json" yaml:"plan_json"`
}

func (m MsgCreateFundPlan) Route() string { return RouterKey }
func (m MsgCreateFundPlan) Type() string  { return TypeMsgCreateFundPlan }
func (m MsgCreateFundPlan) GetSigners() []sdk.AccAddress {
	a, _ := sdk.AccAddressFromBech32(m.Creator)
	return []sdk.AccAddress{a}
}
func (m MsgCreateFundPlan) ValidateBasic() error {
	if m.Creator == "" {
		return fmt.Errorf("creator required")
	}
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return err
	}
	if m.FundAddress == "" {
		return fmt.Errorf("fund_address required")
	}
	if m.Title == "" {
		return fmt.Errorf("title required")
	}
	if len(m.PlanJSON) == 0 {
		return fmt.Errorf("plan_json required")
	}
	return nil
}

type MsgSubmitFundPlanAsProposal struct {
	Creator string `json:"creator" yaml:"creator"`
	PlanId  uint64 `json:"plan_id" yaml:"plan_id"`
	Title   string `json:"title" yaml:"title"`
	Summary string `json:"summary" yaml:"summary"`
}

func (m MsgSubmitFundPlanAsProposal) Route() string { return RouterKey }
func (m MsgSubmitFundPlanAsProposal) Type() string  { return TypeMsgSubmitFundPlanProposal }
func (m MsgSubmitFundPlanAsProposal) GetSigners() []sdk.AccAddress {
	a, _ := sdk.AccAddressFromBech32(m.Creator)
	return []sdk.AccAddress{a}
}
func (m MsgSubmitFundPlanAsProposal) ValidateBasic() error {
	if m.Creator == "" {
		return fmt.Errorf("creator required")
	}
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return err
	}
	if m.PlanId == 0 {
		return fmt.Errorf("plan_id required")
	}
	if m.Title == "" {
		return fmt.Errorf("title required")
	}
	if m.Summary == "" {
		return fmt.Errorf("summary required")
	}
	return nil
}

// MsgExecuteFundPlan is executed by x/gov when a proposal passes (gov v1 message proposals).
type MsgExecuteFundPlan struct {
	Authority string `json:"authority" yaml:"authority"`
	PlanId    uint64 `json:"plan_id" yaml:"plan_id"`
}

func (m MsgExecuteFundPlan) Route() string { return RouterKey }
func (m MsgExecuteFundPlan) Type() string  { return TypeMsgExecuteFundPlan }
func (m MsgExecuteFundPlan) GetSigners() []sdk.AccAddress {
	a, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{a}
}
func (m MsgExecuteFundPlan) ValidateBasic() error {
	if m.Authority == "" {
		return fmt.Errorf("authority required")
	}
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return err
	}
	if m.PlanId == 0 {
		return fmt.Errorf("plan_id required")
	}
	return nil
}
