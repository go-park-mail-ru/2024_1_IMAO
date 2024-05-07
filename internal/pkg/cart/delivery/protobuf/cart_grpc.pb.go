// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: cart.proto

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CartClient is the client API for Cart service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CartClient interface {
	GetCartByUserID(ctx context.Context, in *UserIdRequest, opts ...grpc.CallOption) (*ReturningAdvertList, error)
	DeleteAdvByIDs(ctx context.Context, in *UserIdAdvertIdRequest, opts ...grpc.CallOption) (*DeleteAdvResponse, error)
	AppendAdvByIDs(ctx context.Context, in *UserIdAdvertIdRequest, opts ...grpc.CallOption) (*AppendAdvResponse, error)
}

type cartClient struct {
	cc grpc.ClientConnInterface
}

func NewCartClient(cc grpc.ClientConnInterface) CartClient {
	return &cartClient{cc}
}

func (c *cartClient) GetCartByUserID(ctx context.Context, in *UserIdRequest, opts ...grpc.CallOption) (*ReturningAdvertList, error) {
	out := new(ReturningAdvertList)
	err := c.cc.Invoke(ctx, "/Cart/GetCartByUserID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cartClient) DeleteAdvByIDs(ctx context.Context, in *UserIdAdvertIdRequest, opts ...grpc.CallOption) (*DeleteAdvResponse, error) {
	out := new(DeleteAdvResponse)
	err := c.cc.Invoke(ctx, "/Cart/DeleteAdvByIDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cartClient) AppendAdvByIDs(ctx context.Context, in *UserIdAdvertIdRequest, opts ...grpc.CallOption) (*AppendAdvResponse, error) {
	out := new(AppendAdvResponse)
	err := c.cc.Invoke(ctx, "/Cart/AppendAdvByIDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CartServer is the server API for Cart service.
// All implementations must embed UnimplementedCartServer
// for forward compatibility
type CartServer interface {
	GetCartByUserID(context.Context, *UserIdRequest) (*ReturningAdvertList, error)
	DeleteAdvByIDs(context.Context, *UserIdAdvertIdRequest) (*DeleteAdvResponse, error)
	AppendAdvByIDs(context.Context, *UserIdAdvertIdRequest) (*AppendAdvResponse, error)
	mustEmbedUnimplementedCartServer()
}

// UnimplementedCartServer must be embedded to have forward compatible implementations.
type UnimplementedCartServer struct {
}

func (UnimplementedCartServer) GetCartByUserID(context.Context, *UserIdRequest) (*ReturningAdvertList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCartByUserID not implemented")
}
func (UnimplementedCartServer) DeleteAdvByIDs(context.Context, *UserIdAdvertIdRequest) (*DeleteAdvResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAdvByIDs not implemented")
}
func (UnimplementedCartServer) AppendAdvByIDs(context.Context, *UserIdAdvertIdRequest) (*AppendAdvResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendAdvByIDs not implemented")
}
func (UnimplementedCartServer) mustEmbedUnimplementedCartServer() {}

// UnsafeCartServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CartServer will
// result in compilation errors.
type UnsafeCartServer interface {
	mustEmbedUnimplementedCartServer()
}

func RegisterCartServer(s grpc.ServiceRegistrar, srv CartServer) {
	s.RegisterService(&Cart_ServiceDesc, srv)
}

func _Cart_GetCartByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CartServer).GetCartByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Cart/GetCartByUserID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CartServer).GetCartByUserID(ctx, req.(*UserIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cart_DeleteAdvByIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIdAdvertIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CartServer).DeleteAdvByIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Cart/DeleteAdvByIDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CartServer).DeleteAdvByIDs(ctx, req.(*UserIdAdvertIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cart_AppendAdvByIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIdAdvertIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CartServer).AppendAdvByIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Cart/AppendAdvByIDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CartServer).AppendAdvByIDs(ctx, req.(*UserIdAdvertIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Cart_ServiceDesc is the grpc.ServiceDesc for Cart service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cart_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Cart",
	HandlerType: (*CartServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCartByUserID",
			Handler:    _Cart_GetCartByUserID_Handler,
		},
		{
			MethodName: "DeleteAdvByIDs",
			Handler:    _Cart_DeleteAdvByIDs_Handler,
		},
		{
			MethodName: "AppendAdvByIDs",
			Handler:    _Cart_AppendAdvByIDs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cart.proto",
}