# CLI workflows

## 1) Set national president (admin only)

Tx: `MsgSetPresident`
- authority = ugov.params.admin
- role_type = NATIONAL
- region_id = ""

## 2) Create a plan (president)

Tx: `MsgCreateFundPlan`
- creator = president
- fund_address = fund module fund account address
- plan_json = JSON encoding of your intended fund plan (temporary scaffold)

Result: plan stored as DRAFT.

## 3) Submit governance proposal (gov v1) to execute the plan

Create a proposal with message:

- `ugov.MsgExecuteFundPlan{ authority: <gov module address>, plan_id: <id> }`

## 4) Vote & pass

When passed:
- gov executes msg
- ugov verifies authority == gov authority
- ugov calls x/fund ExecuteFundPlan(...)
- plan status becomes EXECUTED
