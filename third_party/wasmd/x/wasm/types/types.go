package types

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
	Permission AccessType `json:"permission" yaml:"permission"`
}

// Module describes the wasm module wiring configuration.
type Module struct {
	UploadPermission AccessConfig `json:"upload_permission" yaml:"upload_permission"`
}
