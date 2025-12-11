
# Agents in Ukraine Growth (UAG)

This file describes **who the agents are**, **what each of them can and cannot do on-chain**, and **which modules enforce those rules**.

If the code and this file ever disagree, treat this as the **design spec** and open an issue / PR to realign.

---

## 1. System at a glance

- **Token:** UAG (base denom `uuag`, 6 decimals)
- **Total supply:** 370M UAG
- **Inflation:** 7% / year from day one

Supply split:

- **70%** → 27 regional funds (one per oblast)
- **20%** → Ukraine-level fund
- **10%** → projects fund for UAG/WAW products and infrastructure

Rules from the pitch:

- Regional **70%** and Ukraine **20%** can only:
  1. **Delegate coins to validators**
  2. **Pay limited salaries and grants** to people who work on Ukraine Growth
- The **size of these limits** (how much can be delegated, how big salaries and grants can be) depends on **how well Ukraine and each region actually grow over time**.
- Growth is measured by **transparent indicators**: taxes, GDP, exports and other **official statistics** for each region and nationally.
- The **10% projects fund** is used **only** for ecosystem projects chosen by a global community vote.
- On both levels, **presidents are servants, not owners**: the **president of Ukraine** and each **region president** prepare plans, explain them, and their communities vote. If people approve, presidents simply execute the decision and cannot secretly move coins or invent extra payouts.

Extra product layer from the pitch:

- Ukraine Growth will offer **ready-made CRM systems** for businesses.
- Each CRM has an **optional UAG mode**:
  - The company can mark some entries as **critical locked records** and send them to the UAG chain.
  - Once anchored on-chain, **nobody can secretly change or delete** those records − not devops, admins, employees, or owners.
  - If someone later “fixes the database” to hide fraud, the **CRM and the chain no longer match**, and this mismatch is visible.

---

## 2. Economic agents

### 2.1 Protocol-owned funds

These are **on-chain accounts without private owners**. They are controlled only by module logic.

#### 2.1.1 Regional funds (27)

- **Who:** 27 protocol-owned accounts, one per oblast.
- **Holds:** UAG only (`uuag`).
- **Can do:**
  - Delegate coins to validators.
  - Pay **limited** salaries and grants to people working on Ukraine Growth in that region.
- **Cannot do:**
  - Send coins arbitrarily (no random transfers, no “gifts”).
  - Delegate or pay **beyond limits** computed from real-world growth indicators.
- **Limits depend on growth:**
  - `x/growth` turns **regional indicators** (taxes, GDP, exports, official stats) into:
    - Max delegation per region.
    - Max yearly / period budgets for salaries and grants.
- **Controlled by:**
  - `x/fund` − registers these accounts as “regional funds” and enforces their rules.
  - `x/growth` − stores and updates per-region limits based on transparent indicators.
  - `x/ugov` − ensures any action (delegation / payment) follows an **approved regional plan**.

#### 2.1.2 Ukraine fund (national level)

- **Who:** One protocol-owned account for the whole country.
- **Holds:** UAG only.
- **Can do:**
  - Delegate to validators.
  - Pay **limited** salaries and grants to national-level contributors.
- **Cannot do:**
  - Bypass national limits or send coins arbitrarily.
- **Limits depend on growth:**
  - `x/growth` aggregates **national indicators** (taxes, GDP, exports, national stats) and computes:
    - National max delegation.
    - National salary / grant ceilings.
- **Controlled by:**
  - `x/fund` − marks it as “Ukraine fund”.
  - `x/growth` − provides national limits.
  - `x/ugov` − handles national plans and votes.

#### 2.1.3 Projects fund (ecosystem)

- **Who:** One protocol-owned account for UAG/WAW ecosystem and infrastructure.
- **Holds:** UAG only.
- **Can do:**
  - Pay **ecosystem project grants and infra costs**.
- **Cannot do:**
  - Delegate to validators.
  - Pay salaries or grants outside **approved project proposals**.
- **Governance constraint:**
  - Every movement from the projects fund comes only **after a global on-chain vote** by the community (validators + delegators).
- **Controlled by:**
  - `x/fund` − marks it as “projects fund” with special rules.
  - `x/ugov` − project proposals and global votes.

---

### 2.2 Human & social agents

#### 2.2.1 Validators

- **Who:** Node operators that propose and validate blocks.
- **Can do:**
  - Create and manage validator operators, commissions, etc.
  - Receive delegations from:
    - Normal users (delegators)
    - Regional funds
    - Ukraine fund
    - (Optionally) special global staking pools
  - Vote in governance with their stake.
- **Cannot do:**
  - Directly move protocol-owned funds.
  - Bypass limits derived from `x/growth`.
- **Modules:**
  - `x/staking`, `x/slashing`, `x/distribution`, `x/gov`, `x/ugov`.

#### 2.2.2 Delegators (UAG holders)

- **Who:** Normal users holding UAG.
- **Can do:**
  - Send, delegate, undelegate, redelegate, withdraw rewards.
  - Vote in governance (`x/gov`, `x/ugov`).
- **Cannot do:**
  - Directly control protocol-owned funds (only via governance).
- **Modules:**
  - `x/bank`, `x/staking`, `x/gov`, `x/ugov`.

#### 2.2.3 Presidents (national and regional)

We have:

- One **president of Ukraine**.
- 27 **regional presidents** (one per oblast).

They are **executors**, not owners.

- **Can do (always with checks):**
  - Draft **plans** for:
    - Delegating regional or national funds.
    - Paying defined salaries and grants.
    - Supporting specific grants or projects.
  - Submit these plans into `x/ugov` for their communities to vote on.
  - After **approval**, call `x/fund` messages to execute the plan within:
    - Fund type rules.
    - Limits from `x/growth`.
- **Cannot do:**
  - Move a single coin from any fund without:
    - A valid plan.
    - A successful community vote.
  - Invent extra payouts or “side payments” outside the approved plan.
- **Modules:**
  - `x/ugov` − defines who is president and how plans and votes work.
  - `x/fund` − enforces hard safety around actual movements.

#### 2.2.4 Communities

- **Who:**
  - **National community:** UAG stakers (and possibly broader holders, TBD).
  - **Regional communities:** token holders / stakeholders for each oblast.
- **Can do:**
  - Vote on:
    - National plans affecting the Ukraine fund and projects fund.
    - Regional plans affecting their regional fund.
    - Project funding from the projects fund.
  - Replace presidents if governance rules allow (future `x/ugov` logic).
- **Cannot do:**
  - Edit on-chain balances directly; they act via governance.
- **Modules:**
  - `x/gov`, `x/ugov`.

#### 2.2.5 Contributors (salary & grant recipients)

- **Who:** People and teams working on Ukraine Growth (dev, community, operations, etc.).
- **Can do:**
  - Receive salaries or grants from:
    - Their regional fund.
    - The Ukraine fund.
    - The projects fund (for ecosystem work).
- **Cannot do:**
  - Pull money themselves; they only receive payouts defined in approved plans.

#### 2.2.6 Project teams

- **Who:** Teams building UAG/WAW products, infrastructure or ecosystem tools.
- **Can do:**
  - Submit project proposals to access the **projects fund** (and potentially regional/Ukraine funds if design allows).
  - Receive funds once the community approves proposals.
- **Cannot do:**
  - Reconfigure fund rules or bypass governance.

---

## 2.3 Business & CRM agents (critical locked records)

This section reflects the **CRM + UAG mode** from the pitch.

#### 2.3.1 Businesses using UAG-enabled CRM

- **Who:** Companies running UAG-integrated CRM systems.
- **Can do:**
  - Run their CRM in **normal mode** (pure database).
  - Enable **optional UAG mode** for selected entries:
    - Mark some records as **critical locked records**.
    - Send a hash / snapshot of those records to the UAG chain.
  - Later compare:
    - Current CRM data vs.
    - Anchored records on the chain.
- **Cannot do:**
  - Edit or delete records that already live on-chain; they can only update their internal database. Any mismatch with the chain becomes evidence.
- **Chain interaction:**
  - They typically send **anchoring transactions** via:
    - A dedicated anchoring module (future `x/anchor` or similar), or
    - A standard module designed for business integration.

#### 2.3.2 Internal staff: admins, devops, employees

- **Who:** People with privileged access to the company’s infrastructure and CRM.
- **Can do:**
  - Edit the CRM database (subject to company rules).
  - Attempt to manipulate data off-chain (this is the classic fraud threat).
- **Cannot do:**
  - Change already-anchored facts on the UAG chain.
  - Rewrite historical on-chain records, even if they control servers and databases.
- **Security effect:**
  - Any attempt to “fix the database” to hide fraud produces a **mismatch**:
    - CRM says one thing, chain says another.
  - This mismatch is visible to auditors, owners, or clients.

#### 2.3.3 External verifiers (auditors, partners, regulators)

- **Who:** Third parties who want to verify data integrity.
- **Can do:**
  - Ask the business for CRM snapshots or exports.
  - Compare them to:
    - On-chain **critical locked records**.
  - Detect if the CRM has been quietly edited after anchoring.
- **Cannot do:**
  - Change the chain or the company’s CRM, unless they have their own credentials.
- **Value:**
  - The chain acts as a **public, immutable reference** for the records the company decided are important.

> Note: **Critical locked records are not agents**, they are objects. The agents here are the businesses, their internal staff, and external verifiers who interact with those records and with the chain.

---

## 3. Technical agents (accounts, modules, indicators)

### 3.1 Account types

- **User accounts (normal wallets):**
  - Standard Cosmos accounts with private keys.
  - Hold UAG, can send, delegate, vote, pay fees.
- **Validator operator accounts:**
  - Special accounts registered in `x/staking`.
- **Protocol-owned fund accounts:**
  - Stored in `x/fund` state, without private keys.
  - Only modules can move coins in or out, under strict rules.
- **Module accounts:**
  - `x/fund`, `x/growth`, `x/ugov`, and standard modules each may hold their own module account balances.

### 3.2 Module responsibilities (short)

- `x/fund`
  - Registers and manages **fund metadata** (type, region, permissions).
  - Exposes a **small safe API**:
    - “delegate from fund”
    - “pay salary/grant from fund”
    - “pay project grant from projects fund”
  - Rejects any action not matching:
    - Fund type (regional / Ukraine / projects).
    - Active governance decision.
    - Limits from `x/growth`.

- `x/growth`
  - Consumes **transparent indicators**:
    - Per region: taxes, GDP, exports, official regional stats.
    - National: taxes, GDP, exports, official national stats.
  - Computes and updates:
    - Max delegation per fund.
    - Max yearly / period budgets for salaries and grants.
  - Exposes these limits to `x/fund` and `x/ugov`.

- `x/ugov`
  - Builds on standard `x/gov`.
  - Defines:
    - **Presidents** (Ukraine + regions).
    - **Plans** (delegation plans, salary/grant plans, project funding plans).
    - Voting flows and communities involved (national vs. regional).
  - After a plan is approved:
    - Allows presidents to call specific `x/fund` actions **only within**:
      - The plan’s scope.
      - Fund type rules.
      - Limits from `x/growth`.

- (Future) `x/anchor` or CRM integration module
  - Receives hashes / snapshots of **critical locked records** from business CRMs.
  - Stores them on-chain in a format that:
    - Is cheap to store (hashes / Merkle roots).
    - Is easy to audit against off-chain data.

---

## 4. Agent–capability matrix (summary)

| Agent                          | Holds UAG | Can delegate | Can pay salaries / grants | Can fund projects | Can anchor CRM records | Can vote | Can change rules directly |
|--------------------------------|----------:|-------------:|--------------------------:|------------------:|-----------------------:|---------:|---------------------------|
| User / holder                  |     ✔     |      ✔       |             ✖            |         ✖         |           ✖            |    ✔    |            ✖              |
| Delegator                      |     ✔     |      ✔       |             ✖            |         ✖         |           ✖            |    ✔    |            ✖              |
| Validator                      |     ✔     |  receives    |             ✖            |         ✖         |           ✖            |    ✔    |            ✖              |
| Regional fund                  |     ✔     |      ✔       |     ✔ (limited, by data) |   maybe (TBD)     |           ✖            |    ✖    |            ✖              |
| Ukraine fund                   |     ✔     |      ✔       |     ✔ (limited, by data) |   maybe (TBD)     |           ✖            |    ✖    |            ✖              |
| Projects fund                  |     ✔     |      ✖       |             ✖            |   ✔ (by votes)    |           ✖            |    ✖    |            ✖              |
| Presidents (UA, regional)      |     ✖     | via plans    |        via plans         |    via plans      |           ✖            | maybe   |            ✖              |
| Communities (national, region) |   varies  |      ✔       |        via votes         |    via votes      |           ✖            |    ✔    |    via governance only    |
| Businesses using CRM           |   varies  |      ✖       |             ✖            |         ✖         |           ✔            |    ✖    |            ✖              |
| Internal staff (admins/devops) |     ✖     |      ✖       |             ✖            |         ✖         |   via business infra   |    ✖    |            ✖              |
| External verifiers             |     ✖     |      ✖       |             ✖            |         ✖         |  read-only comparison  |    ✖    |            ✖              |

“via plans / via votes” = the action only happens after `x/ugov` + `x/fund` checks.  
“limited, by data” = bounded by limits that come from **transparent growth indicators** in `x/growth`.

---

## 5. How to extend this file

When we add or change something that affects **who can do what**:

1. Add / update the agent definition here.
2. Link to the ADR / spec that introduced the change.
3. During review, explicitly check that:
   - Code matches this file.
   - Limits and roles are consistent with the latest pitch.

This keeps `AGENTS.md` as the **first place** a newcomer opens to understand the full economic + governance + CRM model of Ukraine Growth.

