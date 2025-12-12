package types

import "google.golang.org/grpc"

// RegisterMsgServer wires the ugov message server into the gRPC router.
// TODO: replace the placeholder implementation once protobuf services are added.
func RegisterMsgServer(_ grpc.ServiceRegistrar, _ interface{}) {}

// RegisterQueryServer wires the ugov query server into the gRPC router.
// TODO: replace the placeholder implementation once protobuf services are added.
func RegisterQueryServer(_ grpc.ServiceRegistrar, _ interface{}) {}
