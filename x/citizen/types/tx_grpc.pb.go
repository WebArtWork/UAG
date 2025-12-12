package types

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MsgClient is the client API for Msg service.
type MsgClient interface {
	SetCitizenRegion(ctx context.Context, in *MsgSetCitizenRegion, opts ...grpc.CallOption) (*MsgSetCitizenRegionResponse, error)
	ClearCitizenRegion(ctx context.Context, in *MsgClearCitizenRegion, opts ...grpc.CallOption) (*MsgClearCitizenRegionResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) SetCitizenRegion(ctx context.Context, in *MsgSetCitizenRegion, opts ...grpc.CallOption) (*MsgSetCitizenRegionResponse, error) {
	out := new(MsgSetCitizenRegionResponse)
	err := c.cc.Invoke(ctx, "/uagd.citizen.v1.Msg/SetCitizenRegion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ClearCitizenRegion(ctx context.Context, in *MsgClearCitizenRegion, opts ...grpc.CallOption) (*MsgClearCitizenRegionResponse, error) {
	out := new(MsgClearCitizenRegionResponse)
	err := c.cc.Invoke(ctx, "/uagd.citizen.v1.Msg/ClearCitizenRegion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	SetCitizenRegion(context.Context, *MsgSetCitizenRegion) (*MsgSetCitizenRegionResponse, error)
	ClearCitizenRegion(context.Context, *MsgClearCitizenRegion) (*MsgClearCitizenRegionResponse, error)
	mustEmbedUnimplementedMsgServer()
}

type UnimplementedMsgServer struct{}

func (UnimplementedMsgServer) SetCitizenRegion(context.Context, *MsgSetCitizenRegion) (*MsgSetCitizenRegionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetCitizenRegion not implemented")
}
func (UnimplementedMsgServer) ClearCitizenRegion(context.Context, *MsgClearCitizenRegion) (*MsgClearCitizenRegionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearCitizenRegion not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_SetCitizenRegion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetCitizenRegion)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetCitizenRegion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/uagd.citizen.v1.Msg/SetCitizenRegion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetCitizenRegion(ctx, req.(*MsgSetCitizenRegion))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClearCitizenRegion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClearCitizenRegion)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClearCitizenRegion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/uagd.citizen.v1.Msg/ClearCitizenRegion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClearCitizenRegion(ctx, req.(*MsgClearCitizenRegion))
	}
	return interceptor(ctx, in, info, handler)
}

var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "uagd.citizen.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetCitizenRegion",
			Handler:    _Msg_SetCitizenRegion_Handler,
		},
		{
			MethodName: "ClearCitizenRegion",
			Handler:    _Msg_ClearCitizenRegion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "uagd/citizen/v1/tx.proto",
}
