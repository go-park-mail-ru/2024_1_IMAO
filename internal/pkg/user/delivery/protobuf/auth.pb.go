//protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v4.25.1
// source: auth.proto

package grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ExistedUserData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email    string `protobuf:"bytes,1,opt,name=Email,proto3" json:"Email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=Password,proto3" json:"Password,omitempty"`
}

func (x *ExistedUserData) Reset() {
	*x = ExistedUserData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExistedUserData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExistedUserData) ProtoMessage() {}

func (x *ExistedUserData) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExistedUserData.ProtoReflect.Descriptor instead.
func (*ExistedUserData) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{0}
}

func (x *ExistedUserData) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *ExistedUserData) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type NewUserData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email          string `protobuf:"bytes,1,opt,name=Email,proto3" json:"Email,omitempty"`
	Password       string `protobuf:"bytes,2,opt,name=Password,proto3" json:"Password,omitempty"`
	PasswordRepeat string `protobuf:"bytes,3,opt,name=PasswordRepeat,proto3" json:"PasswordRepeat,omitempty"`
}

func (x *NewUserData) Reset() {
	*x = NewUserData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewUserData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewUserData) ProtoMessage() {}

func (x *NewUserData) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewUserData.ProtoReflect.Descriptor instead.
func (*NewUserData) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{1}
}

func (x *NewUserData) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *NewUserData) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *NewUserData) GetPasswordRepeat() string {
	if x != nil {
		return x.PasswordRepeat
	}
	return ""
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID           uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Email        string `protobuf:"bytes,2,opt,name=Email,proto3" json:"Email,omitempty"`
	PasswordHash string `protobuf:"bytes,3,opt,name=PasswordHash,proto3" json:"PasswordHash,omitempty"`
	Avatar       string `protobuf:"bytes,4,opt,name=Avatar,proto3" json:"Avatar,omitempty"`
	IsAuth       bool   `protobuf:"varint,5,opt,name=IsAuth,proto3" json:"IsAuth,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{2}
}

func (x *User) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetPasswordHash() string {
	if x != nil {
		return x.PasswordHash
	}
	return ""
}

func (x *User) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

func (x *User) GetIsAuth() bool {
	if x != nil {
		return x.IsAuth
	}
	return false
}

type LoggedUser struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID           uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Email        string `protobuf:"bytes,2,opt,name=Email,proto3" json:"Email,omitempty"`
	PasswordHash string `protobuf:"bytes,3,opt,name=PasswordHash,proto3" json:"PasswordHash,omitempty"`
	SessionID    string `protobuf:"bytes,4,opt,name=SessionID,proto3" json:"SessionID,omitempty"`
	IsAuth       bool   `protobuf:"varint,5,opt,name=IsAuth,proto3" json:"IsAuth,omitempty"`
}

func (x *LoggedUser) Reset() {
	*x = LoggedUser{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoggedUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoggedUser) ProtoMessage() {}

func (x *LoggedUser) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoggedUser.ProtoReflect.Descriptor instead.
func (*LoggedUser) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{3}
}

func (x *LoggedUser) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *LoggedUser) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *LoggedUser) GetPasswordHash() string {
	if x != nil {
		return x.PasswordHash
	}
	return ""
}

func (x *LoggedUser) GetSessionID() string {
	if x != nil {
		return x.SessionID
	}
	return ""
}

func (x *LoggedUser) GetIsAuth() bool {
	if x != nil {
		return x.IsAuth
	}
	return false
}

type SessionData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionID string `protobuf:"bytes,1,opt,name=SessionID,proto3" json:"SessionID,omitempty"`
}

func (x *SessionData) Reset() {
	*x = SessionData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SessionData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SessionData) ProtoMessage() {}

func (x *SessionData) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SessionData.ProtoReflect.Descriptor instead.
func (*SessionData) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{4}
}

func (x *SessionData) GetSessionID() string {
	if x != nil {
		return x.SessionID
	}
	return ""
}

type EditEmailRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionID string `protobuf:"bytes,1,opt,name=SessionID,proto3" json:"SessionID,omitempty"`
	Email     string `protobuf:"bytes,2,opt,name=Email,proto3" json:"Email,omitempty"`
}

func (x *EditEmailRequest) Reset() {
	*x = EditEmailRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EditEmailRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EditEmailRequest) ProtoMessage() {}

func (x *EditEmailRequest) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EditEmailRequest.ProtoReflect.Descriptor instead.
func (*EditEmailRequest) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{5}
}

func (x *EditEmailRequest) GetSessionID() string {
	if x != nil {
		return x.SessionID
	}
	return ""
}

func (x *EditEmailRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type BoolResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=Ok,proto3" json:"Ok,omitempty"`
}

func (x *BoolResponse) Reset() {
	*x = BoolResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoolResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoolResponse) ProtoMessage() {}

func (x *BoolResponse) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoolResponse.ProtoReflect.Descriptor instead.
func (*BoolResponse) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{6}
}

func (x *BoolResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

var File_auth_proto protoreflect.FileDescriptor

var file_auth_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x43, 0x0a, 0x0f, 0x45, 0x78, 0x69,
	0x73, 0x74, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05,
	0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x6d, 0x61,
	0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x67,
	0x0a, 0x0b, 0x4e, 0x65, 0x77, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a,
	0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x6d,
	0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12,
	0x26, 0x0a, 0x0e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x70, 0x65, 0x61,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x22, 0x80, 0x01, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x44,
	0x12, 0x14, 0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x22, 0x0a, 0x0c, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x48, 0x61, 0x73, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x76,
	0x61, 0x74, 0x61, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x41, 0x76, 0x61, 0x74,
	0x61, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x49, 0x73, 0x41, 0x75, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x49, 0x73, 0x41, 0x75, 0x74, 0x68, 0x22, 0x8c, 0x01, 0x0a, 0x0a, 0x4c,
	0x6f, 0x67, 0x67, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x6d, 0x61,
	0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12,
	0x22, 0x0a, 0x0c, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x48, 0x61, 0x73, 0x68, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x48,
	0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x44,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49,
	0x44, 0x12, 0x16, 0x0a, 0x06, 0x49, 0x73, 0x41, 0x75, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x49, 0x73, 0x41, 0x75, 0x74, 0x68, 0x22, 0x2b, 0x0a, 0x0b, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x22, 0x46, 0x0a, 0x10, 0x45, 0x64, 0x69, 0x74, 0x45, 0x6d,
	0x61, 0x69, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x1e,
	0x0a, 0x0c, 0x42, 0x6f, 0x6f, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x4f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x4f, 0x6b, 0x32, 0xcc,
	0x01, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12, 0x26, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x12, 0x10, 0x2e, 0x45, 0x78, 0x69, 0x73, 0x74, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61,
	0x74, 0x61, 0x1a, 0x0b, 0x2e, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x12,
	0x23, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x12, 0x0c, 0x2e, 0x4e, 0x65, 0x77, 0x55,
	0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x0b, 0x2e, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x64,
	0x55, 0x73, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x06, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x12, 0x0c,
	0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x12, 0x20, 0x0a, 0x09, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x41, 0x75, 0x74,
	0x68, 0x12, 0x0c, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x1a,
	0x05, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x09, 0x45, 0x64, 0x69, 0x74, 0x45, 0x6d,
	0x61, 0x69, 0x6c, 0x12, 0x11, 0x2e, 0x45, 0x64, 0x69, 0x74, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x05, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x42, 0x23, 0x5a,
	0x21, 0x2e, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x2f, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_proto_rawDescOnce sync.Once
	file_auth_proto_rawDescData = file_auth_proto_rawDesc
)

func file_auth_proto_rawDescGZIP() []byte {
	file_auth_proto_rawDescOnce.Do(func() {
		file_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_proto_rawDescData)
	})
	return file_auth_proto_rawDescData
}

var file_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_auth_proto_goTypes = []interface{}{
	(*ExistedUserData)(nil),  // 0: ExistedUserData
	(*NewUserData)(nil),      // 1: NewUserData
	(*User)(nil),             // 2: User
	(*LoggedUser)(nil),       // 3: LoggedUser
	(*SessionData)(nil),      // 4: SessionData
	(*EditEmailRequest)(nil), // 5: EditEmailRequest
	(*BoolResponse)(nil),     // 6: BoolResponse
	(*emptypb.Empty)(nil),    // 7: google.protobuf.Empty
}
var file_auth_proto_depIdxs = []int32{
	0, // 0: Auth.Login:input_type -> ExistedUserData
	1, // 1: Auth.Signup:input_type -> NewUserData
	4, // 2: Auth.Logout:input_type -> SessionData
	4, // 3: Auth.CheckAuth:input_type -> SessionData
	5, // 4: Auth.EditEmail:input_type -> EditEmailRequest
	3, // 5: Auth.Login:output_type -> LoggedUser
	3, // 6: Auth.Signup:output_type -> LoggedUser
	7, // 7: Auth.Logout:output_type -> google.protobuf.Empty
	2, // 8: Auth.CheckAuth:output_type -> User
	2, // 9: Auth.EditEmail:output_type -> User
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_auth_proto_init() }
func file_auth_proto_init() {
	if File_auth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_auth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExistedUserData); i {
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
		file_auth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewUserData); i {
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
		file_auth_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_auth_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoggedUser); i {
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
		file_auth_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SessionData); i {
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
		file_auth_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EditEmailRequest); i {
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
		file_auth_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoolResponse); i {
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
			RawDescriptor: file_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_auth_proto_goTypes,
		DependencyIndexes: file_auth_proto_depIdxs,
		MessageInfos:      file_auth_proto_msgTypes,
	}.Build()
	File_auth_proto = out.File
	file_auth_proto_rawDesc = nil
	file_auth_proto_goTypes = nil
	file_auth_proto_depIdxs = nil
}
