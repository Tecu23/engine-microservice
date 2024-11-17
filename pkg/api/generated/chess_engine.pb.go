// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.21.12
// source: pkg/api/chess_engine.proto

package generated

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

type MoveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Fen        string `protobuf:"bytes,2,opt,name=fen,proto3" json:"fen,omitempty"`
	EngineType string `protobuf:"bytes,3,opt,name=engine_type,json=engineType,proto3" json:"engine_type,omitempty"`
	Depth      int32  `protobuf:"varint,4,opt,name=depth,proto3" json:"depth,omitempty"`
}

func (x *MoveRequest) Reset() {
	*x = MoveRequest{}
	mi := &file_pkg_api_chess_engine_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MoveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MoveRequest) ProtoMessage() {}

func (x *MoveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_chess_engine_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MoveRequest.ProtoReflect.Descriptor instead.
func (*MoveRequest) Descriptor() ([]byte, []int) {
	return file_pkg_api_chess_engine_proto_rawDescGZIP(), []int{0}
}

func (x *MoveRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MoveRequest) GetFen() string {
	if x != nil {
		return x.Fen
	}
	return ""
}

func (x *MoveRequest) GetEngineType() string {
	if x != nil {
		return x.EngineType
	}
	return ""
}

func (x *MoveRequest) GetDepth() int32 {
	if x != nil {
		return x.Depth
	}
	return 0
}

type MoveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BestMove   string `protobuf:"bytes,1,opt,name=best_move,json=bestMove,proto3" json:"best_move,omitempty"`
	EngineInfo string `protobuf:"bytes,2,opt,name=engine_info,json=engineInfo,proto3" json:"engine_info,omitempty"`
}

func (x *MoveResponse) Reset() {
	*x = MoveResponse{}
	mi := &file_pkg_api_chess_engine_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MoveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MoveResponse) ProtoMessage() {}

func (x *MoveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_chess_engine_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MoveResponse.ProtoReflect.Descriptor instead.
func (*MoveResponse) Descriptor() ([]byte, []int) {
	return file_pkg_api_chess_engine_proto_rawDescGZIP(), []int{1}
}

func (x *MoveResponse) GetBestMove() string {
	if x != nil {
		return x.BestMove
	}
	return ""
}

func (x *MoveResponse) GetEngineInfo() string {
	if x != nil {
		return x.EngineInfo
	}
	return ""
}

var File_pkg_api_chess_engine_proto protoreflect.FileDescriptor

var file_pkg_api_chess_engine_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x68, 0x65, 0x73, 0x73, 0x5f,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x66, 0x0a, 0x0b,
	0x4d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x66,
	0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x66, 0x65, 0x6e, 0x12, 0x1f, 0x0a,
	0x0b, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x64, 0x65, 0x70, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x64,
	0x65, 0x70, 0x74, 0x68, 0x22, 0x4c, 0x0a, 0x0c, 0x4d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x65, 0x73, 0x74, 0x5f, 0x6d, 0x6f, 0x76,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x65, 0x73, 0x74, 0x4d, 0x6f, 0x76,
	0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x5f, 0x69, 0x6e, 0x66, 0x6f,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x49, 0x6e,
	0x66, 0x6f, 0x32, 0x3f, 0x0a, 0x0b, 0x43, 0x68, 0x65, 0x73, 0x73, 0x45, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x12, 0x30, 0x0a, 0x11, 0x43, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x42, 0x65,
	0x73, 0x74, 0x4d, 0x6f, 0x76, 0x65, 0x12, 0x0c, 0x2e, 0x4d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x4d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64,
	0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_api_chess_engine_proto_rawDescOnce sync.Once
	file_pkg_api_chess_engine_proto_rawDescData = file_pkg_api_chess_engine_proto_rawDesc
)

func file_pkg_api_chess_engine_proto_rawDescGZIP() []byte {
	file_pkg_api_chess_engine_proto_rawDescOnce.Do(func() {
		file_pkg_api_chess_engine_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_api_chess_engine_proto_rawDescData)
	})
	return file_pkg_api_chess_engine_proto_rawDescData
}

var file_pkg_api_chess_engine_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_api_chess_engine_proto_goTypes = []any{
	(*MoveRequest)(nil),  // 0: MoveRequest
	(*MoveResponse)(nil), // 1: MoveResponse
}
var file_pkg_api_chess_engine_proto_depIdxs = []int32{
	0, // 0: ChessEngine.CalculateBestMove:input_type -> MoveRequest
	1, // 1: ChessEngine.CalculateBestMove:output_type -> MoveResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_api_chess_engine_proto_init() }
func file_pkg_api_chess_engine_proto_init() {
	if File_pkg_api_chess_engine_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_api_chess_engine_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_api_chess_engine_proto_goTypes,
		DependencyIndexes: file_pkg_api_chess_engine_proto_depIdxs,
		MessageInfos:      file_pkg_api_chess_engine_proto_msgTypes,
	}.Build()
	File_pkg_api_chess_engine_proto = out.File
	file_pkg_api_chess_engine_proto_rawDesc = nil
	file_pkg_api_chess_engine_proto_goTypes = nil
	file_pkg_api_chess_engine_proto_depIdxs = nil
}