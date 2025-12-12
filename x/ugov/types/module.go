package types

import gogoproto "github.com/cosmos/gogoproto/proto"

// Module defines the module config type needed for app wiring.
type Module struct{}

var _ gogoproto.Message = (*Module)(nil)

func (m *Module) Reset()         {}
func (m *Module) String() string { return ModuleName }
func (m *Module) ProtoMessage()  {}
