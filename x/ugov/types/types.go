package types

import "encoding/json"

type President struct {
	RoleType    PresidentRoleType `json:"role_type" yaml:"role_type"`
	RegionId    string            `json:"region_id" yaml:"region_id"`
	Address     string            `json:"address" yaml:"address"`
	Active      bool              `json:"active" yaml:"active"`
	SinceHeight int64             `json:"since_height" yaml:"since_height"`
}

// StoredFundPlan keeps a snapshot of the plan at creation time.
//
// NOTE: this scaffold stores PlanJSON to avoid hard proto coupling.
// Replace PlanJSON with fundtypes.FundPlan (proto) once your x/fund proto is in place.
type StoredFundPlan struct {
	Id               uint64          `json:"id" yaml:"id"`
	FundAddress      string          `json:"fund_address" yaml:"fund_address"`
	Title            string          `json:"title" yaml:"title"`
	Description      string          `json:"description" yaml:"description"`
	Creator          string          `json:"creator" yaml:"creator"`
	Status           FundPlanStatus  `json:"status" yaml:"status"`
	GovProposalId    uint64          `json:"gov_proposal_id" yaml:"gov_proposal_id"`
	CreatedAtHeight  int64           `json:"created_at_height" yaml:"created_at_height"`
	ExecutedAtHeight int64           `json:"executed_at_height" yaml:"executed_at_height"`
	PlanJSON         json.RawMessage `json:"plan_json" yaml:"plan_json"`
}

type GenesisState struct {
	Params     Params           `json:"params" yaml:"params"`
	Presidents []President      `json:"presidents" yaml:"presidents"`
	Plans      []StoredFundPlan `json:"plans" yaml:"plans"`
	PlanSeq    uint64           `json:"plan_seq" yaml:"plan_seq"`
}

func DefaultGenesis() GenesisState {
	return GenesisState{
		Params:     DefaultParams(),
		Presidents: []President{},
		Plans:      []StoredFundPlan{},
		PlanSeq:    0,
	}
}
