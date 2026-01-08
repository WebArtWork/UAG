# AI.md

## Repository purpose

This is the **UAG (Ukraine Growth) blockchain repository**
(Cosmos SDK + CometBFT).

Development is performed mostly via **terminal scripts** located in the repo root.

A separate **tools repository** exists and may be used to support this repo,
but this repository is the **source of truth**.

Chain binary: **`uag`**

---

## Common environment

Minimal env used during development:

```env
RPC=http://127.0.0.1:26657
API_MNEMONIC="..."
```

Common optional env:

```env
CHAIN_ID=uag-test-1
UAG_HOME=~/.uag
PREFIX=uag
GAS_PRICE=0.025muag
FEES=5000muag
```

Local ports:

- RPC: `26657`
- REST: `1317`
- RPC proxy (if used): `26658`

---

## Common terminal scripts

### Chain lifecycle

```bash
./build.sh
```

Build `uag`.

```bash
./init-chain.sh
```

Initialize genesis + validator (no start).

```bash
./serve.sh
```

Start local node.
Auto-initializes chain if home dir is missing.

```bash
./wipe.sh
```

Wipe chain state and regenerate protos.

```bash
./clear.sh
```

Clear caches, binaries, node modules, docker leftovers.

---

### Wallet & authority

```bash
./wallet-new.sh
```

Generate mnemonic + address.

```bash
./wallet-check.sh
```

Check balances and account existence.

```bash
./api-check.sh
```

Verify that `API_MNEMONIC` matches on-chain citizen authority.

```bash
./set-api-wallet.sh
```

Patch genesis to set citizen API authority.

---

### Citizen module

```bash
./citizen-set-region.sh <region_id>
```

Broadcast `MsgSetCitizenRegion` using API authority.

---

### Go / Cosmos diagnostics

```bash
./mod.sh
```

Check Go module health (Cosmos SDK, CometBFT mismatches).

```bash
./test.sh
```

Run `go mod tidy` and `go test`.

---

### Protobuf

```bash
./proto.sh
```

Generate protobufs and format code.

---

### Browser / WSL helper (optional)

```bash
./rpc-proxy.sh
```

Expose RPC with CORS (port `26658`).

---

## Common failure signals

- **account not found** → wallet not funded or chain not running
- **authority mismatch** → wrong `API_MNEMONIC`
- **unknown msg / module** → protos not generated or node not restarted
- **Go build errors** → run `./mod.sh`
- **RPC unreachable** → node not started or wrong RPC
