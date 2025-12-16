package types

func DefaultGenesis() GenesisState {
	return GenesisState{
		Params:     DefaultParams(),
		Plans:      []Plan{},
		NextPlanId: 1,
	}
}
