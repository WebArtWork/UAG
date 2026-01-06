package app

import (
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterIBC is a stub for CLI wiring until IBC modules support depinject.
func RegisterIBC(_ codec.Codec) map[string]appmodule.AppModule {
	return map[string]appmodule.AppModule{}
}
