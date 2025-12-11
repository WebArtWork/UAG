package types

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QueryClient is the client API for Query service.
type QueryClient interface {
	Fund(ctx context.Context, in *QueryFundRequest, opts ...grpc.CallOption) (*QueryFundResponse, error)
	Funds(ctx context.Context, in *QueryFundsRequest, opts ...grpc.CallOption) (*QueryFundsResponse, error)
	FundsByType(ctx context.Context, in *QueryFundsByTypeRequest, opts ...grpc.CallOption) (*QueryFundsByTypeResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Fund(ctx context.Context, in *QueryFundRequest, opts ...grpc.CallOption) (*QueryFundResponse, error) {
	out := new(QueryFundResponse)
	err := c.cc.Invoke(ctx, "/uagd.fund.v1.Query/Fund", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Funds(ctx context.Context, in *QueryFundsRequest, opts ...grpc.CallOption) (*QueryFundsResponse, error) {
	out := new(QueryFundsResponse)
	err := c.cc.Invoke(ctx, "/uagd.fund.v1.Query/Funds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) FundsByType(ctx context.Context, in *QueryFundsByTypeRequest, opts ...grpc.CallOption) (*QueryFundsByTypeResponse, error) {
	out := new(QueryFundsByTypeResponse)
	err := c.cc.Invoke(ctx, "/uagd.fund.v1.Query/FundsByType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	Fund(context.Context, *QueryFundRequest) (*QueryFundResponse, error)
	Funds(context.Context, *QueryFundsRequest) (*QueryFundsResponse, error)
	FundsByType(context.Context, *QueryFundsByTypeRequest) (*QueryFundsByTypeResponse, error)
	mustEmbedUnimplementedQueryServer()
}

type UnimplementedQueryServer struct{}

func (UnimplementedQueryServer) Fund(context.Context, *QueryFundRequest) (*QueryFundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fund not implemented")
}
func (UnimplementedQueryServer) Funds(context.Context, *QueryFundsRequest) (*QueryFundsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Funds not implemented")
}
func (UnimplementedQueryServer) FundsByType(context.Context, *QueryFundsByTypeRequest) (*QueryFundsByTypeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FundsByType not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Fund_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryFundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Fund(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/uagd.fund.v1.Query/Fund"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Fund(ctx, req.(*QueryFundRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Funds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryFundsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Funds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/uagd.fund.v1.Query/Funds"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Funds(ctx, req.(*QueryFundsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_FundsByType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryFundsByTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).FundsByType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/uagd.fund.v1.Query/FundsByType"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).FundsByType(ctx, req.(*QueryFundsByTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "uagd.fund.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Fund", Handler: _Query_Fund_Handler},
		{MethodName: "Funds", Handler: _Query_Funds_Handler},
		{MethodName: "FundsByType", Handler: _Query_FundsByType_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "uagd/fund/v1/query.proto",
}
