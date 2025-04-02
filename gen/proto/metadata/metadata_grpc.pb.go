// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: proto/metadata/metadata.proto

package metadata

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
	MetadataService_AddMetaInfo_FullMethodName    = "/proto.metadata.MetadataService/AddMetaInfo"
	MetadataService_RemoveMetaInfo_FullMethodName = "/proto.metadata.MetadataService/RemoveMetaInfo"
	MetadataService_GetMetaInfo_FullMethodName    = "/proto.metadata.MetadataService/GetMetaInfo"
)

// MetadataServiceClient is the client API for MetadataService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Service for managing metadata
type MetadataServiceClient interface {
	AddMetaInfo(ctx context.Context, in *AddMetaInfoRequest, opts ...grpc.CallOption) (*AddMetaInfoResponse, error)
	RemoveMetaInfo(ctx context.Context, in *RemoveMetaInfoRequest, opts ...grpc.CallOption) (*RemoveMetaInfoResponse, error)
	GetMetaInfo(ctx context.Context, in *GetMetaInfoRequest, opts ...grpc.CallOption) (*GetMetaInfoResponse, error)
}

type metadataServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMetadataServiceClient(cc grpc.ClientConnInterface) MetadataServiceClient {
	return &metadataServiceClient{cc}
}

func (c *metadataServiceClient) AddMetaInfo(ctx context.Context, in *AddMetaInfoRequest, opts ...grpc.CallOption) (*AddMetaInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddMetaInfoResponse)
	err := c.cc.Invoke(ctx, MetadataService_AddMetaInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metadataServiceClient) RemoveMetaInfo(ctx context.Context, in *RemoveMetaInfoRequest, opts ...grpc.CallOption) (*RemoveMetaInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RemoveMetaInfoResponse)
	err := c.cc.Invoke(ctx, MetadataService_RemoveMetaInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metadataServiceClient) GetMetaInfo(ctx context.Context, in *GetMetaInfoRequest, opts ...grpc.CallOption) (*GetMetaInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetMetaInfoResponse)
	err := c.cc.Invoke(ctx, MetadataService_GetMetaInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetadataServiceServer is the server API for MetadataService service.
// All implementations must embed UnimplementedMetadataServiceServer
// for forward compatibility.
//
// Service for managing metadata
type MetadataServiceServer interface {
	AddMetaInfo(context.Context, *AddMetaInfoRequest) (*AddMetaInfoResponse, error)
	RemoveMetaInfo(context.Context, *RemoveMetaInfoRequest) (*RemoveMetaInfoResponse, error)
	GetMetaInfo(context.Context, *GetMetaInfoRequest) (*GetMetaInfoResponse, error)
	mustEmbedUnimplementedMetadataServiceServer()
}

// UnimplementedMetadataServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMetadataServiceServer struct{}

func (UnimplementedMetadataServiceServer) AddMetaInfo(context.Context, *AddMetaInfoRequest) (*AddMetaInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMetaInfo not implemented")
}
func (UnimplementedMetadataServiceServer) RemoveMetaInfo(context.Context, *RemoveMetaInfoRequest) (*RemoveMetaInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveMetaInfo not implemented")
}
func (UnimplementedMetadataServiceServer) GetMetaInfo(context.Context, *GetMetaInfoRequest) (*GetMetaInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetaInfo not implemented")
}
func (UnimplementedMetadataServiceServer) mustEmbedUnimplementedMetadataServiceServer() {}
func (UnimplementedMetadataServiceServer) testEmbeddedByValue()                         {}

// UnsafeMetadataServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetadataServiceServer will
// result in compilation errors.
type UnsafeMetadataServiceServer interface {
	mustEmbedUnimplementedMetadataServiceServer()
}

func RegisterMetadataServiceServer(s grpc.ServiceRegistrar, srv MetadataServiceServer) {
	// If the following call pancis, it indicates UnimplementedMetadataServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MetadataService_ServiceDesc, srv)
}

func _MetadataService_AddMetaInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMetaInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetadataServiceServer).AddMetaInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetadataService_AddMetaInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetadataServiceServer).AddMetaInfo(ctx, req.(*AddMetaInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetadataService_RemoveMetaInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveMetaInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetadataServiceServer).RemoveMetaInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetadataService_RemoveMetaInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetadataServiceServer).RemoveMetaInfo(ctx, req.(*RemoveMetaInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetadataService_GetMetaInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetaInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetadataServiceServer).GetMetaInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetadataService_GetMetaInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetadataServiceServer).GetMetaInfo(ctx, req.(*GetMetaInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MetadataService_ServiceDesc is the grpc.ServiceDesc for MetadataService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetadataService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.metadata.MetadataService",
	HandlerType: (*MetadataServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddMetaInfo",
			Handler:    _MetadataService_AddMetaInfo_Handler,
		},
		{
			MethodName: "RemoveMetaInfo",
			Handler:    _MetadataService_RemoveMetaInfo_Handler,
		},
		{
			MethodName: "GetMetaInfo",
			Handler:    _MetadataService_GetMetaInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/metadata/metadata.proto",
}
