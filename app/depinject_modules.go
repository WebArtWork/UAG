package app

// Blank-import module wiring packages so their init() registers depinject providers.
// SDK v0.53: most SDK modules register from the module root (x/<name>), not x/<name>/module.
// Some split modules live under cosmossdk.io/x/*.

import (
	// Cosmos SDK modules (within github.com/cosmos/cosmos-sdk)
	_ "github.com/cosmos/cosmos-sdk/x/auth"
	_ "github.com/cosmos/cosmos-sdk/x/authz"
	_ "github.com/cosmos/cosmos-sdk/x/bank"
	_ "github.com/cosmos/cosmos-sdk/x/consensus"
	_ "github.com/cosmos/cosmos-sdk/x/crisis"
	_ "github.com/cosmos/cosmos-sdk/x/distribution"
	_ "github.com/cosmos/cosmos-sdk/x/genutil"
	_ "github.com/cosmos/cosmos-sdk/x/gov"
	_ "github.com/cosmos/cosmos-sdk/x/group"
	_ "github.com/cosmos/cosmos-sdk/x/params"
	_ "github.com/cosmos/cosmos-sdk/x/slashing"
	_ "github.com/cosmos/cosmos-sdk/x/staking"

	// Split-out modules (cosmossdk.io/x/*)
	_ "cosmossdk.io/x/circuit"
	_ "cosmossdk.io/x/evidence"
	_ "cosmossdk.io/x/feegrant"
	_ "cosmossdk.io/x/nft/module"
	_ "cosmossdk.io/x/upgrade"

	// UAG modules
	_ "uag/x/citizen/module"
	_ "uag/x/growth/module"
)
