// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: dispatcher.proto

package dispatcher

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

type Service struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Service) Reset() {
	*x = Service{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dispatcher_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Service) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Service) ProtoMessage() {}

func (x *Service) ProtoReflect() protoreflect.Message {
	mi := &file_dispatcher_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Service.ProtoReflect.Descriptor instead.
func (*Service) Descriptor() ([]byte, []int) {
	return file_dispatcher_proto_rawDescGZIP(), []int{0}
}

func (x *Service) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type ServiceStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Svc        *Service               `protobuf:"bytes,1,opt,name=svc,proto3" json:"svc,omitempty"`
	Running    bool                   `protobuf:"varint,2,opt,name=running,proto3" json:"running,omitempty"`
	Pid        uint32                 `protobuf:"varint,3,opt,name=pid,proto3" json:"pid,omitempty"`
	ExitStatus uint32                 `protobuf:"varint,4,opt,name=exit_status,json=exitStatus,proto3" json:"exit_status,omitempty"`
	StartTime  *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime    *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	Success    bool                   `protobuf:"varint,7,opt,name=success,proto3" json:"success,omitempty"`
	Error      string                 `protobuf:"bytes,8,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *ServiceStatus) Reset() {
	*x = ServiceStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dispatcher_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceStatus) ProtoMessage() {}

func (x *ServiceStatus) ProtoReflect() protoreflect.Message {
	mi := &file_dispatcher_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceStatus.ProtoReflect.Descriptor instead.
func (*ServiceStatus) Descriptor() ([]byte, []int) {
	return file_dispatcher_proto_rawDescGZIP(), []int{1}
}

func (x *ServiceStatus) GetSvc() *Service {
	if x != nil {
		return x.Svc
	}
	return nil
}

func (x *ServiceStatus) GetRunning() bool {
	if x != nil {
		return x.Running
	}
	return false
}

func (x *ServiceStatus) GetPid() uint32 {
	if x != nil {
		return x.Pid
	}
	return 0
}

func (x *ServiceStatus) GetExitStatus() uint32 {
	if x != nil {
		return x.ExitStatus
	}
	return 0
}

func (x *ServiceStatus) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *ServiceStatus) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

func (x *ServiceStatus) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *ServiceStatus) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_dispatcher_proto protoreflect.FileDescriptor

var file_dispatcher_proto_rawDesc = []byte{
	0x0a, 0x10, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x1d, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0x9a, 0x02, 0x0a, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x1a, 0x0a, 0x03, 0x73, 0x76, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08,
	0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x03, 0x73, 0x76, 0x63, 0x12, 0x18, 0x0a,
	0x07, 0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x72, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x78, 0x69,
	0x74, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a,
	0x65, 0x78, 0x69, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x80, 0x02, 0x0a,
	0x0a, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x12, 0x2b, 0x0a, 0x05, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x12, 0x08, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x2a, 0x0a, 0x04, 0x53, 0x74, 0x6f, 0x70,
	0x12, 0x08, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x00, 0x12, 0x24, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x08,
	0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x0e, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x06, 0x52, 0x65,
	0x6c, 0x6f, 0x61, 0x64, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x41, 0x6c, 0x6c, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0e, 0x2e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x00, 0x30, 0x01, 0x42,
	0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76, 0x69,
	0x6e, 0x79, 0x6c, 0x2d, 0x6c, 0x69, 0x6e, 0x75, 0x78, 0x2f, 0x76, 0x69, 0x6e, 0x69, 0x74, 0x2f,
	0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_dispatcher_proto_rawDescOnce sync.Once
	file_dispatcher_proto_rawDescData = file_dispatcher_proto_rawDesc
)

func file_dispatcher_proto_rawDescGZIP() []byte {
	file_dispatcher_proto_rawDescOnce.Do(func() {
		file_dispatcher_proto_rawDescData = protoimpl.X.CompressGZIP(file_dispatcher_proto_rawDescData)
	})
	return file_dispatcher_proto_rawDescData
}

var file_dispatcher_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_dispatcher_proto_goTypes = []interface{}{
	(*Service)(nil),               // 0: Service
	(*ServiceStatus)(nil),         // 1: ServiceStatus
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 3: google.protobuf.Empty
}
var file_dispatcher_proto_depIdxs = []int32{
	0, // 0: ServiceStatus.svc:type_name -> Service
	2, // 1: ServiceStatus.start_time:type_name -> google.protobuf.Timestamp
	2, // 2: ServiceStatus.end_time:type_name -> google.protobuf.Timestamp
	0, // 3: Dispatcher.Start:input_type -> Service
	0, // 4: Dispatcher.Stop:input_type -> Service
	0, // 5: Dispatcher.Status:input_type -> Service
	3, // 6: Dispatcher.Reload:input_type -> google.protobuf.Empty
	3, // 7: Dispatcher.StatusAll:input_type -> google.protobuf.Empty
	3, // 8: Dispatcher.Start:output_type -> google.protobuf.Empty
	3, // 9: Dispatcher.Stop:output_type -> google.protobuf.Empty
	1, // 10: Dispatcher.Status:output_type -> ServiceStatus
	3, // 11: Dispatcher.Reload:output_type -> google.protobuf.Empty
	1, // 12: Dispatcher.StatusAll:output_type -> ServiceStatus
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_dispatcher_proto_init() }
func file_dispatcher_proto_init() {
	if File_dispatcher_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dispatcher_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Service); i {
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
		file_dispatcher_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceStatus); i {
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
			RawDescriptor: file_dispatcher_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_dispatcher_proto_goTypes,
		DependencyIndexes: file_dispatcher_proto_depIdxs,
		MessageInfos:      file_dispatcher_proto_msgTypes,
	}.Build()
	File_dispatcher_proto = out.File
	file_dispatcher_proto_rawDesc = nil
	file_dispatcher_proto_goTypes = nil
	file_dispatcher_proto_depIdxs = nil
}
