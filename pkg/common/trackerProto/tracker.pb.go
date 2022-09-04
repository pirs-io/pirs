// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: tracker.proto

package trackerProto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type InstanceRegisterStatus int32

const (
	InstanceRegisterStatus_SUCCESS InstanceRegisterStatus = 0
	InstanceRegisterStatus_FAILED  InstanceRegisterStatus = 1
)

// Enum value maps for InstanceRegisterStatus.
var (
	InstanceRegisterStatus_name = map[int32]string{
		0: "SUCCESS",
		1: "FAILED",
	}
	InstanceRegisterStatus_value = map[string]int32{
		"SUCCESS": 0,
		"FAILED":  1,
	}
)

func (x InstanceRegisterStatus) Enum() *InstanceRegisterStatus {
	p := new(InstanceRegisterStatus)
	*p = x
	return p
}

func (x InstanceRegisterStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (InstanceRegisterStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_tracker_proto_enumTypes[0].Descriptor()
}

func (InstanceRegisterStatus) Type() protoreflect.EnumType {
	return &file_tracker_proto_enumTypes[0]
}

func (x InstanceRegisterStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use InstanceRegisterStatus.Descriptor instead.
func (InstanceRegisterStatus) EnumDescriptor() ([]byte, []int) {
	return file_tracker_proto_rawDescGZIP(), []int{0}
}

type PackageInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PackageInfo) Reset() {
	*x = PackageInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tracker_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PackageInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PackageInfo) ProtoMessage() {}

func (x *PackageInfo) ProtoReflect() protoreflect.Message {
	mi := &file_tracker_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PackageInfo.ProtoReflect.Descriptor instead.
func (*PackageInfo) Descriptor() ([]byte, []int) {
	return file_tracker_proto_rawDescGZIP(), []int{0}
}

type PackageRegisterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PackageRegisterResponse) Reset() {
	*x = PackageRegisterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tracker_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PackageRegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PackageRegisterResponse) ProtoMessage() {}

func (x *PackageRegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tracker_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PackageRegisterResponse.ProtoReflect.Descriptor instead.
func (*PackageRegisterResponse) Descriptor() ([]byte, []int) {
	return file_tracker_proto_rawDescGZIP(), []int{1}
}

type LocationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LocationRequest) Reset() {
	*x = LocationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tracker_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocationRequest) ProtoMessage() {}

func (x *LocationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tracker_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocationRequest.ProtoReflect.Descriptor instead.
func (*LocationRequest) Descriptor() ([]byte, []int) {
	return file_tracker_proto_rawDescGZIP(), []int{2}
}

type PackageLocation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PackageLocation) Reset() {
	*x = PackageLocation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tracker_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PackageLocation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PackageLocation) ProtoMessage() {}

func (x *PackageLocation) ProtoReflect() protoreflect.Message {
	mi := &file_tracker_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PackageLocation.ProtoReflect.Descriptor instead.
func (*PackageLocation) Descriptor() ([]byte, []int) {
	return file_tracker_proto_rawDescGZIP(), []int{3}
}

type TrackerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrganizationId       string   `protobuf:"bytes,1,opt,name=organizationId,proto3" json:"organizationId,omitempty"`
	AllowPublicAccess    bool     `protobuf:"varint,2,opt,name=allow_public_access,json=allowPublicAccess,proto3" json:"allow_public_access,omitempty"`
	PublicRegister       bool     `protobuf:"varint,3,opt,name=public_register,json=publicRegister,proto3" json:"public_register,omitempty"`
	PartnerOrganizations []string `protobuf:"bytes,4,rep,name=partnerOrganizations,proto3" json:"partnerOrganizations,omitempty"`
}

func (x *TrackerInfo) Reset() {
	*x = TrackerInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tracker_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TrackerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrackerInfo) ProtoMessage() {}

func (x *TrackerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_tracker_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrackerInfo.ProtoReflect.Descriptor instead.
func (*TrackerInfo) Descriptor() ([]byte, []int) {
	return file_tracker_proto_rawDescGZIP(), []int{4}
}

func (x *TrackerInfo) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *TrackerInfo) GetAllowPublicAccess() bool {
	if x != nil {
		return x.AllowPublicAccess
	}
	return false
}

func (x *TrackerInfo) GetPublicRegister() bool {
	if x != nil {
		return x.PublicRegister
	}
	return false
}

func (x *TrackerInfo) GetPartnerOrganizations() []string {
	if x != nil {
		return x.PartnerOrganizations
	}
	return nil
}

type InstanceRegisterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status InstanceRegisterStatus `protobuf:"varint,1,opt,name=status,proto3,enum=trackerProto.InstanceRegisterStatus" json:"status,omitempty"`
	Error  string                 `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *InstanceRegisterResponse) Reset() {
	*x = InstanceRegisterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tracker_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstanceRegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceRegisterResponse) ProtoMessage() {}

func (x *InstanceRegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tracker_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstanceRegisterResponse.ProtoReflect.Descriptor instead.
func (*InstanceRegisterResponse) Descriptor() ([]byte, []int) {
	return file_tracker_proto_rawDescGZIP(), []int{5}
}

func (x *InstanceRegisterResponse) GetStatus() InstanceRegisterStatus {
	if x != nil {
		return x.Status
	}
	return InstanceRegisterStatus_SUCCESS
}

func (x *InstanceRegisterResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_tracker_proto protoreflect.FileDescriptor

var file_tracker_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0c, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0d, 0x0a,
	0x0b, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x19, 0x0a, 0x17,
	0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x11, 0x0a, 0x0f, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x11, 0x0a, 0x0f, 0x50, 0x61,
	0x63, 0x6b, 0x61, 0x67, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xc2, 0x01,
	0x0a, 0x0b, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x26, 0x0a,
	0x0e, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x2e, 0x0a, 0x13, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x5f, 0x70,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x11, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x27, 0x0a, 0x0f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f,
	0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e,
	0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x32,
	0x0a, 0x14, 0x70, 0x61, 0x72, 0x74, 0x6e, 0x65, 0x72, 0x4f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x14, 0x70, 0x61,
	0x72, 0x74, 0x6e, 0x65, 0x72, 0x4f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x22, 0x6e, 0x0a, 0x18, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x24,
	0x2e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x2a, 0x31, 0x0a, 0x16, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07,
	0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x01, 0x32, 0x9a, 0x02, 0x0a, 0x07, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x65,
	0x72, 0x12, 0x5e, 0x0a, 0x17, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x72, 0x61,
	0x63, 0x6b, 0x65, 0x72, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x19, 0x2e, 0x74,
	0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x61, 0x63,
	0x6b, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x26, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65,
	0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x58, 0x0a, 0x12, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x65, 0x77,
	0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x12, 0x19, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65,
	0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x49, 0x6e,
	0x66, 0x6f, 0x1a, 0x25, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x13, 0x46,
	0x69, 0x6e, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1d, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1d, 0x2e, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0x00, 0x42, 0x33, 0x0a, 0x1b, 0x69, 0x6f, 0x2e, 0x70, 0x69, 0x72, 0x73, 0x2e, 0x64, 0x61,
	0x74, 0x61, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x6d, 0x67, 0x6d, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x0f, 0x2e, 0x2f, 0x3b, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x65, 0x72, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x88, 0x01, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tracker_proto_rawDescOnce sync.Once
	file_tracker_proto_rawDescData = file_tracker_proto_rawDesc
)

func file_tracker_proto_rawDescGZIP() []byte {
	file_tracker_proto_rawDescOnce.Do(func() {
		file_tracker_proto_rawDescData = protoimpl.X.CompressGZIP(file_tracker_proto_rawDescData)
	})
	return file_tracker_proto_rawDescData
}

var file_tracker_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_tracker_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_tracker_proto_goTypes = []interface{}{
	(InstanceRegisterStatus)(0),      // 0: trackerProto.InstanceRegisterStatus
	(*PackageInfo)(nil),              // 1: trackerProto.PackageInfo
	(*PackageRegisterResponse)(nil),  // 2: trackerProto.PackageRegisterResponse
	(*LocationRequest)(nil),          // 3: trackerProto.LocationRequest
	(*PackageLocation)(nil),          // 4: trackerProto.PackageLocation
	(*TrackerInfo)(nil),              // 5: trackerProto.TrackerInfo
	(*InstanceRegisterResponse)(nil), // 6: trackerProto.InstanceRegisterResponse
}
var file_tracker_proto_depIdxs = []int32{
	0, // 0: trackerProto.InstanceRegisterResponse.status:type_name -> trackerProto.InstanceRegisterStatus
	5, // 1: trackerProto.Tracker.RegisterTrackerInstance:input_type -> trackerProto.TrackerInfo
	1, // 2: trackerProto.Tracker.RegisterNewPackage:input_type -> trackerProto.PackageInfo
	3, // 3: trackerProto.Tracker.FindPackageLocation:input_type -> trackerProto.LocationRequest
	6, // 4: trackerProto.Tracker.RegisterTrackerInstance:output_type -> trackerProto.InstanceRegisterResponse
	2, // 5: trackerProto.Tracker.RegisterNewPackage:output_type -> trackerProto.PackageRegisterResponse
	4, // 6: trackerProto.Tracker.FindPackageLocation:output_type -> trackerProto.PackageLocation
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tracker_proto_init() }
func file_tracker_proto_init() {
	if File_tracker_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tracker_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PackageInfo); i {
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
		file_tracker_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PackageRegisterResponse); i {
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
		file_tracker_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LocationRequest); i {
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
		file_tracker_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PackageLocation); i {
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
		file_tracker_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TrackerInfo); i {
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
		file_tracker_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstanceRegisterResponse); i {
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
			RawDescriptor: file_tracker_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tracker_proto_goTypes,
		DependencyIndexes: file_tracker_proto_depIdxs,
		EnumInfos:         file_tracker_proto_enumTypes,
		MessageInfos:      file_tracker_proto_msgTypes,
	}.Build()
	File_tracker_proto = out.File
	file_tracker_proto_rawDesc = nil
	file_tracker_proto_goTypes = nil
	file_tracker_proto_depIdxs = nil
}
