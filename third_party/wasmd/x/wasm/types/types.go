package types

import "fmt"

// ModuleName defines the wasm module name.
const ModuleName = "wasm"

// AccessType mirrors the permission levels for wasm code upload.
type AccessType uint32

const (
	// AccessTypeNobody blocks all uploads.
	AccessTypeNobody AccessType = iota
	// AccessTypeOnlyAddress limits uploads to a single address.
	AccessTypeOnlyAddress
	// AccessTypeAnyOfAddresses allows a predefined set of addresses to upload.
	AccessTypeAnyOfAddresses
	// AccessTypeEvery body allows anyone to upload code.
	AccessTypeEverybody
)

// AccessConfig configures who may upload wasm code.
type AccessConfig struct {
	Permission AccessType `protobuf:"varint,1,opt,name=permission,proto3,enum=cosmwasm.wasm.v1.AccessType" json:"permission,omitempty" yaml:"permission"`
}

// Reset implements the proto.Message interface.
func (m *AccessConfig) Reset() { *m = AccessConfig{} }

// ProtoMessage implements the proto.Message interface.
func (*AccessConfig) ProtoMessage() {}

// String implements the proto.Message interface.
func (m *AccessConfig) String() string {
	return fmt.Sprintf("AccessConfig{Permission:%v}", m.Permission)
}

// Module describes the wasm module wiring configuration.
type Module struct {
	UploadPermission AccessConfig `json:"upload_permission" yaml:"upload_permission"`
}
