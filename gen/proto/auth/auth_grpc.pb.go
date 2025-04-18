// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: proto/auth/auth.proto

package auth

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AuthService_RegisterV1_FullMethodName     = "/proto.auth.AuthService/RegisterV1"
	AuthService_LoginV1_FullMethodName        = "/proto.auth.AuthService/LoginV1"
	AuthService_RefreshTokenV1_FullMethodName = "/proto.auth.AuthService/RefreshTokenV1"
)

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// AuthService provides methods for user registration, authentication,
// and token refreshing.
type AuthServiceClient interface {
	// Register a new user and return tokens and a user key.
	RegisterV1(ctx context.Context, in *RegisterV1Request, opts ...grpc.CallOption) (*RegisterV1Response, error)
	// Authenticate a user with username and password, returning tokens.
	LoginV1(ctx context.Context, in *LoginV1Request, opts ...grpc.CallOption) (*LoginV1Response, error)
	// Refresh authentication tokens using a valid refresh token.
	RefreshTokenV1(ctx context.Context, in *RefreshTokenV1Request, opts ...grpc.CallOption) (*RefreshTokenV1Response, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) RegisterV1(ctx context.Context, in *RegisterV1Request, opts ...grpc.CallOption) (*RegisterV1Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterV1Response)
	err := c.cc.Invoke(ctx, AuthService_RegisterV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) LoginV1(ctx context.Context, in *LoginV1Request, opts ...grpc.CallOption) (*LoginV1Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginV1Response)
	err := c.cc.Invoke(ctx, AuthService_LoginV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) RefreshTokenV1(ctx context.Context, in *RefreshTokenV1Request, opts ...grpc.CallOption) (*RefreshTokenV1Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RefreshTokenV1Response)
	err := c.cc.Invoke(ctx, AuthService_RefreshTokenV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility.
//
// AuthService provides methods for user registration, authentication,
// and token refreshing.
type AuthServiceServer interface {
	// Register a new user and return tokens and a user key.
	RegisterV1(context.Context, *RegisterV1Request) (*RegisterV1Response, error)
	// Authenticate a user with username and password, returning tokens.
	LoginV1(context.Context, *LoginV1Request) (*LoginV1Response, error)
	// Refresh authentication tokens using a valid refresh token.
	RefreshTokenV1(context.Context, *RefreshTokenV1Request) (*RefreshTokenV1Response, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAuthServiceServer struct{}

func (UnimplementedAuthServiceServer) RegisterV1(context.Context, *RegisterV1Request) (*RegisterV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterV1 not implemented")
}
func (UnimplementedAuthServiceServer) LoginV1(context.Context, *LoginV1Request) (*LoginV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginV1 not implemented")
}
func (UnimplementedAuthServiceServer) RefreshTokenV1(context.Context, *RefreshTokenV1Request) (*RefreshTokenV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshTokenV1 not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}
func (UnimplementedAuthServiceServer) testEmbeddedByValue()                     {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	// If the following call pancis, it indicates UnimplementedAuthServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_RegisterV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).RegisterV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_RegisterV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).RegisterV1(ctx, req.(*RegisterV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_LoginV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).LoginV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_LoginV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).LoginV1(ctx, req.(*LoginV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_RefreshTokenV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).RefreshTokenV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_RefreshTokenV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).RefreshTokenV1(ctx, req.(*RefreshTokenV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.auth.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterV1",
			Handler:    _AuthService_RegisterV1_Handler,
		},
		{
			MethodName: "LoginV1",
			Handler:    _AuthService_LoginV1_Handler,
		},
		{
			MethodName: "RefreshTokenV1",
			Handler:    _AuthService_RefreshTokenV1_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/auth/auth.proto",
}
