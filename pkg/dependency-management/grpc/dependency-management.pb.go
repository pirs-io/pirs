// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: dependency-management.proto

package grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DetectRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProcessType int32  `protobuf:"varint,1,opt,name=process_type,json=processType,proto3" json:"process_type,omitempty"`
	ProjectUri  string `protobuf:"bytes,2,opt,name=project_uri,json=projectUri,proto3" json:"project_uri,omitempty"`
	// Types that are assignable to Data:
	//
	//	*DetectRequest_CountAndChecksum
	//	*DetectRequest_ChunkData
	Data isDetectRequest_Data `protobuf_oneof:"data"`
}

func (x *DetectRequest) Reset() {
	*x = DetectRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dependency_management_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DetectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DetectRequest) ProtoMessage() {}

func (x *DetectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dependency_management_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DetectRequest.ProtoReflect.Descriptor instead.
func (*DetectRequest) Descriptor() ([]byte, []int) {
	return file_dependency_management_proto_rawDescGZIP(), []int{0}
}

func (x *DetectRequest) GetProcessType() int32 {
	if x != nil {
		return x.ProcessType
	}
	return 0
}

func (x *DetectRequest) GetProjectUri() string {
	if x != nil {
		return x.ProjectUri
	}
	return ""
}

func (m *DetectRequest) GetData() isDetectRequest_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *DetectRequest) GetCountAndChecksum() string {
	if x, ok := x.GetData().(*DetectRequest_CountAndChecksum); ok {
		return x.CountAndChecksum
	}
	return ""
}

func (x *DetectRequest) GetChunkData() []byte {
	if x, ok := x.GetData().(*DetectRequest_ChunkData); ok {
		return x.ChunkData
	}
	return nil
}

type isDetectRequest_Data interface {
	isDetectRequest_Data()
}

type DetectRequest_CountAndChecksum struct {
	CountAndChecksum string `protobuf:"bytes,3,opt,name=countAndChecksum,proto3,oneof"`
}

type DetectRequest_ChunkData struct {
	ChunkData []byte `protobuf:"bytes,4,opt,name=chunk_data,json=chunkData,proto3,oneof"`
}

func (*DetectRequest_CountAndChecksum) isDetectRequest_Data() {}

func (*DetectRequest_ChunkData) isDetectRequest_Data() {}

type DetectResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message  string           `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Metadata *structpb.Struct `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *DetectResponse) Reset() {
	*x = DetectResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dependency_management_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DetectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DetectResponse) ProtoMessage() {}

func (x *DetectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dependency_management_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DetectResponse.ProtoReflect.Descriptor instead.
func (*DetectResponse) Descriptor() ([]byte, []int) {
	return file_dependency_management_proto_rawDescGZIP(), []int{1}
}

func (x *DetectResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *DetectResponse) GetMetadata() *structpb.Struct {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type ResolveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResolveUri string `protobuf:"bytes,1,opt,name=resolve_uri,json=resolveUri,proto3" json:"resolve_uri,omitempty"`
}

func (x *ResolveRequest) Reset() {
	*x = ResolveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dependency_management_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResolveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResolveRequest) ProtoMessage() {}

func (x *ResolveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_dependency_management_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResolveRequest.ProtoReflect.Descriptor instead.
func (*ResolveRequest) Descriptor() ([]byte, []int) {
	return file_dependency_management_proto_rawDescGZIP(), []int{2}
}

func (x *ResolveRequest) GetResolveUri() string {
	if x != nil {
		return x.ResolveUri
	}
	return ""
}

type ResolveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message  string           `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	Metadata *structpb.Struct `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *ResolveResponse) Reset() {
	*x = ResolveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dependency_management_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResolveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResolveResponse) ProtoMessage() {}

func (x *ResolveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_dependency_management_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResolveResponse.ProtoReflect.Descriptor instead.
func (*ResolveResponse) Descriptor() ([]byte, []int) {
	return file_dependency_management_proto_rawDescGZIP(), []int{3}
}

func (x *ResolveResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ResolveResponse) GetMetadata() *structpb.Struct {
	if x != nil {
		return x.Metadata
	}
	return nil
}

var File_dependency_management_proto protoreflect.FileDescriptor

var file_dependency_management_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x64, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63, 0x79, 0x2d, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67,
	0x72, 0x70, 0x63, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xaa, 0x01, 0x0a, 0x0d, 0x44, 0x65, 0x74, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x55, 0x72, 0x69, 0x12, 0x2c, 0x0a, 0x10, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x41, 0x6e, 0x64, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x10, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x41, 0x6e, 0x64, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x73, 0x75, 0x6d, 0x12, 0x1f, 0x0a, 0x0a, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x09, 0x63, 0x68, 0x75,
	0x6e, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x5f,
	0x0a, 0x0e, 0x44, 0x65, 0x74, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x33, 0x0a, 0x08, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x31, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x5f, 0x75, 0x72, 0x69,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x55,
	0x72, 0x69, 0x22, 0x60, 0x0a, 0x0f, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x33, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x32, 0x8f, 0x01, 0x0a, 0x14, 0x44, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65,
	0x6e, 0x63, 0x79, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x39, 0x0a,
	0x06, 0x44, 0x65, 0x74, 0x65, 0x63, 0x74, 0x12, 0x13, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x44,
	0x65, 0x74, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x44, 0x65, 0x74, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x3c, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x6f,
	0x6c, 0x76, 0x65, 0x12, 0x14, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x6c,
	0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x3b, 0x67, 0x72, 0x70,
	0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_dependency_management_proto_rawDescOnce sync.Once
	file_dependency_management_proto_rawDescData = file_dependency_management_proto_rawDesc
)

func file_dependency_management_proto_rawDescGZIP() []byte {
	file_dependency_management_proto_rawDescOnce.Do(func() {
		file_dependency_management_proto_rawDescData = protoimpl.X.CompressGZIP(file_dependency_management_proto_rawDescData)
	})
	return file_dependency_management_proto_rawDescData
}

var file_dependency_management_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_dependency_management_proto_goTypes = []interface{}{
	(*DetectRequest)(nil),   // 0: grpc.DetectRequest
	(*DetectResponse)(nil),  // 1: grpc.DetectResponse
	(*ResolveRequest)(nil),  // 2: grpc.ResolveRequest
	(*ResolveResponse)(nil), // 3: grpc.ResolveResponse
	(*structpb.Struct)(nil), // 4: google.protobuf.Struct
}
var file_dependency_management_proto_depIdxs = []int32{
	4, // 0: grpc.DetectResponse.metadata:type_name -> google.protobuf.Struct
	4, // 1: grpc.ResolveResponse.metadata:type_name -> google.protobuf.Struct
	0, // 2: grpc.DependencyManagement.Detect:input_type -> grpc.DetectRequest
	2, // 3: grpc.DependencyManagement.Resolve:input_type -> grpc.ResolveRequest
	1, // 4: grpc.DependencyManagement.Detect:output_type -> grpc.DetectResponse
	3, // 5: grpc.DependencyManagement.Resolve:output_type -> grpc.ResolveResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_dependency_management_proto_init() }
func file_dependency_management_proto_init() {
	if File_dependency_management_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dependency_management_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DetectRequest); i {
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
		file_dependency_management_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DetectResponse); i {
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
		file_dependency_management_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResolveRequest); i {
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
		file_dependency_management_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResolveResponse); i {
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
	file_dependency_management_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*DetectRequest_CountAndChecksum)(nil),
		(*DetectRequest_ChunkData)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_dependency_management_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_dependency_management_proto_goTypes,
		DependencyIndexes: file_dependency_management_proto_depIdxs,
		MessageInfos:      file_dependency_management_proto_msgTypes,
	}.Build()
	File_dependency_management_proto = out.File
	file_dependency_management_proto_rawDesc = nil
	file_dependency_management_proto_goTypes = nil
	file_dependency_management_proto_depIdxs = nil
}