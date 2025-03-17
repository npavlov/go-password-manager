// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: proto/password/password.proto

package password

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

type StorePasswordRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Login         string                 `protobuf:"bytes,2,opt,name=login,proto3" json:"login,omitempty"`
	Password      string                 `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StorePasswordRequest) Reset() {
	*x = StorePasswordRequest{}
	mi := &file_proto_password_password_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StorePasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StorePasswordRequest) ProtoMessage() {}

func (x *StorePasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_password_password_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StorePasswordRequest.ProtoReflect.Descriptor instead.
func (*StorePasswordRequest) Descriptor() ([]byte, []int) {
	return file_proto_password_password_proto_rawDescGZIP(), []int{0}
}

func (x *StorePasswordRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *StorePasswordRequest) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *StorePasswordRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type StorePasswordResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PasswordId    string                 `protobuf:"bytes,1,opt,name=password_id,json=passwordId,proto3" json:"password_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StorePasswordResponse) Reset() {
	*x = StorePasswordResponse{}
	mi := &file_proto_password_password_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StorePasswordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StorePasswordResponse) ProtoMessage() {}

func (x *StorePasswordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_password_password_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StorePasswordResponse.ProtoReflect.Descriptor instead.
func (*StorePasswordResponse) Descriptor() ([]byte, []int) {
	return file_proto_password_password_proto_rawDescGZIP(), []int{1}
}

func (x *StorePasswordResponse) GetPasswordId() string {
	if x != nil {
		return x.PasswordId
	}
	return ""
}

type GetPasswordsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPasswordsRequest) Reset() {
	*x = GetPasswordsRequest{}
	mi := &file_proto_password_password_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPasswordsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPasswordsRequest) ProtoMessage() {}

func (x *GetPasswordsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_password_password_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPasswordsRequest.ProtoReflect.Descriptor instead.
func (*GetPasswordsRequest) Descriptor() ([]byte, []int) {
	return file_proto_password_password_proto_rawDescGZIP(), []int{2}
}

type GetPasswordsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Passwords     []*PasswordMeta        `protobuf:"bytes,1,rep,name=passwords,proto3" json:"passwords,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPasswordsResponse) Reset() {
	*x = GetPasswordsResponse{}
	mi := &file_proto_password_password_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPasswordsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPasswordsResponse) ProtoMessage() {}

func (x *GetPasswordsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_password_password_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPasswordsResponse.ProtoReflect.Descriptor instead.
func (*GetPasswordsResponse) Descriptor() ([]byte, []int) {
	return file_proto_password_password_proto_rawDescGZIP(), []int{3}
}

func (x *GetPasswordsResponse) GetPasswords() []*PasswordMeta {
	if x != nil {
		return x.Passwords
	}
	return nil
}

type PasswordMeta struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PasswordId    string                 `protobuf:"bytes,1,opt,name=password_id,json=passwordId,proto3" json:"password_id,omitempty"`
	LastUpdate    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=last_update,json=lastUpdate,proto3" json:"last_update,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PasswordMeta) Reset() {
	*x = PasswordMeta{}
	mi := &file_proto_password_password_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PasswordMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PasswordMeta) ProtoMessage() {}

func (x *PasswordMeta) ProtoReflect() protoreflect.Message {
	mi := &file_proto_password_password_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PasswordMeta.ProtoReflect.Descriptor instead.
func (*PasswordMeta) Descriptor() ([]byte, []int) {
	return file_proto_password_password_proto_rawDescGZIP(), []int{4}
}

func (x *PasswordMeta) GetPasswordId() string {
	if x != nil {
		return x.PasswordId
	}
	return ""
}

func (x *PasswordMeta) GetLastUpdate() *timestamppb.Timestamp {
	if x != nil {
		return x.LastUpdate
	}
	return nil
}

type GetPasswordRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PasswordId    string                 `protobuf:"bytes,1,opt,name=password_id,json=passwordId,proto3" json:"password_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPasswordRequest) Reset() {
	*x = GetPasswordRequest{}
	mi := &file_proto_password_password_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPasswordRequest) ProtoMessage() {}

func (x *GetPasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_password_password_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPasswordRequest.ProtoReflect.Descriptor instead.
func (*GetPasswordRequest) Descriptor() ([]byte, []int) {
	return file_proto_password_password_proto_rawDescGZIP(), []int{5}
}

func (x *GetPasswordRequest) GetPasswordId() string {
	if x != nil {
		return x.PasswordId
	}
	return ""
}

type GetPasswordResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Password      *PasswordData          `protobuf:"bytes,1,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPasswordResponse) Reset() {
	*x = GetPasswordResponse{}
	mi := &file_proto_password_password_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPasswordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPasswordResponse) ProtoMessage() {}

func (x *GetPasswordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_password_password_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPasswordResponse.ProtoReflect.Descriptor instead.
func (*GetPasswordResponse) Descriptor() ([]byte, []int) {
	return file_proto_password_password_proto_rawDescGZIP(), []int{6}
}

func (x *GetPasswordResponse) GetPassword() *PasswordData {
	if x != nil {
		return x.Password
	}
	return nil
}

type UpdatePasswordRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PasswordId    string                 `protobuf:"bytes,1,opt,name=password_id,json=passwordId,proto3" json:"password_id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Login         string                 `protobuf:"bytes,3,opt,name=login,proto3" json:"login,omitempty"`
	Password      string                 `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdatePasswordRequest) Reset() {
	*x = UpdatePasswordRequest{}
	mi := &file_proto_password_password_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdatePasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePasswordRequest) ProtoMessage() {}

func (x *UpdatePasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_password_password_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePasswordRequest.ProtoReflect.Descriptor instead.
func (*UpdatePasswordRequest) Descriptor() ([]byte, []int) {
	return file_proto_password_password_proto_rawDescGZIP(), []int{7}
}

func (x *UpdatePasswordRequest) GetPasswordId() string {
	if x != nil {
		return x.PasswordId
	}
	return ""
}

func (x *UpdatePasswordRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdatePasswordRequest) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *UpdatePasswordRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type UpdatePasswordResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PasswordId    string                 `protobuf:"bytes,1,opt,name=password_id,json=passwordId,proto3" json:"password_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdatePasswordResponse) Reset() {
	*x = UpdatePasswordResponse{}
	mi := &file_proto_password_password_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdatePasswordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePasswordResponse) ProtoMessage() {}

func (x *UpdatePasswordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_password_password_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePasswordResponse.ProtoReflect.Descriptor instead.
func (*UpdatePasswordResponse) Descriptor() ([]byte, []int) {
	return file_proto_password_password_proto_rawDescGZIP(), []int{8}
}

func (x *UpdatePasswordResponse) GetPasswordId() string {
	if x != nil {
		return x.PasswordId
	}
	return ""
}

type PasswordData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Login         string                 `protobuf:"bytes,2,opt,name=login,proto3" json:"login,omitempty"`
	Password      string                 `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	LastUpdate    *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=last_update,json=lastUpdate,proto3" json:"last_update,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PasswordData) Reset() {
	*x = PasswordData{}
	mi := &file_proto_password_password_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PasswordData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PasswordData) ProtoMessage() {}

func (x *PasswordData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_password_password_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PasswordData.ProtoReflect.Descriptor instead.
func (*PasswordData) Descriptor() ([]byte, []int) {
	return file_proto_password_password_proto_rawDescGZIP(), []int{9}
}

func (x *PasswordData) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PasswordData) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *PasswordData) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *PasswordData) GetLastUpdate() *timestamppb.Timestamp {
	if x != nil {
		return x.LastUpdate
	}
	return nil
}

var File_proto_password_password_proto protoreflect.FileDescriptor

var file_proto_password_password_proto_rawDesc = string([]byte{
	0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x2f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x1a,
	0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x77, 0x0a,
	0x14, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x03, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x03, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69,
	0x6e, 0x12, 0x23, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x03, 0x52, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x38, 0x0a, 0x15, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x49, 0x64,
	0x22, 0x15, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x52, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3a, 0x0a, 0x09, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x4d, 0x65, 0x74, 0x61,
	0x52, 0x09, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x22, 0x6c, 0x0a, 0x0c, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x1f, 0x0a, 0x0b, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x49, 0x64, 0x12, 0x3b, 0x0a, 0x0b,
	0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x6c,
	0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x22, 0x3f, 0x0a, 0x12, 0x47, 0x65, 0x74,
	0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x29, 0x0a, 0x0b, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x0a,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x49, 0x64, 0x22, 0x4f, 0x0a, 0x13, 0x47, 0x65,
	0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x38, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0xa3, 0x01, 0x0a, 0x15,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29, 0x0a, 0x0b, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72,
	0x03, 0xb0, 0x01, 0x01, 0x52, 0x0a, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x49, 0x64,
	0x12, 0x1b, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x03, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a,
	0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48,
	0x04, 0x72, 0x02, 0x10, 0x03, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x23, 0x0a, 0x08,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x03, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x22, 0x39, 0x0a, 0x16, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x49, 0x64, 0x22, 0x91, 0x01, 0x0a,
	0x0c, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x12, 0x3b, 0x0a, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x32, 0x83, 0x03, 0x0a, 0x0f, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x5c, 0x0a, 0x0d, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x50, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x24, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x50, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2e, 0x53, 0x74, 0x6f,
	0x72, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x56, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x12, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x59, 0x0a, 0x0c, 0x47, 0x65,
	0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x23, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2e, 0x47, 0x65, 0x74, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x24, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x2e, 0x47, 0x65, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5f, 0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0xb3, 0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x42, 0x0d, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x35,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x70, 0x61, 0x76, 0x6c,
	0x6f, 0x76, 0x2f, 0x67, 0x6f, 0x2d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x2d, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0xa2, 0x02, 0x03, 0x50, 0x50, 0x58, 0xaa, 0x02, 0x0e, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0xca, 0x02, 0x0e, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0xe2, 0x02, 0x1a,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5c, 0x47,
	0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0f, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x3a, 0x3a, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_password_password_proto_rawDescOnce sync.Once
	file_proto_password_password_proto_rawDescData []byte
)

func file_proto_password_password_proto_rawDescGZIP() []byte {
	file_proto_password_password_proto_rawDescOnce.Do(func() {
		file_proto_password_password_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_password_password_proto_rawDesc), len(file_proto_password_password_proto_rawDesc)))
	})
	return file_proto_password_password_proto_rawDescData
}

var file_proto_password_password_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_proto_password_password_proto_goTypes = []any{
	(*StorePasswordRequest)(nil),   // 0: proto.password.StorePasswordRequest
	(*StorePasswordResponse)(nil),  // 1: proto.password.StorePasswordResponse
	(*GetPasswordsRequest)(nil),    // 2: proto.password.GetPasswordsRequest
	(*GetPasswordsResponse)(nil),   // 3: proto.password.GetPasswordsResponse
	(*PasswordMeta)(nil),           // 4: proto.password.PasswordMeta
	(*GetPasswordRequest)(nil),     // 5: proto.password.GetPasswordRequest
	(*GetPasswordResponse)(nil),    // 6: proto.password.GetPasswordResponse
	(*UpdatePasswordRequest)(nil),  // 7: proto.password.UpdatePasswordRequest
	(*UpdatePasswordResponse)(nil), // 8: proto.password.UpdatePasswordResponse
	(*PasswordData)(nil),           // 9: proto.password.PasswordData
	(*timestamppb.Timestamp)(nil),  // 10: google.protobuf.Timestamp
}
var file_proto_password_password_proto_depIdxs = []int32{
	4,  // 0: proto.password.GetPasswordsResponse.passwords:type_name -> proto.password.PasswordMeta
	10, // 1: proto.password.PasswordMeta.last_update:type_name -> google.protobuf.Timestamp
	9,  // 2: proto.password.GetPasswordResponse.password:type_name -> proto.password.PasswordData
	10, // 3: proto.password.PasswordData.last_update:type_name -> google.protobuf.Timestamp
	0,  // 4: proto.password.PasswordService.StorePassword:input_type -> proto.password.StorePasswordRequest
	5,  // 5: proto.password.PasswordService.GetPassword:input_type -> proto.password.GetPasswordRequest
	2,  // 6: proto.password.PasswordService.GetPasswords:input_type -> proto.password.GetPasswordsRequest
	7,  // 7: proto.password.PasswordService.UpdatePassword:input_type -> proto.password.UpdatePasswordRequest
	1,  // 8: proto.password.PasswordService.StorePassword:output_type -> proto.password.StorePasswordResponse
	6,  // 9: proto.password.PasswordService.GetPassword:output_type -> proto.password.GetPasswordResponse
	3,  // 10: proto.password.PasswordService.GetPasswords:output_type -> proto.password.GetPasswordsResponse
	8,  // 11: proto.password.PasswordService.UpdatePassword:output_type -> proto.password.UpdatePasswordResponse
	8,  // [8:12] is the sub-list for method output_type
	4,  // [4:8] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_proto_password_password_proto_init() }
func file_proto_password_password_proto_init() {
	if File_proto_password_password_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_password_password_proto_rawDesc), len(file_proto_password_password_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_password_password_proto_goTypes,
		DependencyIndexes: file_proto_password_password_proto_depIdxs,
		MessageInfos:      file_proto_password_password_proto_msgTypes,
	}.Build()
	File_proto_password_password_proto = out.File
	file_proto_password_password_proto_goTypes = nil
	file_proto_password_password_proto_depIdxs = nil
}
