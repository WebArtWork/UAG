# AGENTS.md – How humans and machines work on UAGD

This document defines **roles, responsibilities, and rules** for everyone (and everything) interacting with the **UAGD – Ukraine Growth blockchain** codebase.

It applies to:

- human contributors,
- maintainers,
- automated agents (AI, bots, CI),
- external auditors and integrators.

UAGD is a **public protocol**, not a startup codebase.
Rules matter.

---

## 1. Core principle

**Code is law. Agents are servants.**

No agent − human or automated − has the right to:

- bypass protocol rules,
- introduce hidden control paths,
- weaken immutability guarantees,
- add discretionary power where growth-based logic is required.

If behavior is not enforced by code and governance − **it must not exist**.

---

## 2. Agent types

### 2.1 Human contributors

Includes:

- core protocol developers,
- module authors,
- reviewers,
- documentation writers.

Allowed:

- propose changes via PRs,
- discuss designs publicly,
- implement features **exactly as specified**.

Not allowed:

- merge own PRs without review,
- add backdoors, admin shortcuts, or emergency keys,
- introduce logic that allows silent fund movement.

---

### 2.2 Maintainers

Maintainers are **operators**, not owners.

They can:

- review and merge PRs,
- manage releases,
- coordinate upgrades.

They cannot:

- change economic rules unilaterally,
- override governance outcomes,
- deploy code that contradicts README philosophy.

All major changes must be:

- documented,
- auditable,
- traceable to governance decisions.

---

### 2.3 Automated agents (CI, bots, AI)

Includes:

- GitHub Actions,
- code generators,
- AI coding assistants,
- test automation.

Rules:

- automated agents **never merge directly to main**
- generated code must be reviewable and deterministic
- AI output is treated as **untrusted until reviewed**

AI agents must not:

- invent protocol rules,
- simplify economic logic for convenience,
- replace explicit constraints with comments or assumptions.

---

## 3. Protocol boundaries (non-negotiable)

Agents **must not violate** the following invariants:

### Funds

- All UAG funds are protocol-owned
- No private treasury keys
- No discretionary spending paths
- All limits derive from `x/growth`

### Identity

- No PII on-chain
- `x/citizen` stores only:

  - address
  - region id
  - verification flags
  - KYC source tag
  - proof hash

- Only registrar accounts can write

### Governance

- Votes decide
- Presidents execute
- No hidden parameters
- No emergency overrides

If a change weakens any of the above − it must be rejected.

---

## 4. Contribution rules

### 4.1 Pull requests

Every PR must:

- have a clear scope
- reference affected modules
- explain **why** the change aligns with UAG philosophy
- include tests for economic logic

PRs that:

- add shortcuts,
- centralize power,
- blur responsibility boundaries

will be closed.

---

### 4.2 Tests are mandatory

For core modules (`x/fund`, `x/growth`, `x/citizen`, `x/gov`):

- happy paths
- failure paths
- limit enforcement
- permission checks

Untested economic logic is considered **unsafe**.

---

## 5. Versioning and upgrades

Protocol upgrades:

- must be explicit,
- must be reviewable,
- must not silently change economic meaning.

Breaking changes require:

- migration logic,
- clear documentation,
- governance approval (where applicable).

---

## 6. CRM & business integrations

Agents working on:

- CRM anchoring,
- CosmWasm contracts,
- external tooling

must respect:

- the company chooses what to anchor
- anchored records are immutable
- off-chain systems must **fail loudly** if chain state mismatches

No agent may introduce:

- silent resync logic,
- “admin fix” buttons,
- overwrite paths for anchored data.

---

## 7. Security posture

Assume:

- contributors can be compromised,
- infrastructure can be attacked,
- insiders can act maliciously.

Design accordingly:

- minimal permissions,
- explicit checks,
- no trust by role.

If an agent needs power − it must be justified by code and votes.

---

## 8. Final rule

If you are an agent working on UAGD:

> **You do not decide outcomes.
> You implement rules that make outcomes unavoidable.**

Anything else is out of scope for this project.
