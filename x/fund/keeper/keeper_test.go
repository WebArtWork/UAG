package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"uagd/x/fund/types"
)

func TestGenesisState_Validate_Default(t *testing.T) {
	gs := types.DefaultGenesis()
	require.NoError(t, gs.Validate())
}
