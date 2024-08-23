// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.12.4
// source: room/service/v1/room_error.proto

package v1

import (
	_ "github.com/go-kratos/kratos/v2/errors"
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

type UserServiceErrorReason int32

const (
	UserServiceErrorReason_UNKNOWN_ERROR              UserServiceErrorReason = 0
	UserServiceErrorReason_ROOM_NOT_FOUND             UserServiceErrorReason = 2
	UserServiceErrorReason_ROOM_NOT_ACCESSIBLE        UserServiceErrorReason = 3
	UserServiceErrorReason_CREATE_ROOM_FAILED         UserServiceErrorReason = 4
	UserServiceErrorReason_GET_ROOM_FAILED            UserServiceErrorReason = 5
	UserServiceErrorReason_JOIN_ROOM_FAILED           UserServiceErrorReason = 6
	UserServiceErrorReason_CREATE_ROOM_SESSION_FAILED UserServiceErrorReason = 7
	UserServiceErrorReason_UPDATE_ROOM_FAILED         UserServiceErrorReason = 8
	UserServiceErrorReason_DELETE_ROOM_FAILED         UserServiceErrorReason = 9
	UserServiceErrorReason_GET_ROOM_MEMBER_FAILED     UserServiceErrorReason = 10
	UserServiceErrorReason_UPDATE_MEMBER_FAILED       UserServiceErrorReason = 11
	UserServiceErrorReason_DELETE_MEMBER_FAILED       UserServiceErrorReason = 12
)

// Enum value maps for UserServiceErrorReason.
var (
	UserServiceErrorReason_name = map[int32]string{
		0:  "UNKNOWN_ERROR",
		2:  "ROOM_NOT_FOUND",
		3:  "ROOM_NOT_ACCESSIBLE",
		4:  "CREATE_ROOM_FAILED",
		5:  "GET_ROOM_FAILED",
		6:  "JOIN_ROOM_FAILED",
		7:  "CREATE_ROOM_SESSION_FAILED",
		8:  "UPDATE_ROOM_FAILED",
		9:  "DELETE_ROOM_FAILED",
		10: "GET_ROOM_MEMBER_FAILED",
		11: "UPDATE_MEMBER_FAILED",
		12: "DELETE_MEMBER_FAILED",
	}
	UserServiceErrorReason_value = map[string]int32{
		"UNKNOWN_ERROR":              0,
		"ROOM_NOT_FOUND":             2,
		"ROOM_NOT_ACCESSIBLE":        3,
		"CREATE_ROOM_FAILED":         4,
		"GET_ROOM_FAILED":            5,
		"JOIN_ROOM_FAILED":           6,
		"CREATE_ROOM_SESSION_FAILED": 7,
		"UPDATE_ROOM_FAILED":         8,
		"DELETE_ROOM_FAILED":         9,
		"GET_ROOM_MEMBER_FAILED":     10,
		"UPDATE_MEMBER_FAILED":       11,
		"DELETE_MEMBER_FAILED":       12,
	}
)

func (x UserServiceErrorReason) Enum() *UserServiceErrorReason {
	p := new(UserServiceErrorReason)
	*p = x
	return p
}

func (x UserServiceErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserServiceErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_room_service_v1_room_error_proto_enumTypes[0].Descriptor()
}

func (UserServiceErrorReason) Type() protoreflect.EnumType {
	return &file_room_service_v1_room_error_proto_enumTypes[0]
}

func (x UserServiceErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserServiceErrorReason.Descriptor instead.
func (UserServiceErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_room_service_v1_room_error_proto_rawDescGZIP(), []int{0}
}

var File_room_service_v1_room_error_proto protoreflect.FileDescriptor

var file_room_service_v1_room_error_proto_rawDesc = []byte{
	0x0a, 0x20, 0x72, 0x6f, 0x6f, 0x6d, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76,
	0x31, 0x2f, 0x72, 0x6f, 0x6f, 0x6d, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x07, 0x72, 0x6f, 0x6f, 0x6d, 0x2e, 0x76, 0x31, 0x1a, 0x13, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2a, 0xd3, 0x02, 0x0a, 0x16, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x11, 0x0a, 0x0d, 0x55,
	0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x00, 0x12, 0x18,
	0x0a, 0x0e, 0x52, 0x4f, 0x4f, 0x4d, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44,
	0x10, 0x02, 0x1a, 0x04, 0xa8, 0x45, 0x94, 0x03, 0x12, 0x1d, 0x0a, 0x13, 0x52, 0x4f, 0x4f, 0x4d,
	0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x49, 0x42, 0x4c, 0x45, 0x10,
	0x03, 0x1a, 0x04, 0xa8, 0x45, 0x93, 0x03, 0x12, 0x1c, 0x0a, 0x12, 0x43, 0x52, 0x45, 0x41, 0x54,
	0x45, 0x5f, 0x52, 0x4f, 0x4f, 0x4d, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x04, 0x1a,
	0x04, 0xa8, 0x45, 0x90, 0x03, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x45, 0x54, 0x5f, 0x52, 0x4f, 0x4f,
	0x4d, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x05, 0x12, 0x14, 0x0a, 0x10, 0x4a, 0x4f,
	0x49, 0x4e, 0x5f, 0x52, 0x4f, 0x4f, 0x4d, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x06,
	0x12, 0x1e, 0x0a, 0x1a, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x5f, 0x52, 0x4f, 0x4f, 0x4d, 0x5f,
	0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x07,
	0x12, 0x16, 0x0a, 0x12, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x52, 0x4f, 0x4f, 0x4d, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x08, 0x12, 0x16, 0x0a, 0x12, 0x44, 0x45, 0x4c, 0x45,
	0x54, 0x45, 0x5f, 0x52, 0x4f, 0x4f, 0x4d, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x09,
	0x12, 0x1a, 0x0a, 0x16, 0x47, 0x45, 0x54, 0x5f, 0x52, 0x4f, 0x4f, 0x4d, 0x5f, 0x4d, 0x45, 0x4d,
	0x42, 0x45, 0x52, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0a, 0x12, 0x18, 0x0a, 0x14,
	0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x4d, 0x45, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x46, 0x41,
	0x49, 0x4c, 0x45, 0x44, 0x10, 0x0b, 0x12, 0x18, 0x0a, 0x14, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45,
	0x5f, 0x4d, 0x45, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0c,
	0x1a, 0x04, 0xa0, 0x45, 0xf4, 0x03, 0x42, 0x20, 0x5a, 0x1e, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e,
	0x64, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x6f, 0x6f, 0x6d, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_room_service_v1_room_error_proto_rawDescOnce sync.Once
	file_room_service_v1_room_error_proto_rawDescData = file_room_service_v1_room_error_proto_rawDesc
)

func file_room_service_v1_room_error_proto_rawDescGZIP() []byte {
	file_room_service_v1_room_error_proto_rawDescOnce.Do(func() {
		file_room_service_v1_room_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_room_service_v1_room_error_proto_rawDescData)
	})
	return file_room_service_v1_room_error_proto_rawDescData
}

var file_room_service_v1_room_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_room_service_v1_room_error_proto_goTypes = []any{
	(UserServiceErrorReason)(0), // 0: room.v1.UserServiceErrorReason
}
var file_room_service_v1_room_error_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_room_service_v1_room_error_proto_init() }
func file_room_service_v1_room_error_proto_init() {
	if File_room_service_v1_room_error_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_room_service_v1_room_error_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_room_service_v1_room_error_proto_goTypes,
		DependencyIndexes: file_room_service_v1_room_error_proto_depIdxs,
		EnumInfos:         file_room_service_v1_room_error_proto_enumTypes,
	}.Build()
	File_room_service_v1_room_error_proto = out.File
	file_room_service_v1_room_error_proto_rawDesc = nil
	file_room_service_v1_room_error_proto_goTypes = nil
	file_room_service_v1_room_error_proto_depIdxs = nil
}