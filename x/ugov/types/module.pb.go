// Code generated manually to mirror protoc-gen-go output for the ugov module config.
// This file defines the Module config type required for app wiring.
package types

import (
	_ "cosmossdk.io/depinject/appconfig/v1alpha1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	"reflect"
	"sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Module is the config object for the ugov module.
type Module struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Authority     string                 `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Module) Reset() {
	*x = Module{}
	mi := &file_uagd_ugov_module_v1_module_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Module) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Module) ProtoMessage() {}

func (x *Module) ProtoReflect() protoreflect.Message {
	mi := &file_uagd_ugov_module_v1_module_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Module.ProtoReflect.Descriptor instead.
func (*Module) Descriptor() ([]byte, []int) {
	return file_uagd_ugov_module_v1_module_proto_rawDescGZIP(), []int{0}
}

func (x *Module) GetAuthority() string {
	if x != nil {
		return x.Authority
	}
	return ""
}

var File_uagd_ugov_module_v1_module_proto protoreflect.FileDescriptor

var file_uagd_ugov_module_v1_module_proto_rawDesc = func() []byte {
	fd := &descriptorpb.FileDescriptorProto{
		Syntax:  proto.String("proto3"),
		Name:    proto.String("uagd/ugov/module/v1/module.proto"),
		Package: proto.String("uagd.ugov.module.v1"),
		Dependency: []string{
			"cosmos/app/v1alpha1/module.proto",
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("Module"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("authority"),
						JsonName: proto.String("authority"),
						Number:   proto.Int32(1),
						Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
					},
				},
			},
		},
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("uagd/x/ugov/types"),
		},
	}

	b, err := proto.Marshal(fd)
	if err != nil {
		panic(err)
	}
	return b
}()

var (
	file_uagd_ugov_module_v1_module_proto_rawDescOnce sync.Once
	file_uagd_ugov_module_v1_module_proto_rawDescData []byte
)

func file_uagd_ugov_module_v1_module_proto_rawDescGZIP() []byte {
	file_uagd_ugov_module_v1_module_proto_rawDescOnce.Do(func() {
		file_uagd_ugov_module_v1_module_proto_rawDescData = protoimpl.X.CompressGZIP(file_uagd_ugov_module_v1_module_proto_rawDesc)
	})
	return file_uagd_ugov_module_v1_module_proto_rawDescData
}

var file_uagd_ugov_module_v1_module_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_uagd_ugov_module_v1_module_proto_goTypes = []any{
	(*Module)(nil), // 0: uagd.ugov.module.v1.Module
}
var file_uagd_ugov_module_v1_module_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_uagd_ugov_module_v1_module_proto_init() }
func file_uagd_ugov_module_v1_module_proto_init() {
	if File_uagd_ugov_module_v1_module_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_uagd_ugov_module_v1_module_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_uagd_ugov_module_v1_module_proto_goTypes,
		DependencyIndexes: file_uagd_ugov_module_v1_module_proto_depIdxs,
		MessageInfos:      file_uagd_ugov_module_v1_module_proto_msgTypes,
	}.Build()

	File_uagd_ugov_module_v1_module_proto = out.File
	file_uagd_ugov_module_v1_module_proto_rawDesc = nil
	file_uagd_ugov_module_v1_module_proto_goTypes = nil
	file_uagd_ugov_module_v1_module_proto_depIdxs = nil
}
