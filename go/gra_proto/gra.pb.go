// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: gra.proto

package gra_proto

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

type WizytowkaGracza struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nazwa string `protobuf:"bytes,1,opt,name=nazwa,proto3" json:"nazwa,omitempty"`
}

func (x *WizytowkaGracza) Reset() {
	*x = WizytowkaGracza{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gra_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WizytowkaGracza) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WizytowkaGracza) ProtoMessage() {}

func (x *WizytowkaGracza) ProtoReflect() protoreflect.Message {
	mi := &file_gra_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WizytowkaGracza.ProtoReflect.Descriptor instead.
func (*WizytowkaGracza) Descriptor() ([]byte, []int) {
	return file_gra_proto_rawDescGZIP(), []int{0}
}

func (x *WizytowkaGracza) GetNazwa() string {
	if x != nil {
		return x.Nazwa
	}
	return ""
}

type StanGry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IdGry             string `protobuf:"bytes,1,opt,name=idGry,proto3" json:"idGry,omitempty"`
	SytuacjaNaPlanszy string `protobuf:"bytes,2,opt,name=sytuacjaNaPlanszy,proto3" json:"sytuacjaNaPlanszy,omitempty"`
	TwojeKarty        string `protobuf:"bytes,3,opt,name=twojeKarty,proto3" json:"twojeKarty,omitempty"`
}

func (x *StanGry) Reset() {
	*x = StanGry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gra_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StanGry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StanGry) ProtoMessage() {}

func (x *StanGry) ProtoReflect() protoreflect.Message {
	mi := &file_gra_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StanGry.ProtoReflect.Descriptor instead.
func (*StanGry) Descriptor() ([]byte, []int) {
	return file_gra_proto_rawDescGZIP(), []int{1}
}

func (x *StanGry) GetIdGry() string {
	if x != nil {
		return x.IdGry
	}
	return ""
}

func (x *StanGry) GetSytuacjaNaPlanszy() string {
	if x != nil {
		return x.SytuacjaNaPlanszy
	}
	return ""
}

func (x *StanGry) GetTwojeKarty() string {
	if x != nil {
		return x.TwojeKarty
	}
	return ""
}

type RuchGracza struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IdGry        string `protobuf:"bytes,1,opt,name=idGry,proto3" json:"idGry,omitempty"`
	ZagranaKarta string `protobuf:"bytes,2,opt,name=zagranaKarta,proto3" json:"zagranaKarta,omitempty"`
}

func (x *RuchGracza) Reset() {
	*x = RuchGracza{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gra_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RuchGracza) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RuchGracza) ProtoMessage() {}

func (x *RuchGracza) ProtoReflect() protoreflect.Message {
	mi := &file_gra_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RuchGracza.ProtoReflect.Descriptor instead.
func (*RuchGracza) Descriptor() ([]byte, []int) {
	return file_gra_proto_rawDescGZIP(), []int{2}
}

func (x *RuchGracza) GetIdGry() string {
	if x != nil {
		return x.IdGry
	}
	return ""
}

func (x *RuchGracza) GetZagranaKarta() string {
	if x != nil {
		return x.ZagranaKarta
	}
	return ""
}

var File_gra_proto protoreflect.FileDescriptor

var file_gra_proto_rawDesc = []byte{
	0x0a, 0x09, 0x67, 0x72, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x67, 0x72, 0x61,
	0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x27, 0x0a, 0x0f, 0x57, 0x69, 0x7a, 0x79, 0x74, 0x6f,
	0x77, 0x6b, 0x61, 0x47, 0x72, 0x61, 0x63, 0x7a, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x61, 0x7a,
	0x77, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6e, 0x61, 0x7a, 0x77, 0x61, 0x22,
	0x6d, 0x0a, 0x07, 0x53, 0x74, 0x61, 0x6e, 0x47, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x64,
	0x47, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x64, 0x47, 0x72, 0x79,
	0x12, 0x2c, 0x0a, 0x11, 0x73, 0x79, 0x74, 0x75, 0x61, 0x63, 0x6a, 0x61, 0x4e, 0x61, 0x50, 0x6c,
	0x61, 0x6e, 0x73, 0x7a, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x73, 0x79, 0x74,
	0x75, 0x61, 0x63, 0x6a, 0x61, 0x4e, 0x61, 0x50, 0x6c, 0x61, 0x6e, 0x73, 0x7a, 0x79, 0x12, 0x1e,
	0x0a, 0x0a, 0x74, 0x77, 0x6f, 0x6a, 0x65, 0x4b, 0x61, 0x72, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x74, 0x77, 0x6f, 0x6a, 0x65, 0x4b, 0x61, 0x72, 0x74, 0x79, 0x22, 0x46,
	0x0a, 0x0a, 0x52, 0x75, 0x63, 0x68, 0x47, 0x72, 0x61, 0x63, 0x7a, 0x61, 0x12, 0x14, 0x0a, 0x05,
	0x69, 0x64, 0x47, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x64, 0x47,
	0x72, 0x79, 0x12, 0x22, 0x0a, 0x0c, 0x7a, 0x61, 0x67, 0x72, 0x61, 0x6e, 0x61, 0x4b, 0x61, 0x72,
	0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x7a, 0x61, 0x67, 0x72, 0x61, 0x6e,
	0x61, 0x4b, 0x61, 0x72, 0x74, 0x61, 0x32, 0x7b, 0x0a, 0x03, 0x47, 0x72, 0x61, 0x12, 0x3c, 0x0a,
	0x08, 0x4e, 0x6f, 0x77, 0x79, 0x4d, 0x65, 0x63, 0x7a, 0x12, 0x1a, 0x2e, 0x67, 0x72, 0x61, 0x5f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57, 0x69, 0x7a, 0x79, 0x74, 0x6f, 0x77, 0x6b, 0x61, 0x47,
	0x72, 0x61, 0x63, 0x7a, 0x61, 0x1a, 0x12, 0x2e, 0x67, 0x72, 0x61, 0x5f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x53, 0x74, 0x61, 0x6e, 0x47, 0x72, 0x79, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x07, 0x4d,
	0x6f, 0x6a, 0x52, 0x75, 0x63, 0x68, 0x12, 0x15, 0x2e, 0x67, 0x72, 0x61, 0x5f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x52, 0x75, 0x63, 0x68, 0x47, 0x72, 0x61, 0x63, 0x7a, 0x61, 0x1a, 0x12, 0x2e,
	0x67, 0x72, 0x61, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x6e, 0x47, 0x72,
	0x79, 0x22, 0x00, 0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x6c, 0x61, 0x72, 0x61, 0x7a, 0x2f, 0x74, 0x75, 0x72, 0x6e, 0x69, 0x65, 0x6a,
	0x2f, 0x67, 0x6f, 0x2f, 0x67, 0x72, 0x61, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gra_proto_rawDescOnce sync.Once
	file_gra_proto_rawDescData = file_gra_proto_rawDesc
)

func file_gra_proto_rawDescGZIP() []byte {
	file_gra_proto_rawDescOnce.Do(func() {
		file_gra_proto_rawDescData = protoimpl.X.CompressGZIP(file_gra_proto_rawDescData)
	})
	return file_gra_proto_rawDescData
}

var file_gra_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_gra_proto_goTypes = []interface{}{
	(*WizytowkaGracza)(nil), // 0: gra_proto.WizytowkaGracza
	(*StanGry)(nil),         // 1: gra_proto.StanGry
	(*RuchGracza)(nil),      // 2: gra_proto.RuchGracza
}
var file_gra_proto_depIdxs = []int32{
	0, // 0: gra_proto.Gra.NowyMecz:input_type -> gra_proto.WizytowkaGracza
	2, // 1: gra_proto.Gra.MojRuch:input_type -> gra_proto.RuchGracza
	1, // 2: gra_proto.Gra.NowyMecz:output_type -> gra_proto.StanGry
	1, // 3: gra_proto.Gra.MojRuch:output_type -> gra_proto.StanGry
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gra_proto_init() }
func file_gra_proto_init() {
	if File_gra_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gra_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WizytowkaGracza); i {
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
		file_gra_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StanGry); i {
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
		file_gra_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RuchGracza); i {
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
			RawDescriptor: file_gra_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gra_proto_goTypes,
		DependencyIndexes: file_gra_proto_depIdxs,
		MessageInfos:      file_gra_proto_msgTypes,
	}.Build()
	File_gra_proto = out.File
	file_gra_proto_rawDesc = nil
	file_gra_proto_goTypes = nil
	file_gra_proto_depIdxs = nil
}