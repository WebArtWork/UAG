package module

import (
	"cosmossdk.io/depinject/appconfig"

	modulev1 "github.com/CosmWasm/wasmd/api/wasm/module/v1"
)

// init registers the wasm module configuration with the app wiring system.
// This minimal stub only needs to expose the module config type so the
// application can compose its dependency graph.
func init() {
	appconfig.RegisterModule(&modulev1.Module{})
}
