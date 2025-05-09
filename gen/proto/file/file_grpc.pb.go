// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: proto/file/file.proto

package file

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
	FileService_UploadFileV1_FullMethodName   = "/proto.file.FileService/UploadFileV1"
	FileService_GetFileV1_FullMethodName      = "/proto.file.FileService/GetFileV1"
	FileService_GetFilesV1_FullMethodName     = "/proto.file.FileService/GetFilesV1"
	FileService_DownloadFileV1_FullMethodName = "/proto.file.FileService/DownloadFileV1"
	FileService_DeleteFileV1_FullMethodName   = "/proto.file.FileService/DeleteFileV1"
)

// FileServiceClient is the client API for FileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// FileService provides functionality to upload, download, manage metadata,
// and delete user files securely.
type FileServiceClient interface {
	// Upload a file using a client-streaming RPC.
	UploadFileV1(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadFileV1Request, UploadFileV1Response], error)
	// Retrieve metadata for a specific file by its ID.
	GetFileV1(ctx context.Context, in *GetFileV1Request, opts ...grpc.CallOption) (*GetFileV1Response, error)
	// Retrieve metadata for all stored files.
	GetFilesV1(ctx context.Context, in *GetFilesV1Request, opts ...grpc.CallOption) (*GetFilesV1Response, error)
	// Download a file using a server-streaming RPC.
	DownloadFileV1(ctx context.Context, in *DownloadFileV1Request, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadFileV1Response], error)
	// Delete a file by its ID.
	DeleteFileV1(ctx context.Context, in *DeleteFileV1Request, opts ...grpc.CallOption) (*DeleteFileV1Response, error)
}

type fileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileServiceClient(cc grpc.ClientConnInterface) FileServiceClient {
	return &fileServiceClient{cc}
}

func (c *fileServiceClient) UploadFileV1(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadFileV1Request, UploadFileV1Response], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FileService_ServiceDesc.Streams[0], FileService_UploadFileV1_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[UploadFileV1Request, UploadFileV1Response]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileService_UploadFileV1Client = grpc.ClientStreamingClient[UploadFileV1Request, UploadFileV1Response]

func (c *fileServiceClient) GetFileV1(ctx context.Context, in *GetFileV1Request, opts ...grpc.CallOption) (*GetFileV1Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFileV1Response)
	err := c.cc.Invoke(ctx, FileService_GetFileV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) GetFilesV1(ctx context.Context, in *GetFilesV1Request, opts ...grpc.CallOption) (*GetFilesV1Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFilesV1Response)
	err := c.cc.Invoke(ctx, FileService_GetFilesV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) DownloadFileV1(ctx context.Context, in *DownloadFileV1Request, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadFileV1Response], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FileService_ServiceDesc.Streams[1], FileService_DownloadFileV1_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DownloadFileV1Request, DownloadFileV1Response]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileService_DownloadFileV1Client = grpc.ServerStreamingClient[DownloadFileV1Response]

func (c *fileServiceClient) DeleteFileV1(ctx context.Context, in *DeleteFileV1Request, opts ...grpc.CallOption) (*DeleteFileV1Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteFileV1Response)
	err := c.cc.Invoke(ctx, FileService_DeleteFileV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileServiceServer is the server API for FileService service.
// All implementations must embed UnimplementedFileServiceServer
// for forward compatibility.
//
// FileService provides functionality to upload, download, manage metadata,
// and delete user files securely.
type FileServiceServer interface {
	// Upload a file using a client-streaming RPC.
	UploadFileV1(grpc.ClientStreamingServer[UploadFileV1Request, UploadFileV1Response]) error
	// Retrieve metadata for a specific file by its ID.
	GetFileV1(context.Context, *GetFileV1Request) (*GetFileV1Response, error)
	// Retrieve metadata for all stored files.
	GetFilesV1(context.Context, *GetFilesV1Request) (*GetFilesV1Response, error)
	// Download a file using a server-streaming RPC.
	DownloadFileV1(*DownloadFileV1Request, grpc.ServerStreamingServer[DownloadFileV1Response]) error
	// Delete a file by its ID.
	DeleteFileV1(context.Context, *DeleteFileV1Request) (*DeleteFileV1Response, error)
	mustEmbedUnimplementedFileServiceServer()
}

// UnimplementedFileServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFileServiceServer struct{}

func (UnimplementedFileServiceServer) UploadFileV1(grpc.ClientStreamingServer[UploadFileV1Request, UploadFileV1Response]) error {
	return status.Errorf(codes.Unimplemented, "method UploadFileV1 not implemented")
}
func (UnimplementedFileServiceServer) GetFileV1(context.Context, *GetFileV1Request) (*GetFileV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileV1 not implemented")
}
func (UnimplementedFileServiceServer) GetFilesV1(context.Context, *GetFilesV1Request) (*GetFilesV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFilesV1 not implemented")
}
func (UnimplementedFileServiceServer) DownloadFileV1(*DownloadFileV1Request, grpc.ServerStreamingServer[DownloadFileV1Response]) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFileV1 not implemented")
}
func (UnimplementedFileServiceServer) DeleteFileV1(context.Context, *DeleteFileV1Request) (*DeleteFileV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFileV1 not implemented")
}
func (UnimplementedFileServiceServer) mustEmbedUnimplementedFileServiceServer() {}
func (UnimplementedFileServiceServer) testEmbeddedByValue()                     {}

// UnsafeFileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileServiceServer will
// result in compilation errors.
type UnsafeFileServiceServer interface {
	mustEmbedUnimplementedFileServiceServer()
}

func RegisterFileServiceServer(s grpc.ServiceRegistrar, srv FileServiceServer) {
	// If the following call pancis, it indicates UnimplementedFileServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FileService_ServiceDesc, srv)
}

func _FileService_UploadFileV1_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileServiceServer).UploadFileV1(&grpc.GenericServerStream[UploadFileV1Request, UploadFileV1Response]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileService_UploadFileV1Server = grpc.ClientStreamingServer[UploadFileV1Request, UploadFileV1Response]

func _FileService_GetFileV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFileV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).GetFileV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_GetFileV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).GetFileV1(ctx, req.(*GetFileV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_GetFilesV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFilesV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).GetFilesV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_GetFilesV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).GetFilesV1(ctx, req.(*GetFilesV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_DownloadFileV1_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadFileV1Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileServiceServer).DownloadFileV1(m, &grpc.GenericServerStream[DownloadFileV1Request, DownloadFileV1Response]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileService_DownloadFileV1Server = grpc.ServerStreamingServer[DownloadFileV1Response]

func _FileService_DeleteFileV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFileV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).DeleteFileV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_DeleteFileV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).DeleteFileV1(ctx, req.(*DeleteFileV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

// FileService_ServiceDesc is the grpc.ServiceDesc for FileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.file.FileService",
	HandlerType: (*FileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFileV1",
			Handler:    _FileService_GetFileV1_Handler,
		},
		{
			MethodName: "GetFilesV1",
			Handler:    _FileService_GetFilesV1_Handler,
		},
		{
			MethodName: "DeleteFileV1",
			Handler:    _FileService_DeleteFileV1_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadFileV1",
			Handler:       _FileService_UploadFileV1_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadFileV1",
			Handler:       _FileService_DownloadFileV1_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/file/file.proto",
}
