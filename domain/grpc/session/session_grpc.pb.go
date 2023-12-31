// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: domain/grpc/session/session.proto

package session

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

// AuthCheckerClient is the client API for AuthChecker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthCheckerClient interface {
	IsAuth(ctx context.Context, in *Token, opts ...grpc.CallOption) (*UserID, error)
}

type authCheckerClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthCheckerClient(cc grpc.ClientConnInterface) AuthCheckerClient {
	return &authCheckerClient{cc}
}

func (c *authCheckerClient) IsAuth(ctx context.Context, in *Token, opts ...grpc.CallOption) (*UserID, error) {
	out := new(UserID)
	err := c.cc.Invoke(ctx, "/session.AuthChecker/IsAuth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthCheckerServer is the server API for AuthChecker service.
// All implementations must embed UnimplementedAuthCheckerServer
// for forward compatibility
type AuthCheckerServer interface {
	IsAuth(context.Context, *Token) (*UserID, error)
	mustEmbedUnimplementedAuthCheckerServer()
}

// UnimplementedAuthCheckerServer must be embedded to have forward compatible implementations.
type UnimplementedAuthCheckerServer struct {
}

func (UnimplementedAuthCheckerServer) IsAuth(context.Context, *Token) (*UserID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsAuth not implemented")
}
func (UnimplementedAuthCheckerServer) mustEmbedUnimplementedAuthCheckerServer() {}

// UnsafeAuthCheckerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthCheckerServer will
// result in compilation errors.
type UnsafeAuthCheckerServer interface {
	mustEmbedUnimplementedAuthCheckerServer()
}

func RegisterAuthCheckerServer(s grpc.ServiceRegistrar, srv AuthCheckerServer) {
	s.RegisterService(&AuthChecker_ServiceDesc, srv)
}

func _AuthChecker_IsAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Token)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthCheckerServer).IsAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/session.AuthChecker/IsAuth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthCheckerServer).IsAuth(ctx, req.(*Token))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthChecker_ServiceDesc is the grpc.ServiceDesc for AuthChecker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthChecker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "session.AuthChecker",
	HandlerType: (*AuthCheckerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsAuth",
			Handler:    _AuthChecker_IsAuth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "domain/grpc/session/session.proto",
}
