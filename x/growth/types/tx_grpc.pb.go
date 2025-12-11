package types

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MsgClient is the client API for Msg service.
type MsgClient interface {
	SetRegionMetric(ctx context.Context, in *MsgSetRegionMetric, opts ...grpc.CallOption) (*MsgSetRegionMetricResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) SetRegionMetric(ctx context.Context, in *MsgSetRegionMetric, opts ...grpc.CallOption) (*MsgSetRegionMetricResponse, error) {
	out := new(MsgSetRegionMetricResponse)
	err := c.cc.Invoke(ctx, "/uagd.growth.v1.Msg/SetRegionMetric", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	SetRegionMetric(context.Context, *MsgSetRegionMetric) (*MsgSetRegionMetricResponse, error)
	mustEmbedUnimplementedMsgServer()
}

type UnimplementedMsgServer struct{}

func (UnimplementedMsgServer) SetRegionMetric(context.Context, *MsgSetRegionMetric) (*MsgSetRegionMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetRegionMetric not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_SetRegionMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetRegionMetric)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetRegionMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/uagd.growth.v1.Msg/SetRegionMetric",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetRegionMetric(ctx, req.(*MsgSetRegionMetric))
	}
	return interceptor(ctx, in, info, handler)
}

var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "uagd.growth.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetRegionMetric",
			Handler:    _Msg_SetRegionMetric_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "uagd/growth/v1/tx.proto",
}
