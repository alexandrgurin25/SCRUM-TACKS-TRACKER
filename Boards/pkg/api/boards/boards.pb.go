// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: api/boards/boards.proto

package boards

import (
	types "boards/pkg/api/types"
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

type CreateBoardRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Title         string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateBoardRequest) Reset() {
	*x = CreateBoardRequest{}
	mi := &file_api_boards_boards_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateBoardRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateBoardRequest) ProtoMessage() {}

func (x *CreateBoardRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_boards_boards_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateBoardRequest.ProtoReflect.Descriptor instead.
func (*CreateBoardRequest) Descriptor() ([]byte, []int) {
	return file_api_boards_boards_proto_rawDescGZIP(), []int{0}
}

func (x *CreateBoardRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *CreateBoardRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

type CreateBoardResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateBoardResponse) Reset() {
	*x = CreateBoardResponse{}
	mi := &file_api_boards_boards_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateBoardResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateBoardResponse) ProtoMessage() {}

func (x *CreateBoardResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_boards_boards_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateBoardResponse.ProtoReflect.Descriptor instead.
func (*CreateBoardResponse) Descriptor() ([]byte, []int) {
	return file_api_boards_boards_proto_rawDescGZIP(), []int{1}
}

func (x *CreateBoardResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetBoardRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetBoardRequest) Reset() {
	*x = GetBoardRequest{}
	mi := &file_api_boards_boards_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetBoardRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBoardRequest) ProtoMessage() {}

func (x *GetBoardRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_boards_boards_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBoardRequest.ProtoReflect.Descriptor instead.
func (*GetBoardRequest) Descriptor() ([]byte, []int) {
	return file_api_boards_boards_proto_rawDescGZIP(), []int{2}
}

func (x *GetBoardRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetBoardResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Tasks         []*types.Task          `protobuf:"bytes,3,rep,name=tasks,proto3" json:"tasks,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetBoardResponse) Reset() {
	*x = GetBoardResponse{}
	mi := &file_api_boards_boards_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetBoardResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBoardResponse) ProtoMessage() {}

func (x *GetBoardResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_boards_boards_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBoardResponse.ProtoReflect.Descriptor instead.
func (*GetBoardResponse) Descriptor() ([]byte, []int) {
	return file_api_boards_boards_proto_rawDescGZIP(), []int{3}
}

func (x *GetBoardResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GetBoardResponse) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *GetBoardResponse) GetTasks() []*types.Task {
	if x != nil {
		return x.Tasks
	}
	return nil
}

type GetAllBoardsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAllBoardsRequest) Reset() {
	*x = GetAllBoardsRequest{}
	mi := &file_api_boards_boards_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAllBoardsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllBoardsRequest) ProtoMessage() {}

func (x *GetAllBoardsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_boards_boards_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllBoardsRequest.ProtoReflect.Descriptor instead.
func (*GetAllBoardsRequest) Descriptor() ([]byte, []int) {
	return file_api_boards_boards_proto_rawDescGZIP(), []int{4}
}

func (x *GetAllBoardsRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetAllBoardsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Boards        []*types.Board         `protobuf:"bytes,1,rep,name=boards,proto3" json:"boards,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAllBoardsResponse) Reset() {
	*x = GetAllBoardsResponse{}
	mi := &file_api_boards_boards_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAllBoardsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllBoardsResponse) ProtoMessage() {}

func (x *GetAllBoardsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_boards_boards_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllBoardsResponse.ProtoReflect.Descriptor instead.
func (*GetAllBoardsResponse) Descriptor() ([]byte, []int) {
	return file_api_boards_boards_proto_rawDescGZIP(), []int{5}
}

func (x *GetAllBoardsResponse) GetBoards() []*types.Board {
	if x != nil {
		return x.Boards
	}
	return nil
}

var File_api_boards_boards_proto protoreflect.FileDescriptor

var file_api_boards_boards_proto_rawDesc = string([]byte{
	0x0a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x2f, 0x62, 0x6f, 0x61,
	0x72, 0x64, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x62, 0x6f, 0x61, 0x72, 0x64,
	0x73, 0x1a, 0x15, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x74, 0x79, 0x70,
	0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x49, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x49, 0x64, 0x22, 0x25, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x61,
	0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x21, 0x0a, 0x0f, 0x47, 0x65,
	0x74, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x5b, 0x0a,
	0x10, 0x47, 0x65, 0x74, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x21, 0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x54,
	0x61, 0x73, 0x6b, 0x52, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x22, 0x25, 0x0a, 0x13, 0x47, 0x65,
	0x74, 0x41, 0x6c, 0x6c, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x3c, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x42, 0x6f, 0x61, 0x72, 0x64,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x06, 0x62, 0x6f, 0x61,
	0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x06, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x32,
	0xda, 0x01, 0x0a, 0x06, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x12, 0x46, 0x0a, 0x0b, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x12, 0x1a, 0x2e, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x3d, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x12, 0x17,
	0x2e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x6f, 0x61, 0x72, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x73,
	0x2e, 0x47, 0x65, 0x74, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x49, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x42, 0x6f, 0x61, 0x72, 0x64,
	0x73, 0x12, 0x1b, 0x2e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c,
	0x6c, 0x42, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c,
	0x2e, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x42, 0x6f,
	0x61, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x10, 0x5a, 0x0e,
	0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_api_boards_boards_proto_rawDescOnce sync.Once
	file_api_boards_boards_proto_rawDescData []byte
)

func file_api_boards_boards_proto_rawDescGZIP() []byte {
	file_api_boards_boards_proto_rawDescOnce.Do(func() {
		file_api_boards_boards_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_boards_boards_proto_rawDesc), len(file_api_boards_boards_proto_rawDesc)))
	})
	return file_api_boards_boards_proto_rawDescData
}

var file_api_boards_boards_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_boards_boards_proto_goTypes = []any{
	(*CreateBoardRequest)(nil),   // 0: boards.CreateBoardRequest
	(*CreateBoardResponse)(nil),  // 1: boards.CreateBoardResponse
	(*GetBoardRequest)(nil),      // 2: boards.GetBoardRequest
	(*GetBoardResponse)(nil),     // 3: boards.GetBoardResponse
	(*GetAllBoardsRequest)(nil),  // 4: boards.GetAllBoardsRequest
	(*GetAllBoardsResponse)(nil), // 5: boards.GetAllBoardsResponse
	(*types.Task)(nil),           // 6: types.Task
	(*types.Board)(nil),          // 7: types.Board
}
var file_api_boards_boards_proto_depIdxs = []int32{
	6, // 0: boards.GetBoardResponse.tasks:type_name -> types.Task
	7, // 1: boards.GetAllBoardsResponse.boards:type_name -> types.Board
	0, // 2: boards.Boards.CreateBoard:input_type -> boards.CreateBoardRequest
	2, // 3: boards.Boards.GetBoard:input_type -> boards.GetBoardRequest
	4, // 4: boards.Boards.GetAllBoards:input_type -> boards.GetAllBoardsRequest
	1, // 5: boards.Boards.CreateBoard:output_type -> boards.CreateBoardResponse
	3, // 6: boards.Boards.GetBoard:output_type -> boards.GetBoardResponse
	5, // 7: boards.Boards.GetAllBoards:output_type -> boards.GetAllBoardsResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_boards_boards_proto_init() }
func file_api_boards_boards_proto_init() {
	if File_api_boards_boards_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_boards_boards_proto_rawDesc), len(file_api_boards_boards_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_boards_boards_proto_goTypes,
		DependencyIndexes: file_api_boards_boards_proto_depIdxs,
		MessageInfos:      file_api_boards_boards_proto_msgTypes,
	}.Build()
	File_api_boards_boards_proto = out.File
	file_api_boards_boards_proto_goTypes = nil
	file_api_boards_boards_proto_depIdxs = nil
}
