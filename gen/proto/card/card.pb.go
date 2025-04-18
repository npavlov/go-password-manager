// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: proto/card/card.proto

package card

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

// Request to store a new card.
type StoreCardV1Request struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Card data to be stored.
	Card          *CardData `protobuf:"bytes,1,opt,name=card,proto3" json:"card,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StoreCardV1Request) Reset() {
	*x = StoreCardV1Request{}
	mi := &file_proto_card_card_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StoreCardV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreCardV1Request) ProtoMessage() {}

func (x *StoreCardV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_card_card_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreCardV1Request.ProtoReflect.Descriptor instead.
func (*StoreCardV1Request) Descriptor() ([]byte, []int) {
	return file_proto_card_card_proto_rawDescGZIP(), []int{0}
}

func (x *StoreCardV1Request) GetCard() *CardData {
	if x != nil {
		return x.Card
	}
	return nil
}

// Response after storing a card.
type StoreCardV1Response struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// ID of the newly stored card.
	CardId        string `protobuf:"bytes,1,opt,name=card_id,json=cardId,proto3" json:"card_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StoreCardV1Response) Reset() {
	*x = StoreCardV1Response{}
	mi := &file_proto_card_card_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StoreCardV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreCardV1Response) ProtoMessage() {}

func (x *StoreCardV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_card_card_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreCardV1Response.ProtoReflect.Descriptor instead.
func (*StoreCardV1Response) Descriptor() ([]byte, []int) {
	return file_proto_card_card_proto_rawDescGZIP(), []int{1}
}

func (x *StoreCardV1Response) GetCardId() string {
	if x != nil {
		return x.CardId
	}
	return ""
}

// Request to retrieve all stored cards.
type GetCardsV1Request struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCardsV1Request) Reset() {
	*x = GetCardsV1Request{}
	mi := &file_proto_card_card_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCardsV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCardsV1Request) ProtoMessage() {}

func (x *GetCardsV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_card_card_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCardsV1Request.ProtoReflect.Descriptor instead.
func (*GetCardsV1Request) Descriptor() ([]byte, []int) {
	return file_proto_card_card_proto_rawDescGZIP(), []int{2}
}

// Response containing all stored cards.
type GetCardsV1Response struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// List of cards.
	Cards         []*CardData `protobuf:"bytes,1,rep,name=cards,proto3" json:"cards,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCardsV1Response) Reset() {
	*x = GetCardsV1Response{}
	mi := &file_proto_card_card_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCardsV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCardsV1Response) ProtoMessage() {}

func (x *GetCardsV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_card_card_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCardsV1Response.ProtoReflect.Descriptor instead.
func (*GetCardsV1Response) Descriptor() ([]byte, []int) {
	return file_proto_card_card_proto_rawDescGZIP(), []int{3}
}

func (x *GetCardsV1Response) GetCards() []*CardData {
	if x != nil {
		return x.Cards
	}
	return nil
}

// Request to retrieve a specific card.
type GetCardV1Request struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique card ID (UUID format).
	CardId        string `protobuf:"bytes,1,opt,name=card_id,json=cardId,proto3" json:"card_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCardV1Request) Reset() {
	*x = GetCardV1Request{}
	mi := &file_proto_card_card_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCardV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCardV1Request) ProtoMessage() {}

func (x *GetCardV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_card_card_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCardV1Request.ProtoReflect.Descriptor instead.
func (*GetCardV1Request) Descriptor() ([]byte, []int) {
	return file_proto_card_card_proto_rawDescGZIP(), []int{4}
}

func (x *GetCardV1Request) GetCardId() string {
	if x != nil {
		return x.CardId
	}
	return ""
}

// Response containing a single card and its last update timestamp.
type GetCardV1Response struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Card data.
	Card *CardData `protobuf:"bytes,1,opt,name=card,proto3" json:"card,omitempty"`
	// Timestamp of the last update.
	LastUpdate    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=last_update,json=lastUpdate,proto3" json:"last_update,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCardV1Response) Reset() {
	*x = GetCardV1Response{}
	mi := &file_proto_card_card_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCardV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCardV1Response) ProtoMessage() {}

func (x *GetCardV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_card_card_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCardV1Response.ProtoReflect.Descriptor instead.
func (*GetCardV1Response) Descriptor() ([]byte, []int) {
	return file_proto_card_card_proto_rawDescGZIP(), []int{5}
}

func (x *GetCardV1Response) GetCard() *CardData {
	if x != nil {
		return x.Card
	}
	return nil
}

func (x *GetCardV1Response) GetLastUpdate() *timestamppb.Timestamp {
	if x != nil {
		return x.LastUpdate
	}
	return nil
}

// Request to delete a card.
type DeleteCardV1Request struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique card ID (UUID format).
	CardId        string `protobuf:"bytes,1,opt,name=card_id,json=cardId,proto3" json:"card_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteCardV1Request) Reset() {
	*x = DeleteCardV1Request{}
	mi := &file_proto_card_card_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteCardV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCardV1Request) ProtoMessage() {}

func (x *DeleteCardV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_card_card_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCardV1Request.ProtoReflect.Descriptor instead.
func (*DeleteCardV1Request) Descriptor() ([]byte, []int) {
	return file_proto_card_card_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteCardV1Request) GetCardId() string {
	if x != nil {
		return x.CardId
	}
	return ""
}

// Response after attempting to delete a card.
type DeleteCardV1Response struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Status of the delete operation.
	Ok            bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteCardV1Response) Reset() {
	*x = DeleteCardV1Response{}
	mi := &file_proto_card_card_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteCardV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCardV1Response) ProtoMessage() {}

func (x *DeleteCardV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_card_card_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCardV1Response.ProtoReflect.Descriptor instead.
func (*DeleteCardV1Response) Descriptor() ([]byte, []int) {
	return file_proto_card_card_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteCardV1Response) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

// Request to update a card.
type UpdateCardV1Request struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique card ID (UUID format).
	CardId string `protobuf:"bytes,1,opt,name=card_id,json=cardId,proto3" json:"card_id,omitempty"`
	// Updated card data.
	Data          *CardData `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateCardV1Request) Reset() {
	*x = UpdateCardV1Request{}
	mi := &file_proto_card_card_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateCardV1Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCardV1Request) ProtoMessage() {}

func (x *UpdateCardV1Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_card_card_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCardV1Request.ProtoReflect.Descriptor instead.
func (*UpdateCardV1Request) Descriptor() ([]byte, []int) {
	return file_proto_card_card_proto_rawDescGZIP(), []int{8}
}

func (x *UpdateCardV1Request) GetCardId() string {
	if x != nil {
		return x.CardId
	}
	return ""
}

func (x *UpdateCardV1Request) GetData() *CardData {
	if x != nil {
		return x.Data
	}
	return nil
}

// Response after updating a card.
type UpdateCardV1Response struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// ID of the updated card.
	CardId        string `protobuf:"bytes,1,opt,name=card_id,json=cardId,proto3" json:"card_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateCardV1Response) Reset() {
	*x = UpdateCardV1Response{}
	mi := &file_proto_card_card_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateCardV1Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCardV1Response) ProtoMessage() {}

func (x *UpdateCardV1Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_card_card_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCardV1Response.ProtoReflect.Descriptor instead.
func (*UpdateCardV1Response) Descriptor() ([]byte, []int) {
	return file_proto_card_card_proto_rawDescGZIP(), []int{9}
}

func (x *UpdateCardV1Response) GetCardId() string {
	if x != nil {
		return x.CardId
	}
	return ""
}

// Represents the data structure for a payment card.
type CardData struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Card number (13 to 19 digits).
	CardNumber string `protobuf:"bytes,1,opt,name=card_number,json=cardNumber,proto3" json:"card_number,omitempty"`
	// Card expiry date in MM/YY format.
	ExpiryDate string `protobuf:"bytes,2,opt,name=expiry_date,json=expiryDate,proto3" json:"expiry_date,omitempty"`
	// Card security code (CVV), 3 or 4 digits.
	Cvv string `protobuf:"bytes,3,opt,name=cvv,proto3" json:"cvv,omitempty"`
	// Name of the cardholder (1 to 100 characters).
	CardholderName string `protobuf:"bytes,4,opt,name=cardholder_name,json=cardholderName,proto3" json:"cardholder_name,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *CardData) Reset() {
	*x = CardData{}
	mi := &file_proto_card_card_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CardData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CardData) ProtoMessage() {}

func (x *CardData) ProtoReflect() protoreflect.Message {
	mi := &file_proto_card_card_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CardData.ProtoReflect.Descriptor instead.
func (*CardData) Descriptor() ([]byte, []int) {
	return file_proto_card_card_proto_rawDescGZIP(), []int{10}
}

func (x *CardData) GetCardNumber() string {
	if x != nil {
		return x.CardNumber
	}
	return ""
}

func (x *CardData) GetExpiryDate() string {
	if x != nil {
		return x.ExpiryDate
	}
	return ""
}

func (x *CardData) GetCvv() string {
	if x != nil {
		return x.Cvv
	}
	return ""
}

func (x *CardData) GetCardholderName() string {
	if x != nil {
		return x.CardholderName
	}
	return ""
}

var File_proto_card_card_proto protoreflect.FileDescriptor

var file_proto_card_card_proto_rawDesc = string([]byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x61, 0x72, 0x64, 0x2f, 0x63, 0x61, 0x72,
	0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63,
	0x61, 0x72, 0x64, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x3e, 0x0a, 0x12, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x43, 0x61, 0x72, 0x64, 0x56, 0x31,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x28, 0x0a, 0x04, 0x63, 0x61, 0x72, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x61,
	0x72, 0x64, 0x2e, 0x43, 0x61, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x63, 0x61, 0x72,
	0x64, 0x22, 0x2e, 0x0a, 0x13, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x43, 0x61, 0x72, 0x64, 0x56, 0x31,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x61, 0x72, 0x64,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x61, 0x72, 0x64, 0x49,
	0x64, 0x22, 0x13, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x43, 0x61, 0x72, 0x64, 0x73, 0x56, 0x31, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x40, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x43, 0x61, 0x72,
	0x64, 0x73, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x05,
	0x63, 0x61, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e, 0x43, 0x61, 0x72, 0x64, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x05, 0x63, 0x61, 0x72, 0x64, 0x73, 0x22, 0x35, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x43,
	0x61, 0x72, 0x64, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x07,
	0x63, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba,
	0x48, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x06, 0x63, 0x61, 0x72, 0x64, 0x49, 0x64, 0x22,
	0x7a, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x43, 0x61, 0x72, 0x64, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x63, 0x61, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e,
	0x43, 0x61, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x63, 0x61, 0x72, 0x64, 0x12, 0x3b,
	0x0a, 0x0b, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0a, 0x6c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x22, 0x38, 0x0a, 0x13, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x72, 0x64, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x21, 0x0a, 0x07, 0x63, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x06, 0x63,
	0x61, 0x72, 0x64, 0x49, 0x64, 0x22, 0x26, 0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43,
	0x61, 0x72, 0x64, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x22, 0x62, 0x0a,
	0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x61, 0x72, 0x64, 0x56, 0x31, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x07, 0x63, 0x61, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52,
	0x06, 0x63, 0x61, 0x72, 0x64, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x61,
	0x72, 0x64, 0x2e, 0x43, 0x61, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x22, 0x2f, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x61, 0x72, 0x64, 0x56,
	0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x61, 0x72,
	0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x61, 0x72, 0x64,
	0x49, 0x64, 0x22, 0xe2, 0x01, 0x0a, 0x08, 0x43, 0x61, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x36, 0x0a, 0x0b, 0x63, 0x61, 0x72, 0x64, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x15, 0xba, 0x48, 0x12, 0x72, 0x10, 0x32, 0x0e, 0x5e, 0x5b, 0x30,
	0x2d, 0x39, 0x5d, 0x7b, 0x31, 0x33, 0x2c, 0x31, 0x39, 0x7d, 0x24, 0x52, 0x0a, 0x63, 0x61, 0x72,
	0x64, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x43, 0x0a, 0x0b, 0x65, 0x78, 0x70, 0x69, 0x72,
	0x79, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x22, 0xba, 0x48,
	0x1f, 0x72, 0x1d, 0x32, 0x1b, 0x5e, 0x28, 0x30, 0x5b, 0x31, 0x2d, 0x39, 0x5d, 0x7c, 0x31, 0x5b,
	0x30, 0x2d, 0x32, 0x5d, 0x29, 0x5c, 0x2f, 0x5b, 0x30, 0x2d, 0x39, 0x5d, 0x7b, 0x32, 0x7d, 0x24,
	0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x44, 0x61, 0x74, 0x65, 0x12, 0x25, 0x0a, 0x03,
	0x63, 0x76, 0x76, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x13, 0xba, 0x48, 0x10, 0x72, 0x0e,
	0x32, 0x0c, 0x5e, 0x5b, 0x30, 0x2d, 0x39, 0x5d, 0x7b, 0x33, 0x2c, 0x34, 0x7d, 0x24, 0x52, 0x03,
	0x63, 0x76, 0x76, 0x12, 0x32, 0x0a, 0x0f, 0x63, 0x61, 0x72, 0x64, 0x68, 0x6f, 0x6c, 0x64, 0x65,
	0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xba, 0x48,
	0x06, 0x72, 0x04, 0x10, 0x01, 0x18, 0x64, 0x52, 0x0e, 0x63, 0x61, 0x72, 0x64, 0x68, 0x6f, 0x6c,
	0x64, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x32, 0x9a, 0x03, 0x0a, 0x0b, 0x43, 0x61, 0x72, 0x64,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4e, 0x0a, 0x0b, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0x43, 0x61, 0x72, 0x64, 0x56, 0x31, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63,
	0x61, 0x72, 0x64, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x43, 0x61, 0x72, 0x64, 0x56, 0x31, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63,
	0x61, 0x72, 0x64, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x43, 0x61, 0x72, 0x64, 0x56, 0x31, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4b, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x43, 0x61,
	0x72, 0x64, 0x73, 0x56, 0x31, 0x12, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x61,
	0x72, 0x64, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x72, 0x64, 0x73, 0x56, 0x31, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x61, 0x72,
	0x64, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x61, 0x72, 0x64, 0x73, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x43, 0x61, 0x72, 0x64, 0x56,
	0x31, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e, 0x47,
	0x65, 0x74, 0x43, 0x61, 0x72, 0x64, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e, 0x47, 0x65, 0x74,
	0x43, 0x61, 0x72, 0x64, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51,
	0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x61, 0x72, 0x64, 0x56, 0x31, 0x12, 0x1f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x43, 0x61, 0x72, 0x64, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x43, 0x61, 0x72, 0x64, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x51, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x72, 0x64, 0x56,
	0x31, 0x12, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x72, 0x64, 0x56, 0x31, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x61, 0x72, 0x64, 0x56, 0x31, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x95, 0x01, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x42, 0x09, 0x43, 0x61, 0x72, 0x64, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6e, 0x70, 0x61, 0x76, 0x6c, 0x6f, 0x76, 0x2f, 0x67, 0x6f, 0x2d, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x67, 0x65, 0x6e,
	0x2f, 0x63, 0x61, 0x72, 0x64, 0xa2, 0x02, 0x03, 0x50, 0x43, 0x58, 0xaa, 0x02, 0x0a, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x61, 0x72, 0x64, 0xca, 0x02, 0x0a, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x5c, 0x43, 0x61, 0x72, 0x64, 0xe2, 0x02, 0x16, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x43, 0x61,
	0x72, 0x64, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x0b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x3a, 0x43, 0x61, 0x72, 0x64, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_card_card_proto_rawDescOnce sync.Once
	file_proto_card_card_proto_rawDescData []byte
)

func file_proto_card_card_proto_rawDescGZIP() []byte {
	file_proto_card_card_proto_rawDescOnce.Do(func() {
		file_proto_card_card_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_card_card_proto_rawDesc), len(file_proto_card_card_proto_rawDesc)))
	})
	return file_proto_card_card_proto_rawDescData
}

var file_proto_card_card_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_proto_card_card_proto_goTypes = []any{
	(*StoreCardV1Request)(nil),    // 0: proto.card.StoreCardV1Request
	(*StoreCardV1Response)(nil),   // 1: proto.card.StoreCardV1Response
	(*GetCardsV1Request)(nil),     // 2: proto.card.GetCardsV1Request
	(*GetCardsV1Response)(nil),    // 3: proto.card.GetCardsV1Response
	(*GetCardV1Request)(nil),      // 4: proto.card.GetCardV1Request
	(*GetCardV1Response)(nil),     // 5: proto.card.GetCardV1Response
	(*DeleteCardV1Request)(nil),   // 6: proto.card.DeleteCardV1Request
	(*DeleteCardV1Response)(nil),  // 7: proto.card.DeleteCardV1Response
	(*UpdateCardV1Request)(nil),   // 8: proto.card.UpdateCardV1Request
	(*UpdateCardV1Response)(nil),  // 9: proto.card.UpdateCardV1Response
	(*CardData)(nil),              // 10: proto.card.CardData
	(*timestamppb.Timestamp)(nil), // 11: google.protobuf.Timestamp
}
var file_proto_card_card_proto_depIdxs = []int32{
	10, // 0: proto.card.StoreCardV1Request.card:type_name -> proto.card.CardData
	10, // 1: proto.card.GetCardsV1Response.cards:type_name -> proto.card.CardData
	10, // 2: proto.card.GetCardV1Response.card:type_name -> proto.card.CardData
	11, // 3: proto.card.GetCardV1Response.last_update:type_name -> google.protobuf.Timestamp
	10, // 4: proto.card.UpdateCardV1Request.data:type_name -> proto.card.CardData
	0,  // 5: proto.card.CardService.StoreCardV1:input_type -> proto.card.StoreCardV1Request
	2,  // 6: proto.card.CardService.GetCardsV1:input_type -> proto.card.GetCardsV1Request
	4,  // 7: proto.card.CardService.GetCardV1:input_type -> proto.card.GetCardV1Request
	8,  // 8: proto.card.CardService.UpdateCardV1:input_type -> proto.card.UpdateCardV1Request
	6,  // 9: proto.card.CardService.DeleteCardV1:input_type -> proto.card.DeleteCardV1Request
	1,  // 10: proto.card.CardService.StoreCardV1:output_type -> proto.card.StoreCardV1Response
	3,  // 11: proto.card.CardService.GetCardsV1:output_type -> proto.card.GetCardsV1Response
	5,  // 12: proto.card.CardService.GetCardV1:output_type -> proto.card.GetCardV1Response
	9,  // 13: proto.card.CardService.UpdateCardV1:output_type -> proto.card.UpdateCardV1Response
	7,  // 14: proto.card.CardService.DeleteCardV1:output_type -> proto.card.DeleteCardV1Response
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_proto_card_card_proto_init() }
func file_proto_card_card_proto_init() {
	if File_proto_card_card_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_card_card_proto_rawDesc), len(file_proto_card_card_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_card_card_proto_goTypes,
		DependencyIndexes: file_proto_card_card_proto_depIdxs,
		MessageInfos:      file_proto_card_card_proto_msgTypes,
	}.Build()
	File_proto_card_card_proto = out.File
	file_proto_card_card_proto_goTypes = nil
	file_proto_card_card_proto_depIdxs = nil
}
