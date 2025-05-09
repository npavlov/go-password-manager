// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: proto/file/file.proto

package file

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Chunked file upload request.
type UploadFileV1Request struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Name of the file being uploaded (1–255 characters).
	Filename string `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	// Chunk of file data (can be empty for signaling).
	Data          []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadFileV1Request) Reset() {
	*x = UploadFileV1Request{}
	mi := &file_proto_file_file_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadFileV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileV1Request) ProtoMessage() {}

func (x *UploadFileV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_file_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadFileV1Request.ProtoReflect.Descriptor instead.
func (*UploadFileV1Request) Descriptor() ([]byte, []int) {
	return file_proto_file_file_proto_rawDescGZIP(), []int{0}
}

func (x *UploadFileV1Request) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *UploadFileV1Request) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// Response after successfully uploading a file.
type UploadFileV1Response struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique ID of the uploaded file.
	FileId string `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	// Optional server message or status.
	Message       string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadFileV1Response) Reset() {
	*x = UploadFileV1Response{}
	mi := &file_proto_file_file_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadFileV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileV1Response) ProtoMessage() {}

func (x *UploadFileV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_file_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadFileV1Response.ProtoReflect.Descriptor instead.
func (*UploadFileV1Response) Descriptor() ([]byte, []int) {
	return file_proto_file_file_proto_rawDescGZIP(), []int{1}
}

func (x *UploadFileV1Response) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

func (x *UploadFileV1Response) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// Request to download a file.
type DownloadFileV1Request struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// ID of the file to download (UUID format).
	FileId        string `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DownloadFileV1Request) Reset() {
	*x = DownloadFileV1Request{}
	mi := &file_proto_file_file_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DownloadFileV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadFileV1Request) ProtoMessage() {}

func (x *DownloadFileV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_file_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadFileV1Request.ProtoReflect.Descriptor instead.
func (*DownloadFileV1Request) Descriptor() ([]byte, []int) {
	return file_proto_file_file_proto_rawDescGZIP(), []int{2}
}

func (x *DownloadFileV1Request) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

// Response streaming file data during download.
type DownloadFileV1Response struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Chunk of the file's binary data.
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	// Timestamp of the last update to the file.
	LastUpdate    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=last_update,json=lastUpdate,proto3" json:"last_update,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DownloadFileV1Response) Reset() {
	*x = DownloadFileV1Response{}
	mi := &file_proto_file_file_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DownloadFileV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadFileV1Response) ProtoMessage() {}

func (x *DownloadFileV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_file_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadFileV1Response.ProtoReflect.Descriptor instead.
func (*DownloadFileV1Response) Descriptor() ([]byte, []int) {
	return file_proto_file_file_proto_rawDescGZIP(), []int{3}
}

func (x *DownloadFileV1Response) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *DownloadFileV1Response) GetLastUpdate() *timestamppb.Timestamp {
	if x != nil {
		return x.LastUpdate
	}
	return nil
}

// Request to retrieve a specific file's metadata.
type GetFileV1Request struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// ID of the file (UUID format).
	FileId        string `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFileV1Request) Reset() {
	*x = GetFileV1Request{}
	mi := &file_proto_file_file_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFileV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileV1Request) ProtoMessage() {}

func (x *GetFileV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_file_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileV1Request.ProtoReflect.Descriptor instead.
func (*GetFileV1Request) Descriptor() ([]byte, []int) {
	return file_proto_file_file_proto_rawDescGZIP(), []int{4}
}

func (x *GetFileV1Request) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

// Response with metadata of a specific file.
type GetFileV1Response struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// File metadata.
	File *FileMeta `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	// Timestamp of the last update to the file.
	LastUpdate    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=last_update,json=lastUpdate,proto3" json:"last_update,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFileV1Response) Reset() {
	*x = GetFileV1Response{}
	mi := &file_proto_file_file_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFileV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileV1Response) ProtoMessage() {}

func (x *GetFileV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_file_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileV1Response.ProtoReflect.Descriptor instead.
func (*GetFileV1Response) Descriptor() ([]byte, []int) {
	return file_proto_file_file_proto_rawDescGZIP(), []int{5}
}

func (x *GetFileV1Response) GetFile() *FileMeta {
	if x != nil {
		return x.File
	}
	return nil
}

func (x *GetFileV1Response) GetLastUpdate() *timestamppb.Timestamp {
	if x != nil {
		return x.LastUpdate
	}
	return nil
}

// Request to retrieve metadata for all files.
type GetFilesV1Request struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFilesV1Request) Reset() {
	*x = GetFilesV1Request{}
	mi := &file_proto_file_file_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFilesV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFilesV1Request) ProtoMessage() {}

func (x *GetFilesV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_file_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFilesV1Request.ProtoReflect.Descriptor instead.
func (*GetFilesV1Request) Descriptor() ([]byte, []int) {
	return file_proto_file_file_proto_rawDescGZIP(), []int{6}
}

// Response containing metadata of all stored files.
type GetFilesV1Response struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// List of all file metadata entries.
	Notes         []*FileMeta `protobuf:"bytes,1,rep,name=notes,proto3" json:"notes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFilesV1Response) Reset() {
	*x = GetFilesV1Response{}
	mi := &file_proto_file_file_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFilesV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFilesV1Response) ProtoMessage() {}

func (x *GetFilesV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_file_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFilesV1Response.ProtoReflect.Descriptor instead.
func (*GetFilesV1Response) Descriptor() ([]byte, []int) {
	return file_proto_file_file_proto_rawDescGZIP(), []int{7}
}

func (x *GetFilesV1Response) GetNotes() []*FileMeta {
	if x != nil {
		return x.Notes
	}
	return nil
}

// Request to delete a file.
type DeleteFileV1Request struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// ID of the file to be deleted (UUID format).
	FileId        string `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteFileV1Request) Reset() {
	*x = DeleteFileV1Request{}
	mi := &file_proto_file_file_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteFileV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFileV1Request) ProtoMessage() {}

func (x *DeleteFileV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_file_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFileV1Request.ProtoReflect.Descriptor instead.
func (*DeleteFileV1Request) Descriptor() ([]byte, []int) {
	return file_proto_file_file_proto_rawDescGZIP(), []int{8}
}

func (x *DeleteFileV1Request) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

// Response after deleting a file.
type DeleteFileV1Response struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Indicates whether the delete operation was successful.
	Ok            bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteFileV1Response) Reset() {
	*x = DeleteFileV1Response{}
	mi := &file_proto_file_file_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteFileV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFileV1Response) ProtoMessage() {}

func (x *DeleteFileV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_file_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFileV1Response.ProtoReflect.Descriptor instead.
func (*DeleteFileV1Response) Descriptor() ([]byte, []int) {
	return file_proto_file_file_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteFileV1Response) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

// Metadata structure for a stored file.
type FileMeta struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique file identifier.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Original file name.
	FileName string `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	// Size of the file in bytes.
	FileSize int64 `protobuf:"varint,3,opt,name=file_size,json=fileSize,proto3" json:"file_size,omitempty"`
	// Optional URL to access or download the file.
	FileUrl       string `protobuf:"bytes,4,opt,name=file_url,json=fileUrl,proto3" json:"file_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FileMeta) Reset() {
	*x = FileMeta{}
	mi := &file_proto_file_file_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileMeta) ProtoMessage() {}

func (x *FileMeta) ProtoReflect() protoreflect.Message {
	mi := &file_proto_file_file_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileMeta.ProtoReflect.Descriptor instead.
func (*FileMeta) Descriptor() ([]byte, []int) {
	return file_proto_file_file_proto_rawDescGZIP(), []int{10}
}

func (x *FileMeta) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *FileMeta) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *FileMeta) GetFileSize() int64 {
	if x != nil {
		return x.FileSize
	}
	return 0
}

func (x *FileMeta) GetFileUrl() string {
	if x != nil {
		return x.FileUrl
	}
	return ""
}

var File_proto_file_file_proto protoreflect.FileDescriptor

var file_proto_file_file_proto_rawDesc = string([]byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2f, 0x66, 0x69, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x66,
	0x69, 0x6c, 0x65, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x51, 0x0a, 0x13, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x56,
	0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xba, 0x48, 0x07, 0x72,
	0x05, 0x10, 0x01, 0x18, 0xff, 0x01, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x49, 0x0a, 0x14, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69,
	0x6c, 0x65, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66,
	0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x3a, 0x0a, 0x15, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x56,
	0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03,
	0xb0, 0x01, 0x01, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x69, 0x0a, 0x16, 0x44,
	0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x56, 0x31, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x3b, 0x0a, 0x0b, 0x6c, 0x61, 0x73,
	0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x22, 0x35, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x07, 0x66, 0x69,
	0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05,
	0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x7a, 0x0a,
	0x11, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69,
	0x6c, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x3b, 0x0a, 0x0b,
	0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x6c,
	0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x22, 0x13, 0x0a, 0x11, 0x47, 0x65, 0x74,
	0x46, 0x69, 0x6c, 0x65, 0x73, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x40,
	0x0a, 0x12, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x05, 0x6e, 0x6f, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x66, 0x69, 0x6c, 0x65,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x05, 0x6e, 0x6f, 0x74, 0x65, 0x73,
	0x22, 0x38, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x56, 0x31,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0xb0,
	0x01, 0x01, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x26, 0x0a, 0x14, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02,
	0x6f, 0x6b, 0x22, 0x6f, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b,
	0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66,
	0x69, 0x6c, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x66, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65,
	0x5f, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x65,
	0x55, 0x72, 0x6c, 0x32, 0xa7, 0x03, 0x0a, 0x0b, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x53, 0x0a, 0x0c, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c,
	0x65, 0x56, 0x31, 0x12, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x66, 0x69, 0x6c, 0x65,
	0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x56, 0x31, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x66, 0x69, 0x6c,
	0x65, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x56, 0x31, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x12, 0x48, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x46,
	0x69, 0x6c, 0x65, 0x56, 0x31, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x66, 0x69,
	0x6c, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x66, 0x69, 0x6c, 0x65,
	0x2e, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x4b, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x56, 0x31,
	0x12, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x47, 0x65,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x47, 0x65, 0x74,
	0x46, 0x69, 0x6c, 0x65, 0x73, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x59, 0x0a, 0x0e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x56,
	0x31, 0x12, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x44,
	0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x56, 0x31, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x66, 0x69, 0x6c,
	0x65, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x56, 0x31,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x51, 0x0a, 0x0c, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x56, 0x31, 0x12, 0x1f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x69,
	0x6c, 0x65, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46,
	0x69, 0x6c, 0x65, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x95, 0x01,
	0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x66, 0x69, 0x6c, 0x65,
	0x42, 0x09, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x70, 0x61, 0x76, 0x6c, 0x6f,
	0x76, 0x2f, 0x67, 0x6f, 0x2d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2d, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0xa2, 0x02,
	0x03, 0x50, 0x46, 0x58, 0xaa, 0x02, 0x0a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x69, 0x6c,
	0x65, 0xca, 0x02, 0x0a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x46, 0x69, 0x6c, 0x65, 0xe2, 0x02,
	0x16, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x46, 0x69, 0x6c, 0x65, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x3a,
	0x3a, 0x46, 0x69, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_file_file_proto_rawDescOnce sync.Once
	file_proto_file_file_proto_rawDescData []byte
)

func file_proto_file_file_proto_rawDescGZIP() []byte {
	file_proto_file_file_proto_rawDescOnce.Do(func() {
		file_proto_file_file_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_file_file_proto_rawDesc), len(file_proto_file_file_proto_rawDesc)))
	})
	return file_proto_file_file_proto_rawDescData
}

var file_proto_file_file_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_proto_file_file_proto_goTypes = []any{
	(*UploadFileV1Request)(nil),    // 0: proto.file.UploadFileV1Request
	(*UploadFileV1Response)(nil),   // 1: proto.file.UploadFileV1Response
	(*DownloadFileV1Request)(nil),  // 2: proto.file.DownloadFileV1Request
	(*DownloadFileV1Response)(nil), // 3: proto.file.DownloadFileV1Response
	(*GetFileV1Request)(nil),       // 4: proto.file.GetFileV1Request
	(*GetFileV1Response)(nil),      // 5: proto.file.GetFileV1Response
	(*GetFilesV1Request)(nil),      // 6: proto.file.GetFilesV1Request
	(*GetFilesV1Response)(nil),     // 7: proto.file.GetFilesV1Response
	(*DeleteFileV1Request)(nil),    // 8: proto.file.DeleteFileV1Request
	(*DeleteFileV1Response)(nil),   // 9: proto.file.DeleteFileV1Response
	(*FileMeta)(nil),               // 10: proto.file.FileMeta
	(*timestamppb.Timestamp)(nil),  // 11: google.protobuf.Timestamp
}
var file_proto_file_file_proto_depIdxs = []int32{
	11, // 0: proto.file.DownloadFileV1Response.last_update:type_name -> google.protobuf.Timestamp
	10, // 1: proto.file.GetFileV1Response.file:type_name -> proto.file.FileMeta
	11, // 2: proto.file.GetFileV1Response.last_update:type_name -> google.protobuf.Timestamp
	10, // 3: proto.file.GetFilesV1Response.notes:type_name -> proto.file.FileMeta
	0,  // 4: proto.file.FileService.UploadFileV1:input_type -> proto.file.UploadFileV1Request
	4,  // 5: proto.file.FileService.GetFileV1:input_type -> proto.file.GetFileV1Request
	6,  // 6: proto.file.FileService.GetFilesV1:input_type -> proto.file.GetFilesV1Request
	2,  // 7: proto.file.FileService.DownloadFileV1:input_type -> proto.file.DownloadFileV1Request
	8,  // 8: proto.file.FileService.DeleteFileV1:input_type -> proto.file.DeleteFileV1Request
	1,  // 9: proto.file.FileService.UploadFileV1:output_type -> proto.file.UploadFileV1Response
	5,  // 10: proto.file.FileService.GetFileV1:output_type -> proto.file.GetFileV1Response
	7,  // 11: proto.file.FileService.GetFilesV1:output_type -> proto.file.GetFilesV1Response
	3,  // 12: proto.file.FileService.DownloadFileV1:output_type -> proto.file.DownloadFileV1Response
	9,  // 13: proto.file.FileService.DeleteFileV1:output_type -> proto.file.DeleteFileV1Response
	9,  // [9:14] is the sub-list for method output_type
	4,  // [4:9] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_proto_file_file_proto_init() }
func file_proto_file_file_proto_init() {
	if File_proto_file_file_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_file_file_proto_rawDesc), len(file_proto_file_file_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_file_file_proto_goTypes,
		DependencyIndexes: file_proto_file_file_proto_depIdxs,
		MessageInfos:      file_proto_file_file_proto_msgTypes,
	}.Build()
	File_proto_file_file_proto = out.File
	file_proto_file_file_proto_goTypes = nil
	file_proto_file_file_proto_depIdxs = nil
}
