package modulev1

import (
	_ "google.golang.org/protobuf/reflect/protoreflect"
)

// Module is a minimal stand-in for the wasm module configuration proto.
type Module struct {
	UploadPermission interface{} `protobuf:"bytes,1,opt,name=upload_permission,json=uploadPermission,proto3" json:"upload_permission,omitempty"`
}
