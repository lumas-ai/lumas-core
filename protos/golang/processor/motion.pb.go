// Code generated by protoc-gen-go. DO NOT EDIT.
// source: motion.proto

package processor

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type MotionArea struct {
	X                    int32    `protobuf:"varint,1,opt,name=x,proto3" json:"x,omitempty"`
	Y                    int32    `protobuf:"varint,2,opt,name=y,proto3" json:"y,omitempty"`
	Height               int32    `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
	Width                int32    `protobuf:"varint,4,opt,name=width,proto3" json:"width,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MotionArea) Reset()         { *m = MotionArea{} }
func (m *MotionArea) String() string { return proto.CompactTextString(m) }
func (*MotionArea) ProtoMessage()    {}
func (*MotionArea) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea45689136bdb868, []int{0}
}

func (m *MotionArea) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MotionArea.Unmarshal(m, b)
}
func (m *MotionArea) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MotionArea.Marshal(b, m, deterministic)
}
func (m *MotionArea) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MotionArea.Merge(m, src)
}
func (m *MotionArea) XXX_Size() int {
	return xxx_messageInfo_MotionArea.Size(m)
}
func (m *MotionArea) XXX_DiscardUnknown() {
	xxx_messageInfo_MotionArea.DiscardUnknown(m)
}

var xxx_messageInfo_MotionArea proto.InternalMessageInfo

func (m *MotionArea) GetX() int32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *MotionArea) GetY() int32 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *MotionArea) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *MotionArea) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

type Images struct {
	FirstImage           *Image   `protobuf:"bytes,1,opt,name=firstImage,proto3" json:"firstImage,omitempty"`
	SecondImage          *Image   `protobuf:"bytes,2,opt,name=secondImage,proto3" json:"secondImage,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Images) Reset()         { *m = Images{} }
func (m *Images) String() string { return proto.CompactTextString(m) }
func (*Images) ProtoMessage()    {}
func (*Images) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea45689136bdb868, []int{1}
}

func (m *Images) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Images.Unmarshal(m, b)
}
func (m *Images) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Images.Marshal(b, m, deterministic)
}
func (m *Images) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Images.Merge(m, src)
}
func (m *Images) XXX_Size() int {
	return xxx_messageInfo_Images.Size(m)
}
func (m *Images) XXX_DiscardUnknown() {
	xxx_messageInfo_Images.DiscardUnknown(m)
}

var xxx_messageInfo_Images proto.InternalMessageInfo

func (m *Images) GetFirstImage() *Image {
	if m != nil {
		return m.FirstImage
	}
	return nil
}

func (m *Images) GetSecondImage() *Image {
	if m != nil {
		return m.SecondImage
	}
	return nil
}

type MotionResults struct {
	MotionDetected       bool          `protobuf:"varint,1,opt,name=motionDetected,proto3" json:"motionDetected,omitempty"`
	MotionAreas          []*MotionArea `protobuf:"bytes,2,rep,name=motionAreas,proto3" json:"motionAreas,omitempty"`
	Image                *Image        `protobuf:"bytes,3,opt,name=image,proto3" json:"image,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *MotionResults) Reset()         { *m = MotionResults{} }
func (m *MotionResults) String() string { return proto.CompactTextString(m) }
func (*MotionResults) ProtoMessage()    {}
func (*MotionResults) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea45689136bdb868, []int{2}
}

func (m *MotionResults) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MotionResults.Unmarshal(m, b)
}
func (m *MotionResults) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MotionResults.Marshal(b, m, deterministic)
}
func (m *MotionResults) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MotionResults.Merge(m, src)
}
func (m *MotionResults) XXX_Size() int {
	return xxx_messageInfo_MotionResults.Size(m)
}
func (m *MotionResults) XXX_DiscardUnknown() {
	xxx_messageInfo_MotionResults.DiscardUnknown(m)
}

var xxx_messageInfo_MotionResults proto.InternalMessageInfo

func (m *MotionResults) GetMotionDetected() bool {
	if m != nil {
		return m.MotionDetected
	}
	return false
}

func (m *MotionResults) GetMotionAreas() []*MotionArea {
	if m != nil {
		return m.MotionAreas
	}
	return nil
}

func (m *MotionResults) GetImage() *Image {
	if m != nil {
		return m.Image
	}
	return nil
}

type Image struct {
	Base64Image          string   `protobuf:"bytes,1,opt,name=base64Image,proto3" json:"base64Image,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Image) Reset()         { *m = Image{} }
func (m *Image) String() string { return proto.CompactTextString(m) }
func (*Image) ProtoMessage()    {}
func (*Image) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea45689136bdb868, []int{3}
}

func (m *Image) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Image.Unmarshal(m, b)
}
func (m *Image) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Image.Marshal(b, m, deterministic)
}
func (m *Image) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Image.Merge(m, src)
}
func (m *Image) XXX_Size() int {
	return xxx_messageInfo_Image.Size(m)
}
func (m *Image) XXX_DiscardUnknown() {
	xxx_messageInfo_Image.DiscardUnknown(m)
}

var xxx_messageInfo_Image proto.InternalMessageInfo

func (m *Image) GetBase64Image() string {
	if m != nil {
		return m.Base64Image
	}
	return ""
}

func init() {
	proto.RegisterType((*MotionArea)(nil), "processor.MotionArea")
	proto.RegisterType((*Images)(nil), "processor.Images")
	proto.RegisterType((*MotionResults)(nil), "processor.MotionResults")
	proto.RegisterType((*Image)(nil), "processor.Image")
}

func init() { proto.RegisterFile("motion.proto", fileDescriptor_ea45689136bdb868) }

var fileDescriptor_ea45689136bdb868 = []byte{
	// 300 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0x4f, 0x4b, 0xc3, 0x30,
	0x18, 0xc6, 0x4d, 0x6b, 0x8b, 0x7b, 0x5b, 0x45, 0x5f, 0x54, 0x82, 0xa7, 0xd2, 0xc3, 0xa8, 0x97,
	0x32, 0xaa, 0xe8, 0xc9, 0x83, 0x20, 0x82, 0x07, 0x2f, 0x11, 0xbc, 0x77, 0x6d, 0x5c, 0x0b, 0xb6,
	0x19, 0x49, 0xc4, 0xed, 0x53, 0x78, 0xf5, 0xe3, 0x8a, 0x49, 0xd9, 0xc2, 0xc6, 0x76, 0x7c, 0x9e,
	0xbc, 0x7f, 0x7e, 0xcf, 0x4b, 0x20, 0xee, 0x84, 0x6e, 0x45, 0x9f, 0xcf, 0xa5, 0xd0, 0x02, 0x47,
	0x73, 0x29, 0x2a, 0xae, 0x94, 0x90, 0xe9, 0x3b, 0xc0, 0xab, 0x79, 0x7a, 0x94, 0xbc, 0xc4, 0x18,
	0xc8, 0x82, 0x92, 0x84, 0x64, 0x01, 0x23, 0x8b, 0x7f, 0xb5, 0xa4, 0x9e, 0x55, 0x4b, 0xbc, 0x84,
	0xb0, 0xe1, 0xed, 0xac, 0xd1, 0xd4, 0x37, 0xd6, 0xa0, 0xf0, 0x1c, 0x82, 0xef, 0xb6, 0xd6, 0x0d,
	0x3d, 0x34, 0xb6, 0x15, 0x69, 0x0f, 0xe1, 0x4b, 0x57, 0xce, 0xb8, 0xc2, 0x09, 0xc0, 0x47, 0x2b,
	0x95, 0x36, 0xd2, 0x0c, 0x8f, 0x8a, 0xd3, 0x7c, 0x45, 0x90, 0x1b, 0x9f, 0x39, 0x35, 0x58, 0x40,
	0xa4, 0x78, 0x25, 0xfa, 0xda, 0xb6, 0x78, 0x3b, 0x5a, 0xdc, 0xa2, 0xf4, 0x97, 0xc0, 0xb1, 0x0d,
	0xc2, 0xb8, 0xfa, 0xfa, 0xd4, 0x0a, 0xc7, 0x70, 0x62, 0x43, 0x3f, 0x71, 0xcd, 0x2b, 0xcd, 0x6b,
	0xb3, 0xfb, 0x88, 0x6d, 0xb8, 0x78, 0x0f, 0x51, 0xb7, 0xba, 0x80, 0xa2, 0x5e, 0xe2, 0x67, 0x51,
	0x71, 0xe1, 0x6c, 0x5b, 0xdf, 0x87, 0xb9, 0x95, 0x38, 0x86, 0xa0, 0x35, 0x80, 0xfe, 0x0e, 0x40,
	0xfb, 0x9c, 0x5e, 0x43, 0x60, 0x73, 0x25, 0x10, 0x4d, 0x4b, 0xc5, 0xef, 0x6e, 0xd7, 0xa7, 0x18,
	0x31, 0xd7, 0x2a, 0x7e, 0x08, 0x84, 0x76, 0x1d, 0x3e, 0x40, 0x5c, 0x1b, 0xc4, 0x41, 0x9f, 0x6d,
	0x8e, 0x57, 0x57, 0x74, 0x0b, 0x72, 0xc8, 0x9e, 0x1e, 0xe0, 0x33, 0xa0, 0xdb, 0xfe, 0xa6, 0x25,
	0x2f, 0x3b, 0xdc, 0x62, 0xdc, 0x37, 0x23, 0x23, 0x13, 0x32, 0x0d, 0xcd, 0x8f, 0xb9, 0xf9, 0x0b,
	0x00, 0x00, 0xff, 0xff, 0xe1, 0x76, 0xf1, 0x34, 0x41, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MotionClient is the client API for Motion service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MotionClient interface {
	DetectMotion(ctx context.Context, in *Images, opts ...grpc.CallOption) (*MotionResults, error)
	DetectMotionStream(ctx context.Context, opts ...grpc.CallOption) (Motion_DetectMotionStreamClient, error)
}

type motionClient struct {
	cc *grpc.ClientConn
}

func NewMotionClient(cc *grpc.ClientConn) MotionClient {
	return &motionClient{cc}
}

func (c *motionClient) DetectMotion(ctx context.Context, in *Images, opts ...grpc.CallOption) (*MotionResults, error) {
	out := new(MotionResults)
	err := c.cc.Invoke(ctx, "/processor.Motion/detectMotion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *motionClient) DetectMotionStream(ctx context.Context, opts ...grpc.CallOption) (Motion_DetectMotionStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Motion_serviceDesc.Streams[0], "/processor.Motion/detectMotionStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &motionDetectMotionStreamClient{stream}
	return x, nil
}

type Motion_DetectMotionStreamClient interface {
	Send(*Image) error
	Recv() (*MotionResults, error)
	grpc.ClientStream
}

type motionDetectMotionStreamClient struct {
	grpc.ClientStream
}

func (x *motionDetectMotionStreamClient) Send(m *Image) error {
	return x.ClientStream.SendMsg(m)
}

func (x *motionDetectMotionStreamClient) Recv() (*MotionResults, error) {
	m := new(MotionResults)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MotionServer is the server API for Motion service.
type MotionServer interface {
	DetectMotion(context.Context, *Images) (*MotionResults, error)
	DetectMotionStream(Motion_DetectMotionStreamServer) error
}

// UnimplementedMotionServer can be embedded to have forward compatible implementations.
type UnimplementedMotionServer struct {
}

func (*UnimplementedMotionServer) DetectMotion(ctx context.Context, req *Images) (*MotionResults, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetectMotion not implemented")
}
func (*UnimplementedMotionServer) DetectMotionStream(srv Motion_DetectMotionStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method DetectMotionStream not implemented")
}

func RegisterMotionServer(s *grpc.Server, srv MotionServer) {
	s.RegisterService(&_Motion_serviceDesc, srv)
}

func _Motion_DetectMotion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Images)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MotionServer).DetectMotion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/processor.Motion/DetectMotion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MotionServer).DetectMotion(ctx, req.(*Images))
	}
	return interceptor(ctx, in, info, handler)
}

func _Motion_DetectMotionStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MotionServer).DetectMotionStream(&motionDetectMotionStreamServer{stream})
}

type Motion_DetectMotionStreamServer interface {
	Send(*MotionResults) error
	Recv() (*Image, error)
	grpc.ServerStream
}

type motionDetectMotionStreamServer struct {
	grpc.ServerStream
}

func (x *motionDetectMotionStreamServer) Send(m *MotionResults) error {
	return x.ServerStream.SendMsg(m)
}

func (x *motionDetectMotionStreamServer) Recv() (*Image, error) {
	m := new(Image)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Motion_serviceDesc = grpc.ServiceDesc{
	ServiceName: "processor.Motion",
	HandlerType: (*MotionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "detectMotion",
			Handler:    _Motion_DetectMotion_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "detectMotionStream",
			Handler:       _Motion_DetectMotionStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "motion.proto",
}
