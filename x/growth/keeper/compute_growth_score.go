package keeper

import (
	"uag/x/growth/types"
)

// Matches the call site in msg_server.go: m.ComputeGrowthScore(metric)
func (m msgServer) ComputeGrowthScore(_ types.RegionMetric) error {
	// no-op for now
	return nil
}
