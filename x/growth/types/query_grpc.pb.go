package types

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QueryClient is the client API for Query service.
type QueryClient interface {
	RegionMetric(ctx context.Context, in *QueryRegionMetricRequest, opts ...grpc.CallOption) (*QueryRegionMetricResponse, error)
	GrowthScore(ctx context.Context, in *QueryGrowthScoreRequest, opts ...grpc.CallOption) (*QueryGrowthScoreResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) RegionMetric(ctx context.Context, in *QueryRegionMetricRequest, opts ...grpc.CallOption) (*QueryRegionMetricResponse, error) {
	out := new(QueryRegionMetricResponse)
	err := c.cc.Invoke(ctx, "/uagd.growth.v1.Query/RegionMetric", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GrowthScore(ctx context.Context, in *QueryGrowthScoreRequest, opts ...grpc.CallOption) (*QueryGrowthScoreResponse, error) {
	out := new(QueryGrowthScoreResponse)
	err := c.cc.Invoke(ctx, "/uagd.growth.v1.Query/GrowthScore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	RegionMetric(context.Context, *QueryRegionMetricRequest) (*QueryRegionMetricResponse, error)
	GrowthScore(context.Context, *QueryGrowthScoreRequest) (*QueryGrowthScoreResponse, error)
	mustEmbedUnimplementedQueryServer()
}

type UnimplementedQueryServer struct{}

func (UnimplementedQueryServer) RegionMetric(context.Context, *QueryRegionMetricRequest) (*QueryRegionMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegionMetric not implemented")
}
func (UnimplementedQueryServer) GrowthScore(context.Context, *QueryGrowthScoreRequest) (*QueryGrowthScoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GrowthScore not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_RegionMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRegionMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).RegionMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/uagd.growth.v1.Query/RegionMetric",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).RegionMetric(ctx, req.(*QueryRegionMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GrowthScore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGrowthScoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GrowthScore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/uagd.growth.v1.Query/GrowthScore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GrowthScore(ctx, req.(*QueryGrowthScoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "uagd.growth.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegionMetric",
			Handler:    _Query_RegionMetric_Handler,
		},
		{
			MethodName: "GrowthScore",
			Handler:    _Query_GrowthScore_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "uagd/growth/v1/query.proto",
}
