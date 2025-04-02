// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: proto/metadata/metadata.proto

package metadata

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

// Request to add or update metadata
type AddMetaInfoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ItemId        string                 `protobuf:"bytes,1,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"`                                                                 // The item to add metadata to
	Metadata      map[string]string      `protobuf:"bytes,2,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"` // Key-value metadata pairs
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddMetaInfoRequest) Reset() {
	*x = AddMetaInfoRequest{}
	mi := &file_proto_metadata_metadata_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddMetaInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddMetaInfoRequest) ProtoMessage() {}

func (x *AddMetaInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metadata_metadata_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddMetaInfoRequest.ProtoReflect.Descriptor instead.
func (*AddMetaInfoRequest) Descriptor() ([]byte, []int) {
	return file_proto_metadata_metadata_proto_rawDescGZIP(), []int{0}
}

func (x *AddMetaInfoRequest) GetItemId() string {
	if x != nil {
		return x.ItemId
	}
	return ""
}

func (x *AddMetaInfoRequest) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type AddMetaInfoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddMetaInfoResponse) Reset() {
	*x = AddMetaInfoResponse{}
	mi := &file_proto_metadata_metadata_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddMetaInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddMetaInfoResponse) ProtoMessage() {}

func (x *AddMetaInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metadata_metadata_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddMetaInfoResponse.ProtoReflect.Descriptor instead.
func (*AddMetaInfoResponse) Descriptor() ([]byte, []int) {
	return file_proto_metadata_metadata_proto_rawDescGZIP(), []int{1}
}

func (x *AddMetaInfoResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

// Request to remove metadata
type RemoveMetaInfoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ItemId        string                 `protobuf:"bytes,1,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"` // The item to remove metadata from
	Key           string                 `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`                     // The metadata key to delete
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RemoveMetaInfoRequest) Reset() {
	*x = RemoveMetaInfoRequest{}
	mi := &file_proto_metadata_metadata_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveMetaInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveMetaInfoRequest) ProtoMessage() {}

func (x *RemoveMetaInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metadata_metadata_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveMetaInfoRequest.ProtoReflect.Descriptor instead.
func (*RemoveMetaInfoRequest) Descriptor() ([]byte, []int) {
	return file_proto_metadata_metadata_proto_rawDescGZIP(), []int{2}
}

func (x *RemoveMetaInfoRequest) GetItemId() string {
	if x != nil {
		return x.ItemId
	}
	return ""
}

func (x *RemoveMetaInfoRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type RemoveMetaInfoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RemoveMetaInfoResponse) Reset() {
	*x = RemoveMetaInfoResponse{}
	mi := &file_proto_metadata_metadata_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveMetaInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveMetaInfoResponse) ProtoMessage() {}

func (x *RemoveMetaInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metadata_metadata_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveMetaInfoResponse.ProtoReflect.Descriptor instead.
func (*RemoveMetaInfoResponse) Descriptor() ([]byte, []int) {
	return file_proto_metadata_metadata_proto_rawDescGZIP(), []int{3}
}

func (x *RemoveMetaInfoResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

// Request to get metadata for an item
type GetMetaInfoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ItemId        string                 `protobuf:"bytes,1,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMetaInfoRequest) Reset() {
	*x = GetMetaInfoRequest{}
	mi := &file_proto_metadata_metadata_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMetaInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetaInfoRequest) ProtoMessage() {}

func (x *GetMetaInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metadata_metadata_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetaInfoRequest.ProtoReflect.Descriptor instead.
func (*GetMetaInfoRequest) Descriptor() ([]byte, []int) {
	return file_proto_metadata_metadata_proto_rawDescGZIP(), []int{4}
}

func (x *GetMetaInfoRequest) GetItemId() string {
	if x != nil {
		return x.ItemId
	}
	return ""
}

type GetMetaInfoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Metadata      map[string]string      `protobuf:"bytes,1,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMetaInfoResponse) Reset() {
	*x = GetMetaInfoResponse{}
	mi := &file_proto_metadata_metadata_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMetaInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetaInfoResponse) ProtoMessage() {}

func (x *GetMetaInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metadata_metadata_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetaInfoResponse.ProtoReflect.Descriptor instead.
func (*GetMetaInfoResponse) Descriptor() ([]byte, []int) {
	return file_proto_metadata_metadata_proto_rawDescGZIP(), []int{5}
}

func (x *GetMetaInfoResponse) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

var File_proto_metadata_metadata_proto protoreflect.FileDescriptor

var file_proto_metadata_metadata_proto_rawDesc = string([]byte{
	0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x1a,
	0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc2, 0x01, 0x0a,
	0x12, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x07, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x06,
	0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x4c, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74,
	0x61, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x2f, 0x0a, 0x13, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x22, 0x55, 0x0a, 0x15, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x4d, 0x65, 0x74, 0x61,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x07, 0x69,
	0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48,
	0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x19,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04,
	0x72, 0x02, 0x10, 0x01, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x32, 0x0a, 0x16, 0x52, 0x65, 0x6d,
	0x6f, 0x76, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x37, 0x0a,
	0x12, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x07, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x06,
	0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x22, 0xa1, 0x01, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x4d, 0x65,
	0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d,
	0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x3b, 0x0a,
	0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0xa2, 0x02, 0x0a, 0x0f, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x56,
	0x0a, 0x0b, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x22, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x41,
	0x64, 0x64, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x2e, 0x41, 0x64, 0x64, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5f, 0x0a, 0x0e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x56, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x4d, 0x65,
	0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x47, 0x65, 0x74, 0x4d,
	0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0xb3, 0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x42, 0x0d, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x70, 0x61, 0x76, 0x6c, 0x6f, 0x76, 0x2f, 0x67, 0x6f, 0x2d, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xa2, 0x02,
	0x03, 0x50, 0x4d, 0x58, 0xaa, 0x02, 0x0e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xca, 0x02, 0x0e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xe2, 0x02, 0x1a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0xea, 0x02, 0x0f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x3a, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_metadata_metadata_proto_rawDescOnce sync.Once
	file_proto_metadata_metadata_proto_rawDescData []byte
)

func file_proto_metadata_metadata_proto_rawDescGZIP() []byte {
	file_proto_metadata_metadata_proto_rawDescOnce.Do(func() {
		file_proto_metadata_metadata_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_metadata_metadata_proto_rawDesc), len(file_proto_metadata_metadata_proto_rawDesc)))
	})
	return file_proto_metadata_metadata_proto_rawDescData
}

var file_proto_metadata_metadata_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_metadata_metadata_proto_goTypes = []any{
	(*AddMetaInfoRequest)(nil),     // 0: proto.metadata.AddMetaInfoRequest
	(*AddMetaInfoResponse)(nil),    // 1: proto.metadata.AddMetaInfoResponse
	(*RemoveMetaInfoRequest)(nil),  // 2: proto.metadata.RemoveMetaInfoRequest
	(*RemoveMetaInfoResponse)(nil), // 3: proto.metadata.RemoveMetaInfoResponse
	(*GetMetaInfoRequest)(nil),     // 4: proto.metadata.GetMetaInfoRequest
	(*GetMetaInfoResponse)(nil),    // 5: proto.metadata.GetMetaInfoResponse
	nil,                            // 6: proto.metadata.AddMetaInfoRequest.MetadataEntry
	nil,                            // 7: proto.metadata.GetMetaInfoResponse.MetadataEntry
}
var file_proto_metadata_metadata_proto_depIdxs = []int32{
	6, // 0: proto.metadata.AddMetaInfoRequest.metadata:type_name -> proto.metadata.AddMetaInfoRequest.MetadataEntry
	7, // 1: proto.metadata.GetMetaInfoResponse.metadata:type_name -> proto.metadata.GetMetaInfoResponse.MetadataEntry
	0, // 2: proto.metadata.MetadataService.AddMetaInfo:input_type -> proto.metadata.AddMetaInfoRequest
	2, // 3: proto.metadata.MetadataService.RemoveMetaInfo:input_type -> proto.metadata.RemoveMetaInfoRequest
	4, // 4: proto.metadata.MetadataService.GetMetaInfo:input_type -> proto.metadata.GetMetaInfoRequest
	1, // 5: proto.metadata.MetadataService.AddMetaInfo:output_type -> proto.metadata.AddMetaInfoResponse
	3, // 6: proto.metadata.MetadataService.RemoveMetaInfo:output_type -> proto.metadata.RemoveMetaInfoResponse
	5, // 7: proto.metadata.MetadataService.GetMetaInfo:output_type -> proto.metadata.GetMetaInfoResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_metadata_metadata_proto_init() }
func file_proto_metadata_metadata_proto_init() {
	if File_proto_metadata_metadata_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_metadata_metadata_proto_rawDesc), len(file_proto_metadata_metadata_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_metadata_metadata_proto_goTypes,
		DependencyIndexes: file_proto_metadata_metadata_proto_depIdxs,
		MessageInfos:      file_proto_metadata_metadata_proto_msgTypes,
	}.Build()
	File_proto_metadata_metadata_proto = out.File
	file_proto_metadata_metadata_proto_goTypes = nil
	file_proto_metadata_metadata_proto_depIdxs = nil
}
