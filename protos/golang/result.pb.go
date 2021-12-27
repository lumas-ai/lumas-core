// Code generated by protoc-gen-go. DO NOT EDIT.
// source: result.proto

package lumas

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
	return fileDescriptor_4feee897733d2100, []int{0}
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

func init() {
	proto.RegisterType((*Result)(nil), "lumas.Result")
}

func init() { proto.RegisterFile("result.proto", fileDescriptor_4feee897733d2100) }

var fileDescriptor_4feee897733d2100 = []byte{
	// 119 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x4a, 0x2d, 0x2e,
	0xcd, 0x29, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcd, 0x29, 0xcd, 0x4d, 0x2c, 0x56,
	0x4a, 0xe0, 0x62, 0x0b, 0x02, 0x0b, 0x0b, 0xc9, 0x71, 0x71, 0x15, 0x97, 0x26, 0x27, 0xa7, 0x16,
	0x17, 0xa7, 0x95, 0xe6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x04, 0x21, 0x89, 0x08, 0xc9, 0x70,
	0x71, 0xa6, 0x16, 0x15, 0xe5, 0x17, 0x79, 0x67, 0xe6, 0xa5, 0x48, 0x30, 0x29, 0x30, 0x6a, 0x70,
	0x06, 0x21, 0x04, 0x84, 0x24, 0xb8, 0xd8, 0x73, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53, 0x25, 0x98,
	0xc1, 0x72, 0x30, 0x6e, 0x12, 0x1b, 0xd8, 0x3e, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf5,
	0x3e, 0x02, 0x1a, 0x7f, 0x00, 0x00, 0x00,
}