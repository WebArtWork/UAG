package types

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QueryClient is the client API for Query service.
type QueryClient interface {
	RegionByAddress(ctx context.Context, in *QueryRegionByAddressRequest, opts ...grpc.CallOption) (*QueryRegionByAddressResponse, error)
	AddressesByRegion(ctx context.Context, in *QueryAddressesByRegionRequest, opts ...grpc.CallOption) (*QueryAddressesByRegionResponse, error)
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) RegionByAddress(ctx context.Context, in *QueryRegionByAddressRequest, opts ...grpc.CallOption) (*QueryRegionByAddressResponse, error) {
	out := new(QueryRegionByAddressResponse)
	err := c.cc.Invoke(ctx, "/uagd.citizen.v1.Query/RegionByAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AddressesByRegion(ctx context.Context, in *QueryAddressesByRegionRequest, opts ...grpc.CallOption) (*QueryAddressesByRegionResponse, error) {
	out := new(QueryAddressesByRegionResponse)
	err := c.cc.Invoke(ctx, "/uagd.citizen.v1.Query/AddressesByRegion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/uagd.citizen.v1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	RegionByAddress(context.Context, *QueryRegionByAddressRequest) (*QueryRegionByAddressResponse, error)
	AddressesByRegion(context.Context, *QueryAddressesByRegionRequest) (*QueryAddressesByRegionResponse, error)
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	mustEmbedUnimplementedQueryServer()
}

type UnimplementedQueryServer struct{}

func (UnimplementedQueryServer) RegionByAddress(context.Context, *QueryRegionByAddressRequest) (*QueryRegionByAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegionByAddress not implemented")
}
func (UnimplementedQueryServer) AddressesByRegion(context.Context, *QueryAddressesByRegionRequest) (*QueryAddressesByRegionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddressesByRegion not implemented")
}
func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_RegionByAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRegionByAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).RegionByAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/uagd.citizen.v1.Query/RegionByAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).RegionByAddress(ctx, req.(*QueryRegionByAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AddressesByRegion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAddressesByRegionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AddressesByRegion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/uagd.citizen.v1.Query/AddressesByRegion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AddressesByRegion(ctx, req.(*QueryAddressesByRegionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/uagd.citizen.v1.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "uagd.citizen.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegionByAddress",
			Handler:    _Query_RegionByAddress_Handler,
		},
		{
			MethodName: "AddressesByRegion",
			Handler:    _Query_AddressesByRegion_Handler,
		},
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "uagd/citizen/v1/query.proto",
}
