package types

type PresidentRoleType int32

const (
	PRESIDENT_TYPE_NATIONAL PresidentRoleType = 0
	PRESIDENT_TYPE_REGION   PresidentRoleType = 1
)

type FundPlanStatus int32

const (
	PLAN_STATUS_DRAFT     FundPlanStatus = 0
	PLAN_STATUS_SUBMITTED FundPlanStatus = 1
	PLAN_STATUS_APPROVED  FundPlanStatus = 2
	PLAN_STATUS_REJECTED  FundPlanStatus = 3
	PLAN_STATUS_EXECUTED  FundPlanStatus = 4
)
