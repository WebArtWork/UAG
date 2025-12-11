package types

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MsgClient is the client API for Msg service.
type MsgClient interface {
	ExecuteFundPlan(ctx context.Context, in *MsgExecuteFundPlan, opts ...grpc.CallOption) (*MsgExecuteFundPlanResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) ExecuteFundPlan(ctx context.Context, in *MsgExecuteFundPlan, opts ...grpc.CallOption) (*MsgExecuteFundPlanResponse, error) {
	out := new(MsgExecuteFundPlanResponse)
	err := c.cc.Invoke(ctx, "/uagd.fund.v1.Msg/ExecuteFundPlan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	ExecuteFundPlan(context.Context, *MsgExecuteFundPlan) (*MsgExecuteFundPlanResponse, error)
	mustEmbedUnimplementedMsgServer()
}

type UnimplementedMsgServer struct{}

func (UnimplementedMsgServer) ExecuteFundPlan(context.Context, *MsgExecuteFundPlan) (*MsgExecuteFundPlanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteFundPlan not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_ExecuteFundPlan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgExecuteFundPlan)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ExecuteFundPlan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/uagd.fund.v1.Msg/ExecuteFundPlan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ExecuteFundPlan(ctx, req.(*MsgExecuteFundPlan))
	}
	return interceptor(ctx, in, info, handler)
}

var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "uagd.fund.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExecuteFundPlan",
			Handler:    _Msg_ExecuteFundPlan_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "uagd/fund/v1/tx.proto",
}
