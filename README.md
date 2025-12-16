# UAGD – Ukraine Growth blockchain

`uagd` is the reference implementation of the **Ukraine Growth (UAG)** blockchain, built on **Cosmos SDK + CometBFT** and scaffolded with **Ignite CLI**.

Ukraine Growth is not a typical crypto project.
It is a **public economic coordination layer** designed to tie on-chain rules to **real-world regional growth**, using transparent governance and immutable records.

---

## Network basics

- **Chain ID:** `uag-1`
- **Native token:** `UAG`
- **Base denom:** `uuag`

  - `1 UAG = 1,000,000 uuag`

- **Initial supply:** `370,000,000 UAG`
- **Inflation:** `Without`

The node binary is **`uagdd`**, used by full nodes and validators.

---

## Core idea

UAG is a coordination token whose **usage limits are algorithmically tied to real economic growth** − not political promises or private control. Funds cannot be secretly moved, rewritten, or inflated beyond what Ukraine and its regions demonstrably produce in the real world.

---

## Token distribution

All UAG supply is fixed at genesis (no inflation; validators earn fees/gas) and split into **three protocol-owned buckets**:

### 1. Regional funds – **70%**

- One fund for **each of the 27 regions**
- Coins are owned by the **protocol**, not by people

Each regional fund can **only**:

1. **Delegate** coins to validators
2. **Pay limited salaries and grants** to people working on Ukraine Growth in that region
3. If regional **occupation exceeds 50%**, all delegations, salaries and grants from that fund are locked

The **maximum allowed amounts** depend on how well the region grows over time.

---

### 2. Ukraine-level fund – **20%**

A single national fund that can **only**:

1. Delegate coins to validators
2. Pay limited salaries and grants for Ukraine-level work

Its limits depend on **country-wide growth**, not on political decisions.

---

### 3. Projects fund – **10%**

1. Delegate coins to validators
2. Pay limited salaries and grants for developers

- Used only for **ecosystem projects and tools**
- **Every spend requires a global on-chain vote**
- No unilateral control, no emergency keys
- With maximum allowed amounts

---

## Growth-based limits

Neither regions nor the Ukraine-level fund can freely spend.

Their limits are dynamically derived from **transparent real-world indicators**, such as:

- Taxes collected
- GDP
- Exports

As the country or a region grows, **its on-chain capacity grows**.
If growth stalls − limits tighten automatically.

This makes UAG **anti-populist by design**.

---

## Verified citizens and regions (`x/citizen`)

To connect real people and regions to on-chain rewards, UAG introduces a **minimal on-chain identity layer**:

- A **citizen** is defined as:
  - a wallet address
  - a region identifier
  - a `verified` flag
  - a generic KYC source (e.g. `DIIA_V1`)
  - a hash pointing to an off-chain KYC record

Personal data (full name, passport, etc.) stays **off-chain** in partner systems (e.g. Diia-integrated backends).
On-chain we only store:

- `address → { region_id, verified, kyc_source, proof_hash }`

Only special **registrar accounts** (run by the UAG backend and partners) can create or update citizen records.
This allows:

- rewards to be directed **only to verified citizens**
- rewards and programs to be **filtered by region**
- growth-based logic to reason about “real people, real regions” without leaking PII on-chain.

---

## Governance model (»presidents are servants«)

At both levels − **regional** and **national** − governance works the same way:

1. A president **prepares a plan**
   - how much to delegate
   - whom to pay
   - which grants to support
2. The plan is **publicly explained**
3. The community **votes on-chain**
4. If approved − the president **executes exactly that plan**

Presidents:

- ❌ cannot secretly move coins
- ❌ cannot invent extra payouts
- ❌ cannot bypass votes
- ✅ only execute approved decisions

They are **operators, not owners**.

---

## CRM + immutable business records

Ukraine Growth also provides **ready-made CRM systems** for businesses.

Each CRM can enable **UAG mode**.

### How it works

- A company selects certain records as **critical**:
  - ownership
  - balances
  - contracts
  - inventory checkpoints
  - audit milestones
- These records are **anchored to the UAG chain**:
  - either via **native modules** or via **CosmWasm smart contracts (`x/wasm`)**
- Once written:

  - ❌ cannot be changed
  - ❌ cannot be deleted
  - ❌ cannot be rewritten by devops, admins, employees, or owners

### Why this matters

This directly blocks the classic **»devops + employee« fraud**:

- In a normal database → numbers can be quietly edited
- With UAG → any later database change **no longer matches the chain**

The company decides **what to anchor**.
The chain guarantees **those facts are immutable**.

---

## Smart contracts and builders layer (`x/wasm`)

UAG is not only a protocol for Ukraine Growth itself; it is also a **platform for builders**.

By integrating **CosmWasm** as `x/wasm`, the chain allows anyone with UAG to:

- upload smart contracts (WASM code)
- instantiate them
- execute methods that:
  - manage their own state
  - store arbitrary business data
  - implement custom logic, DAOs, games, audit trails, etc.

Every contract interaction is a normal transaction:

- pays **gas in `uuag`**
- writes to the chain’s immutable state
- can interoperate with UAG-native modules (e.g. read `x/citizen` or react to `x/fund` decisions)

This gives:

- **Businesses** – a way to build tailor-made logic on top of UAG
- **Developers** – a permissionless environment to deploy tools for citizens, regions, and companies

---

## Protocol modules (what can do what)

### `x/citizen`

Minimal identity layer:

- Maps **wallet → region + verified flag**
- Stores only:
  - address
  - region id
  - verified / active flags
  - KYC source tag
  - hash of off-chain KYC payload
- Write access restricted to **registrar accounts**
- Read access is public for other modules and apps

Used by:

- `x/fund` and `x/growth` to direct rewards and programs to **verified citizens** in specific regions.

---

### `x/fund`

Controls all protocol-owned funds.

- 27 Regional funds
- Ukraine-level fund
- Engineering fund

Enforces:

- allowed actions (delegate, limited payroll/grants)
- spending caps
- real life occupation locks
- vote-required execution

No module or account can bypass this logic.

---

### `x/growth`

Connects **off-chain economic indicators** to **on-chain limits**.

- Stores growth metrics per region and country
- Calculates dynamic caps
- Feeds limits into `x/fund`
- Makes spending mathematically dependent on real growth

---

### `x/gov`

Global governance.

- Community votes
- Parameter changes
- Projects fund approvals
- Protocol upgrades

---

### `x/wasm` (CosmWasm smart contracts)

Open smart-contract platform:

- Any address with UAG can:
  - upload contracts
  - instantiate contracts
  - execute contract messages
- Contracts can:
  - store arbitrary state
  - serve as app backends (CRMs, registries, DAOs)
  - integrate with UAG’s native economics

Gas and fees are always paid in **`uuag`**.

#### CosmWasm CLI quickstart

All commands use the `uagdd` binary and the `uuag` fee token.

1. **Store a contract** (any account can upload):

   ```bash
   uagdd tx wasm store path/to/contract.wasm \\
     --from <key-name> \\
     --chain-id uag-test-1 \\
     --gas auto --gas-adjustment 1.3 --fees 7500uuag
   ```

2. **Instantiate** the uploaded code (replace `<code-id>` from the store result):

   ```bash
   uagdd tx wasm instantiate <code-id> '{"count":0}' \\
     --label "demo-counter" \\
     --admin <admin-address> \\
     --from <key-name> --chain-id uag-test-1 \\
     --gas auto --gas-adjustment 1.3 --fees 6000uuag
   ```

3. **Execute** a contract message:

   ```bash
   uagdd tx wasm execute <contract-address> '{"increment":{}}' \\
     --from <key-name> --chain-id uag-test-1 \\
     --gas auto --gas-adjustment 1.3 --fees 5000uuag
   ```

4. **Query** contract state:

   ```bash
   uagdd query wasm contract-state smart <contract-address> '{"get_count":{}}'
   ```

Uploads are permissionless, and state writes pay gas in `uuag`.

---

### `x/staking` (Cosmos SDK)

Standard validator and delegation logic.

- Validators
- Delegators
- Slashing
- Rewards

Used by `x/fund`, but not controlled by it.

---

## What this chain explicitly forbids

- ❌ Private treasury keys
- ❌ Hidden admin balances
- ❌ Manual minting
- ❌ Silent database edits
- ❌ Political overrides of protocol rules

If it’s not voted and enforced by code − **it doesn’t happen**.

---

## Philosophy (short)

Ukraine Growth is not about speculation.
It is about **making growth measurable, enforceable, and irreversible**.

- ✅ Code > promises
- ✅ Votes > power
- ✅ Growth > speculation
