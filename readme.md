
# UAGD – Ukraine Growth blockchain

`uagd` is the reference implementation of the **Ukraine Growth (UAG)** blockchain, built on **Cosmos SDK + CometBFT** and scaffolded with [Ignite CLI](https://ignite.com/cli).

- **Chain ID:** `uag-1`  
- **Base denom:** `uuag` (1 UAG = 1,000,000 uuag)

This repo contains the source code for the node binary **`uagdd`**, which is used to run full nodes and validators.

---

## Releases & versioning

Releases follow semantic versioning:

- `v0.x.y` – testnets / early public networks  
- `v1.0.0` – first mainnet release (planned)

Each tag `v*` triggers CI to build binaries for:

- Linux (amd64)
- macOS (amd64, arm64)
- Windows (amd64)

You can find them in **GitHub → Releases**.

---

## Install `uagdd` (validator / full node)

> Always pick the archive that matches your OS and CPU from the latest release.

### 1. Linux (amd64)

1. Download the latest Linux archive from **Releases**.
2. Extract it and move the binary:

```bash
tar -xzf <archive-name>.tar.gz
cd <extracted-folder>
chmod +x uagdd
sudo mv uagdd /usr/local/bin/uagdd
uagdd version
````

You should see the version printed without errors.

---

### 2. macOS (Intel & Apple Silicon)

1. Download the `darwin-amd64` (Intel) or `darwin-arm64` (Apple Silicon) archive.
2. Extract and install:

```bash
tar -xzf <archive-name>.tar.gz
cd <extracted-folder>
chmod +x uagdd
sudo mv uagdd /usr/local/bin/uagdd
uagdd version
```

On macOS you might need to allow the binary in **System Settings → Privacy & Security** if Gatekeeper complains.

---

### 3. Windows (amd64)

1. Download the Windows archive (`windows-amd64`).
2. Extract it (e.g. with 7-Zip).
3. Move `uagdd.exe` somewhere in your `PATH` (or run from the extracted folder):

```powershell
uagdd.exe version
```

---

## Running a node

You can run `uagdd` either as a **full node** or a **validator**.
The basic flow is the same for all OSes.

### 1. Initialize node home

```bash
uagdd init <moniker> --chain-id uag-1
```

This creates the config and data directory, usually:

* Linux/macOS: `~/.uagd`
* Windows: `%USERPROFILE%\.uagd`

### 2. Put the correct `genesis.json`

Download the official `genesis.json` for the current UAG network (testnet or mainnet) and place it into:

```bash
# Linux/macOS
cp genesis.json ~/.uagd/config/genesis.json
```

On Windows, place it into `%USERPROFILE%\.uagd\config\genesis.json`.

### 3. Configure networking & gas prices

Edit:

* `~/.uagd/config/config.toml`

  * set `persistent_peers` / `seeds`
* `~/.uagd/config/app.toml`

  * set `minimum-gas-prices` (e.g. `0.025uuag`)

Save and exit.

### 4. Start as a full node

```bash
uagdd start
```

If everything is correct, the node will start syncing blocks.

---

## Becoming a validator (basics)

1. **Create a key (validator operator):**

```bash
uagdd keys add validator
```

2. **Fund the address** with `uuag` (from faucet or another wallet).

3. **Create validator transaction** (example):

```bash
uagdd tx staking create-validator \
  --from validator \
  --amount 1000000uuag \
  --pubkey "$(uagdd tendermint show-validator)" \
  --moniker "<your-moniker>" \
  --chain-id uag-1 \
  --commission-rate "0.05" \
  --commission-max-rate "0.20" \
  --commission-max-change-rate "0.01" \
  --min-self-delegation "1"
```

4. After the tx is included and your node is healthy, you will enter the active validator set when you have enough voting power.

> **Note:** For production we will publish a separate “Validator Guide” with recommended hardware, security practices, and upgrade process.

---

## Development

For local development (not for validators), you can still use Ignite’s dev workflow:

```bash
ignite chain serve
```

This will build, init, and run a local single-node devnet.

---

## Learn more

* [Ignite CLI](https://ignite.com/cli)
* [Cosmos SDK docs](https://docs.cosmos.network)
* [Ignite CLI docs](https://docs.ignite.com)
