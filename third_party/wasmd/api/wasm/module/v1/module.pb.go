package modulev1

import (
	"fmt"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	_ "google.golang.org/protobuf/reflect/protoreflect"
)

// Module is a minimal stand-in for the wasm module configuration proto.
type Module struct {
	UploadPermission wasmtypes.AccessConfig `protobuf:"bytes,1,opt,name=upload_permission,json=uploadPermission,proto3" json:"upload_permission,omitempty"`
}

// Reset implements the proto.Message interface.
func (m *Module) Reset() { *m = Module{} }

// ProtoMessage implements the proto.Message interface.
func (*Module) ProtoMessage() {}

// String implements the proto.Message interface.
func (m *Module) String() string {
	return fmt.Sprintf("Module{UploadPermission:%v}", m.UploadPermission)
}
