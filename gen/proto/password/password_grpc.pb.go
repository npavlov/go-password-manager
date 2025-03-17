// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: proto/password/password.proto

package password

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
	PasswordService_StorePassword_FullMethodName  = "/proto.password.PasswordService/StorePassword"
	PasswordService_GetPassword_FullMethodName    = "/proto.password.PasswordService/GetPassword"
	PasswordService_GetPasswords_FullMethodName   = "/proto.password.PasswordService/GetPasswords"
	PasswordService_UpdatePassword_FullMethodName = "/proto.password.PasswordService/UpdatePassword"
)

// PasswordServiceClient is the client API for PasswordService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PasswordServiceClient interface {
	StorePassword(ctx context.Context, in *StorePasswordRequest, opts ...grpc.CallOption) (*StorePasswordResponse, error)
	GetPassword(ctx context.Context, in *GetPasswordRequest, opts ...grpc.CallOption) (*GetPasswordResponse, error)
	GetPasswords(ctx context.Context, in *GetPasswordsRequest, opts ...grpc.CallOption) (*GetPasswordsResponse, error)
	UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*UpdatePasswordResponse, error)
}

type passwordServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPasswordServiceClient(cc grpc.ClientConnInterface) PasswordServiceClient {
	return &passwordServiceClient{cc}
}

func (c *passwordServiceClient) StorePassword(ctx context.Context, in *StorePasswordRequest, opts ...grpc.CallOption) (*StorePasswordResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StorePasswordResponse)
	err := c.cc.Invoke(ctx, PasswordService_StorePassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *passwordServiceClient) GetPassword(ctx context.Context, in *GetPasswordRequest, opts ...grpc.CallOption) (*GetPasswordResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPasswordResponse)
	err := c.cc.Invoke(ctx, PasswordService_GetPassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *passwordServiceClient) GetPasswords(ctx context.Context, in *GetPasswordsRequest, opts ...grpc.CallOption) (*GetPasswordsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPasswordsResponse)
	err := c.cc.Invoke(ctx, PasswordService_GetPasswords_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *passwordServiceClient) UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*UpdatePasswordResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdatePasswordResponse)
	err := c.cc.Invoke(ctx, PasswordService_UpdatePassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PasswordServiceServer is the server API for PasswordService service.
// All implementations must embed UnimplementedPasswordServiceServer
// for forward compatibility.
type PasswordServiceServer interface {
	StorePassword(context.Context, *StorePasswordRequest) (*StorePasswordResponse, error)
	GetPassword(context.Context, *GetPasswordRequest) (*GetPasswordResponse, error)
	GetPasswords(context.Context, *GetPasswordsRequest) (*GetPasswordsResponse, error)
	UpdatePassword(context.Context, *UpdatePasswordRequest) (*UpdatePasswordResponse, error)
	mustEmbedUnimplementedPasswordServiceServer()
}

// UnimplementedPasswordServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPasswordServiceServer struct{}

func (UnimplementedPasswordServiceServer) StorePassword(context.Context, *StorePasswordRequest) (*StorePasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StorePassword not implemented")
}
func (UnimplementedPasswordServiceServer) GetPassword(context.Context, *GetPasswordRequest) (*GetPasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPassword not implemented")
}
func (UnimplementedPasswordServiceServer) GetPasswords(context.Context, *GetPasswordsRequest) (*GetPasswordsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPasswords not implemented")
}
func (UnimplementedPasswordServiceServer) UpdatePassword(context.Context, *UpdatePasswordRequest) (*UpdatePasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePassword not implemented")
}
func (UnimplementedPasswordServiceServer) mustEmbedUnimplementedPasswordServiceServer() {}
func (UnimplementedPasswordServiceServer) testEmbeddedByValue()                         {}

// UnsafePasswordServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PasswordServiceServer will
// result in compilation errors.
type UnsafePasswordServiceServer interface {
	mustEmbedUnimplementedPasswordServiceServer()
}

func RegisterPasswordServiceServer(s grpc.ServiceRegistrar, srv PasswordServiceServer) {
	// If the following call pancis, it indicates UnimplementedPasswordServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PasswordService_ServiceDesc, srv)
}

func _PasswordService_StorePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StorePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PasswordServiceServer).StorePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PasswordService_StorePassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PasswordServiceServer).StorePassword(ctx, req.(*StorePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PasswordService_GetPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PasswordServiceServer).GetPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PasswordService_GetPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PasswordServiceServer).GetPassword(ctx, req.(*GetPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PasswordService_GetPasswords_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPasswordsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PasswordServiceServer).GetPasswords(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PasswordService_GetPasswords_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PasswordServiceServer).GetPasswords(ctx, req.(*GetPasswordsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PasswordService_UpdatePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PasswordServiceServer).UpdatePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PasswordService_UpdatePassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PasswordServiceServer).UpdatePassword(ctx, req.(*UpdatePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PasswordService_ServiceDesc is the grpc.ServiceDesc for PasswordService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PasswordService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.password.PasswordService",
	HandlerType: (*PasswordServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StorePassword",
			Handler:    _PasswordService_StorePassword_Handler,
		},
		{
			MethodName: "GetPassword",
			Handler:    _PasswordService_GetPassword_Handler,
		},
		{
			MethodName: "GetPasswords",
			Handler:    _PasswordService_GetPasswords_Handler,
		},
		{
			MethodName: "UpdatePassword",
			Handler:    _PasswordService_UpdatePassword_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/password/password.proto",
}
