# CLI workflows

## 1) Create a plan (president)

Tx: `MsgCreatePlan`

- creator = president
- fund_address = fund module fund account address
- position = the FundPosition (delegations + payouts)

Result: plan stored as DRAFT.

## 2) Submit governance proposal (gov v1) to execute the plan

Create a proposal with message:

- `ugov.MsgExecuteFundPosition{ authority: <gov module address>, plan_id: <id> }`

## 3) Vote & pass

When passed:

- gov executes msg
- ugov verifies authority == gov authority
- ugov calls x/fund ExecuteFundPosition(...)
- plan status becomes EXECUTED
