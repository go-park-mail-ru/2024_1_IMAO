//protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v4.25.1
// source: profile.proto

package grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ProfileIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *ProfileIDRequest) Reset() {
	*x = ProfileIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfileIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileIDRequest) ProtoMessage() {}

func (x *ProfileIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileIDRequest.ProtoReflect.Descriptor instead.
func (*ProfileIDRequest) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{0}
}

func (x *ProfileIDRequest) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

type ProfileData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID              uint64                 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	UserID          uint64                 `protobuf:"varint,2,opt,name=UserID,proto3" json:"UserID,omitempty"`
	Name            string                 `protobuf:"bytes,3,opt,name=Name,proto3" json:"Name,omitempty"`
	Surname         string                 `protobuf:"bytes,4,opt,name=Surname,proto3" json:"Surname,omitempty"`
	CityID          uint64                 `protobuf:"varint,5,opt,name=CityID,proto3" json:"CityID,omitempty"`
	CityName        string                 `protobuf:"bytes,6,opt,name=CityName,proto3" json:"CityName,omitempty"`
	Translation     string                 `protobuf:"bytes,7,opt,name=Translation,proto3" json:"Translation,omitempty"`
	Phone           string                 `protobuf:"bytes,8,opt,name=Phone,proto3" json:"Phone,omitempty"`
	Avatar          string                 `protobuf:"bytes,9,opt,name=Avatar,proto3" json:"Avatar,omitempty"`
	RegisterTime    *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=RegisterTime,proto3" json:"RegisterTime,omitempty"`
	Rating          float32                `protobuf:"fixed32,11,opt,name=Rating,proto3" json:"Rating,omitempty"`
	ReactionsCount  float32                `protobuf:"fixed32,12,opt,name=ReactionsCount,proto3" json:"ReactionsCount,omitempty"`
	Approved        bool                   `protobuf:"varint,13,opt,name=Approved,proto3" json:"Approved,omitempty"`
	MerchantsName   string                 `protobuf:"bytes,14,opt,name=MerchantsName,proto3" json:"MerchantsName,omitempty"`
	SubersCount     int64                  `protobuf:"varint,15,opt,name=SubersCount,proto3" json:"SubersCount,omitempty"`
	SubonsCount     int64                  `protobuf:"varint,16,opt,name=SubonsCount,proto3" json:"SubonsCount,omitempty"`
	AvatarIMG       string                 `protobuf:"bytes,17,opt,name=AvatarIMG,proto3" json:"AvatarIMG,omitempty"`
	ActiveAddsCount int64                  `protobuf:"varint,18,opt,name=ActiveAddsCount,proto3" json:"ActiveAddsCount,omitempty"`
	SoldAddsCount   int64                  `protobuf:"varint,19,opt,name=SoldAddsCount,proto3" json:"SoldAddsCount,omitempty"`
}

func (x *ProfileData) Reset() {
	*x = ProfileData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfileData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileData) ProtoMessage() {}

func (x *ProfileData) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileData.ProtoReflect.Descriptor instead.
func (*ProfileData) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{1}
}

func (x *ProfileData) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *ProfileData) GetUserID() uint64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *ProfileData) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ProfileData) GetSurname() string {
	if x != nil {
		return x.Surname
	}
	return ""
}

func (x *ProfileData) GetCityID() uint64 {
	if x != nil {
		return x.CityID
	}
	return 0
}

func (x *ProfileData) GetCityName() string {
	if x != nil {
		return x.CityName
	}
	return ""
}

func (x *ProfileData) GetTranslation() string {
	if x != nil {
		return x.Translation
	}
	return ""
}

func (x *ProfileData) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *ProfileData) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *ProfileData) GetRegisterTime() *timestamppb.Timestamp {
	if x != nil {
		return x.RegisterTime
	}
	return nil
}

func (x *ProfileData) GetRating() float32 {
	if x != nil {
		return x.Rating
	}
	return 0
}

func (x *ProfileData) GetReactionsCount() float32 {
	if x != nil {
		return x.ReactionsCount
	}
	return 0
}

func (x *ProfileData) GetApproved() bool {
	if x != nil {
		return x.Approved
	}
	return false
}

func (x *ProfileData) GetMerchantsName() string {
	if x != nil {
		return x.MerchantsName
	}
	return ""
}

func (x *ProfileData) GetSubersCount() int64 {
	if x != nil {
		return x.SubersCount
	}
	return 0
}

func (x *ProfileData) GetSubonsCount() int64 {
	if x != nil {
		return x.SubonsCount
	}
	return 0
}

func (x *ProfileData) GetAvatarIMG() string {
	if x != nil {
		return x.AvatarIMG
	}
	return ""
}

func (x *ProfileData) GetActiveAddsCount() int64 {
	if x != nil {
		return x.ActiveAddsCount
	}
	return 0
}

func (x *ProfileData) GetSoldAddsCount() int64 {
	if x != nil {
		return x.SoldAddsCount
	}
	return 0
}

type SetCityRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID          uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	CityID      uint64 `protobuf:"varint,2,opt,name=CityID,proto3" json:"CityID,omitempty"`
	CityName    string `protobuf:"bytes,3,opt,name=CityName,proto3" json:"CityName,omitempty"`
	Translation string `protobuf:"bytes,4,opt,name=Translation,proto3" json:"Translation,omitempty"`
}

func (x *SetCityRequest) Reset() {
	*x = SetCityRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetCityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetCityRequest) ProtoMessage() {}

func (x *SetCityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetCityRequest.ProtoReflect.Descriptor instead.
func (*SetCityRequest) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{2}
}

func (x *SetCityRequest) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *SetCityRequest) GetCityID() uint64 {
	if x != nil {
		return x.CityID
	}
	return 0
}

func (x *SetCityRequest) GetCityName() string {
	if x != nil {
		return x.CityName
	}
	return ""
}

func (x *SetCityRequest) GetTranslation() string {
	if x != nil {
		return x.Translation
	}
	return ""
}

type SetPhoneRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID    uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Phone string `protobuf:"bytes,2,opt,name=Phone,proto3" json:"Phone,omitempty"`
}

func (x *SetPhoneRequest) Reset() {
	*x = SetPhoneRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetPhoneRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetPhoneRequest) ProtoMessage() {}

func (x *SetPhoneRequest) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetPhoneRequest.ProtoReflect.Descriptor instead.
func (*SetPhoneRequest) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{3}
}

func (x *SetPhoneRequest) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *SetPhoneRequest) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

type EditProfileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID      uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Avatar  string `protobuf:"bytes,3,opt,name=Avatar,proto3" json:"Avatar,omitempty"`
	Surname string `protobuf:"bytes,4,opt,name=Surname,proto3" json:"Surname,omitempty"`
}

func (x *EditProfileRequest) Reset() {
	*x = EditProfileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_profile_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EditProfileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EditProfileRequest) ProtoMessage() {}

func (x *EditProfileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_profile_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EditProfileRequest.ProtoReflect.Descriptor instead.
func (*EditProfileRequest) Descriptor() ([]byte, []int) {
	return file_profile_proto_rawDescGZIP(), []int{4}
}

func (x *EditProfileRequest) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *EditProfileRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *EditProfileRequest) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *EditProfileRequest) GetSurname() string {
	if x != nil {
		return x.Surname
	}
	return ""
}

var File_profile_proto protoreflect.FileDescriptor

var file_profile_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x22, 0x0a, 0x10, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x49, 0x44, 0x22, 0xdb, 0x04, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x53, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x69,
	0x74, 0x79, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x43, 0x69, 0x74, 0x79,
	0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x43, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x14, 0x0a, 0x05, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x12, 0x3e,
	0x0a, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x0c, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x02, 0x52, 0x06,
	0x52, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x26, 0x0a, 0x0e, 0x52, 0x65, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0e,
	0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x41, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x08, 0x41, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x65, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x4d, 0x65,
	0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x73, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x53, 0x75, 0x62, 0x65, 0x72, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x0f, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x53, 0x75, 0x62, 0x65, 0x72, 0x73, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x75, 0x62, 0x6f, 0x6e, 0x73, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x10, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x53, 0x75, 0x62, 0x6f, 0x6e, 0x73, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x49, 0x4d,
	0x47, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x49,
	0x4d, 0x47, 0x12, 0x28, 0x0a, 0x0f, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x41, 0x64, 0x64, 0x73,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x12, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x41, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x41, 0x64, 0x64, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x0d,
	0x53, 0x6f, 0x6c, 0x64, 0x41, 0x64, 0x64, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x13, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0d, 0x53, 0x6f, 0x6c, 0x64, 0x41, 0x64, 0x64, 0x73, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x22, 0x76, 0x0a, 0x0e, 0x53, 0x65, 0x74, 0x43, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x69, 0x74, 0x79, 0x49, 0x44, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x43, 0x69, 0x74, 0x79, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08,
	0x43, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x43, 0x69, 0x74, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x37, 0x0a, 0x0f, 0x53, 0x65,
	0x74, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x44, 0x12, 0x14, 0x0a,
	0x05, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x50, 0x68,
	0x6f, 0x6e, 0x65, 0x22, 0x6a, 0x0a, 0x12, 0x45, 0x64, 0x69, 0x74, 0x50, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x41,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x53, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x32,
	0x80, 0x02, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x2d, 0x0a, 0x0a, 0x47,
	0x65, 0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x11, 0x2e, 0x50, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x50,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x30, 0x0a, 0x0d, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x11, 0x2e, 0x50, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c,
	0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x2f, 0x0a, 0x0e,
	0x53, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x43, 0x69, 0x74, 0x79, 0x12, 0x0f,
	0x2e, 0x53, 0x65, 0x74, 0x43, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x0c, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x31, 0x0a,
	0x0f, 0x53, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x68, 0x6f, 0x6e, 0x65,
	0x12, 0x10, 0x2e, 0x53, 0x65, 0x74, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x30, 0x0a, 0x0b, 0x45, 0x64, 0x69, 0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12,
	0x13, 0x2e, 0x45, 0x64, 0x69, 0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x44, 0x61,
	0x74, 0x61, 0x42, 0x26, 0x5a, 0x24, 0x2e, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2f, 0x64, 0x65, 0x6c,
	0x69, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_profile_proto_rawDescOnce sync.Once
	file_profile_proto_rawDescData = file_profile_proto_rawDesc
)

func file_profile_proto_rawDescGZIP() []byte {
	file_profile_proto_rawDescOnce.Do(func() {
		file_profile_proto_rawDescData = protoimpl.X.CompressGZIP(file_profile_proto_rawDescData)
	})
	return file_profile_proto_rawDescData
}

var file_profile_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_profile_proto_goTypes = []interface{}{
	(*ProfileIDRequest)(nil),      // 0: ProfileIDRequest
	(*ProfileData)(nil),           // 1: ProfileData
	(*SetCityRequest)(nil),        // 2: SetCityRequest
	(*SetPhoneRequest)(nil),       // 3: SetPhoneRequest
	(*EditProfileRequest)(nil),    // 4: EditProfileRequest
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_profile_proto_depIdxs = []int32{
	5, // 0: ProfileData.RegisterTime:type_name -> google.protobuf.Timestamp
	0, // 1: Profile.GetProfile:input_type -> ProfileIDRequest
	0, // 2: Profile.CreateProfile:input_type -> ProfileIDRequest
	2, // 3: Profile.SetProfileCity:input_type -> SetCityRequest
	3, // 4: Profile.SetProfilePhone:input_type -> SetPhoneRequest
	4, // 5: Profile.EditProfile:input_type -> EditProfileRequest
	1, // 6: Profile.GetProfile:output_type -> ProfileData
	1, // 7: Profile.CreateProfile:output_type -> ProfileData
	1, // 8: Profile.SetProfileCity:output_type -> ProfileData
	1, // 9: Profile.SetProfilePhone:output_type -> ProfileData
	1, // 10: Profile.EditProfile:output_type -> ProfileData
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_profile_proto_init() }
func file_profile_proto_init() {
	if File_profile_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_profile_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfileIDRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_profile_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfileData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_profile_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetCityRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_profile_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetPhoneRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_profile_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EditProfileRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_profile_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_profile_proto_goTypes,
		DependencyIndexes: file_profile_proto_depIdxs,
		MessageInfos:      file_profile_proto_msgTypes,
	}.Build()
	File_profile_proto = out.File
	file_profile_proto_rawDesc = nil
	file_profile_proto_goTypes = nil
	file_profile_proto_depIdxs = nil
}
