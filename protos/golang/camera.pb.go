// Code generated by protoc-gen-go. DO NOT EDIT.
// source: camera.proto

package lumas

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_struct "github.com/golang/protobuf/ptypes/struct"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

//It's unlikely we'll ever want to pass arguments to list
// but this message is hear just in case
type ListRequest struct {
	CameraID             []*CameraID `protobuf:"bytes,1,rep,name=cameraID,proto3" json:"cameraID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ListRequest) Reset()         { *m = ListRequest{} }
func (m *ListRequest) String() string { return proto.CompactTextString(m) }
func (*ListRequest) ProtoMessage()    {}
func (*ListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f444de3b23b55d1, []int{0}
}

func (m *ListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRequest.Unmarshal(m, b)
}
func (m *ListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRequest.Marshal(b, m, deterministic)
}
func (m *ListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRequest.Merge(m, src)
}
func (m *ListRequest) XXX_Size() int {
	return xxx_messageInfo_ListRequest.Size(m)
}
func (m *ListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListRequest proto.InternalMessageInfo

func (m *ListRequest) GetCameraID() []*CameraID {
	if m != nil {
		return m.CameraID
	}
	return nil
}

type CameraID struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CameraID) Reset()         { *m = CameraID{} }
func (m *CameraID) String() string { return proto.CompactTextString(m) }
func (*CameraID) ProtoMessage()    {}
func (*CameraID) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f444de3b23b55d1, []int{1}
}

func (m *CameraID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CameraID.Unmarshal(m, b)
}
func (m *CameraID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CameraID.Marshal(b, m, deterministic)
}
func (m *CameraID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CameraID.Merge(m, src)
}
func (m *CameraID) XXX_Size() int {
	return xxx_messageInfo_CameraID.Size(m)
}
func (m *CameraID) XXX_DiscardUnknown() {
	xxx_messageInfo_CameraID.DiscardUnknown(m)
}

var xxx_messageInfo_CameraID proto.InternalMessageInfo

func (m *CameraID) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type Result struct {
	Successful           bool     `protobuf:"varint,1,opt,name=successful,proto3" json:"successful,omitempty"`
	ErrorKind            string   `protobuf:"bytes,2,opt,name=errorKind,proto3" json:"errorKind,omitempty"`
	Message              string   `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Result) Reset()         { *m = Result{} }
func (m *Result) String() string { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()    {}
func (*Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f444de3b23b55d1, []int{2}
}

func (m *Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Result.Unmarshal(m, b)
}
func (m *Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Result.Marshal(b, m, deterministic)
}
func (m *Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Result.Merge(m, src)
}
func (m *Result) XXX_Size() int {
	return xxx_messageInfo_Result.Size(m)
}
func (m *Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Result proto.InternalMessageInfo

func (m *Result) GetSuccessful() bool {
	if m != nil {
		return m.Successful
	}
	return false
}

func (m *Result) GetErrorKind() string {
	if m != nil {
		return m.ErrorKind
	}
	return ""
}

func (m *Result) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type CameraConfig struct {
	Id                   int64           `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string          `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Provider             string          `protobuf:"bytes,3,opt,name=provider,proto3" json:"provider,omitempty"`
	ProviderAddress      string          `protobuf:"bytes,4,opt,name=providerAddress,proto3" json:"providerAddress,omitempty"`
	ProviderConfig       *_struct.Struct `protobuf:"bytes,5,opt,name=providerConfig,proto3" json:"providerConfig,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *CameraConfig) Reset()         { *m = CameraConfig{} }
func (m *CameraConfig) String() string { return proto.CompactTextString(m) }
func (*CameraConfig) ProtoMessage()    {}
func (*CameraConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_2f444de3b23b55d1, []int{3}
}

func (m *CameraConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CameraConfig.Unmarshal(m, b)
}
func (m *CameraConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CameraConfig.Marshal(b, m, deterministic)
}
func (m *CameraConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CameraConfig.Merge(m, src)
}
func (m *CameraConfig) XXX_Size() int {
	return xxx_messageInfo_CameraConfig.Size(m)
}
func (m *CameraConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_CameraConfig.DiscardUnknown(m)
}

var xxx_messageInfo_CameraConfig proto.InternalMessageInfo

func (m *CameraConfig) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CameraConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CameraConfig) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

func (m *CameraConfig) GetProviderAddress() string {
	if m != nil {
		return m.ProviderAddress
	}
	return ""
}

func (m *CameraConfig) GetProviderConfig() *_struct.Struct {
	if m != nil {
		return m.ProviderConfig
	}
	return nil
}

func init() {
	proto.RegisterType((*ListRequest)(nil), "lumas.ListRequest")
	proto.RegisterType((*CameraID)(nil), "lumas.CameraID")
	proto.RegisterType((*Result)(nil), "lumas.Result")
	proto.RegisterType((*CameraConfig)(nil), "lumas.CameraConfig")
}

func init() { proto.RegisterFile("camera.proto", fileDescriptor_2f444de3b23b55d1) }

var fileDescriptor_2f444de3b23b55d1 = []byte{
	// 350 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0xc1, 0x6a, 0xab, 0x40,
	0x14, 0x86, 0x99, 0x68, 0x8c, 0x39, 0xc9, 0x4d, 0xe0, 0xdc, 0xc5, 0x15, 0x09, 0x17, 0x71, 0x51,
	0x0c, 0x05, 0xd3, 0xa6, 0xbb, 0x6e, 0x4a, 0x49, 0x37, 0xa5, 0x5d, 0xd9, 0x17, 0xa8, 0x71, 0x26,
	0x22, 0x68, 0xc6, 0xce, 0x8c, 0x79, 0xb7, 0x3e, 0x4f, 0x5f, 0xa4, 0x74, 0x46, 0xd3, 0x54, 0x0a,
	0xd9, 0xcd, 0x7c, 0x7e, 0xfe, 0x73, 0xf8, 0x0f, 0x4c, 0xb3, 0xb4, 0x62, 0x22, 0x8d, 0x6b, 0xc1,
	0x15, 0xc7, 0x61, 0xd9, 0x54, 0xa9, 0xf4, 0x17, 0x39, 0xe7, 0x79, 0xc9, 0x56, 0x1a, 0x6e, 0x9b,
	0xdd, 0x4a, 0x2a, 0xd1, 0x64, 0xca, 0x48, 0xe1, 0x2d, 0x4c, 0x9e, 0x0b, 0xa9, 0x12, 0xf6, 0xd6,
	0x30, 0xa9, 0xf0, 0x12, 0x5c, 0x93, 0xf1, 0xf8, 0xe0, 0x91, 0xc0, 0x8a, 0x26, 0xeb, 0x79, 0xac,
	0x63, 0xe2, 0x4d, 0x8b, 0x93, 0xa3, 0x10, 0xfa, 0xe0, 0x76, 0x14, 0x67, 0x30, 0x28, 0xa8, 0x47,
	0x02, 0x12, 0x59, 0xc9, 0xa0, 0xa0, 0xe1, 0x2b, 0x38, 0x09, 0x93, 0x4d, 0xa9, 0xf0, 0x3f, 0x80,
	0x6c, 0xb2, 0x8c, 0x49, 0xb9, 0x6b, 0x4a, 0x6d, 0xb8, 0xc9, 0x09, 0xc1, 0x05, 0x8c, 0x99, 0x10,
	0x5c, 0x3c, 0x15, 0x7b, 0xea, 0x0d, 0x02, 0x12, 0x8d, 0x93, 0x6f, 0x80, 0x1e, 0x8c, 0x2a, 0x26,
	0x65, 0x9a, 0x33, 0xcf, 0xd2, 0xdf, 0xba, 0x6b, 0xf8, 0x4e, 0x60, 0x6a, 0x9e, 0xdf, 0xf0, 0xfd,
	0xae, 0xc8, 0xfb, 0x23, 0x20, 0x82, 0xbd, 0x4f, 0x2b, 0xd6, 0x66, 0xea, 0x33, 0xfa, 0xe0, 0xd6,
	0x82, 0x1f, 0x0a, 0xca, 0x44, 0x9b, 0x77, 0xbc, 0x63, 0x04, 0xf3, 0xee, 0x7c, 0x4f, 0xa9, 0x60,
	0x52, 0x7a, 0xb6, 0x56, 0xfa, 0x18, 0xef, 0x60, 0xd6, 0x21, 0xf3, 0xb6, 0x37, 0x0c, 0x48, 0x34,
	0x59, 0xff, 0x8b, 0x4d, 0xd7, 0x71, 0xd7, 0x75, 0xfc, 0xa2, 0xbb, 0x4e, 0x7a, 0xfa, 0xfa, 0x83,
	0x80, 0x63, 0x66, 0xc7, 0x25, 0x58, 0x29, 0xa5, 0xf8, 0xf7, 0x47, 0xcd, 0x46, 0xf3, 0xff, 0xb4,
	0xb0, 0x6d, 0xf2, 0x1a, 0xec, 0xb2, 0x90, 0x0a, 0xb1, 0xc5, 0x27, 0x8b, 0xf3, 0x7f, 0xfb, 0xff,
	0x8a, 0x60, 0x04, 0x8e, 0x60, 0x15, 0x3f, 0x30, 0xec, 0xef, 0xb1, 0x1f, 0x7e, 0x01, 0xb6, 0x54,
	0xbc, 0x3e, 0xeb, 0x2d, 0x61, 0x54, 0x0b, 0xfe, 0xb5, 0xbc, 0x73, 0xea, 0xd6, 0xd1, 0x35, 0xdc,
	0x7c, 0x06, 0x00, 0x00, 0xff, 0xff, 0x2f, 0x59, 0x8f, 0x12, 0x97, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CameraClient is the client API for Camera service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CameraClient interface {
	Add(ctx context.Context, in *CameraConfig, opts ...grpc.CallOption) (*Result, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (Camera_ListClient, error)
	Remove(ctx context.Context, in *CameraID, opts ...grpc.CallOption) (*Result, error)
	Stop(ctx context.Context, in *CameraID, opts ...grpc.CallOption) (*Result, error)
	Process(ctx context.Context, in *CameraID, opts ...grpc.CallOption) (*Result, error)
}

type cameraClient struct {
	cc *grpc.ClientConn
}

func NewCameraClient(cc *grpc.ClientConn) CameraClient {
	return &cameraClient{cc}
}

func (c *cameraClient) Add(ctx context.Context, in *CameraConfig, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/lumas.Camera/add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cameraClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (Camera_ListClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Camera_serviceDesc.Streams[0], "/lumas.Camera/list", opts...)
	if err != nil {
		return nil, err
	}
	x := &cameraListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Camera_ListClient interface {
	Recv() (*CameraConfig, error)
	grpc.ClientStream
}

type cameraListClient struct {
	grpc.ClientStream
}

func (x *cameraListClient) Recv() (*CameraConfig, error) {
	m := new(CameraConfig)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *cameraClient) Remove(ctx context.Context, in *CameraID, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/lumas.Camera/remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cameraClient) Stop(ctx context.Context, in *CameraID, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/lumas.Camera/stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cameraClient) Process(ctx context.Context, in *CameraID, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/lumas.Camera/process", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CameraServer is the server API for Camera service.
type CameraServer interface {
	Add(context.Context, *CameraConfig) (*Result, error)
	List(*ListRequest, Camera_ListServer) error
	Remove(context.Context, *CameraID) (*Result, error)
	Stop(context.Context, *CameraID) (*Result, error)
	Process(context.Context, *CameraID) (*Result, error)
}

func RegisterCameraServer(s *grpc.Server, srv CameraServer) {
	s.RegisterService(&_Camera_serviceDesc, srv)
}

func _Camera_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CameraConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CameraServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lumas.Camera/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CameraServer).Add(ctx, req.(*CameraConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _Camera_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CameraServer).List(m, &cameraListServer{stream})
}

type Camera_ListServer interface {
	Send(*CameraConfig) error
	grpc.ServerStream
}

type cameraListServer struct {
	grpc.ServerStream
}

func (x *cameraListServer) Send(m *CameraConfig) error {
	return x.ServerStream.SendMsg(m)
}

func _Camera_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CameraID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CameraServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lumas.Camera/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CameraServer).Remove(ctx, req.(*CameraID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Camera_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CameraID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CameraServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lumas.Camera/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CameraServer).Stop(ctx, req.(*CameraID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Camera_Process_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CameraID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CameraServer).Process(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lumas.Camera/Process",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CameraServer).Process(ctx, req.(*CameraID))
	}
	return interceptor(ctx, in, info, handler)
}

var _Camera_serviceDesc = grpc.ServiceDesc{
	ServiceName: "lumas.Camera",
	HandlerType: (*CameraServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "add",
			Handler:    _Camera_Add_Handler,
		},
		{
			MethodName: "remove",
			Handler:    _Camera_Remove_Handler,
		},
		{
			MethodName: "stop",
			Handler:    _Camera_Stop_Handler,
		},
		{
			MethodName: "process",
			Handler:    _Camera_Process_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "list",
			Handler:       _Camera_List_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "camera.proto",
}
