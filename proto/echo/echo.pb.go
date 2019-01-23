// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/echo/echo.proto

package echo

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EchoInbound struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EchoInbound) Reset()         { *m = EchoInbound{} }
func (m *EchoInbound) String() string { return proto.CompactTextString(m) }
func (*EchoInbound) ProtoMessage()    {}
func (*EchoInbound) Descriptor() ([]byte, []int) {
	return fileDescriptor_echo_b76e9c412c4e5a98, []int{0}
}
func (m *EchoInbound) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoInbound.Unmarshal(m, b)
}
func (m *EchoInbound) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoInbound.Marshal(b, m, deterministic)
}
func (dst *EchoInbound) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoInbound.Merge(dst, src)
}
func (m *EchoInbound) XXX_Size() int {
	return xxx_messageInfo_EchoInbound.Size(m)
}
func (m *EchoInbound) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoInbound.DiscardUnknown(m)
}

var xxx_messageInfo_EchoInbound proto.InternalMessageInfo

func (m *EchoInbound) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type EchoOutbound struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EchoOutbound) Reset()         { *m = EchoOutbound{} }
func (m *EchoOutbound) String() string { return proto.CompactTextString(m) }
func (*EchoOutbound) ProtoMessage()    {}
func (*EchoOutbound) Descriptor() ([]byte, []int) {
	return fileDescriptor_echo_b76e9c412c4e5a98, []int{1}
}
func (m *EchoOutbound) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoOutbound.Unmarshal(m, b)
}
func (m *EchoOutbound) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoOutbound.Marshal(b, m, deterministic)
}
func (dst *EchoOutbound) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoOutbound.Merge(dst, src)
}
func (m *EchoOutbound) XXX_Size() int {
	return xxx_messageInfo_EchoOutbound.Size(m)
}
func (m *EchoOutbound) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoOutbound.DiscardUnknown(m)
}

var xxx_messageInfo_EchoOutbound proto.InternalMessageInfo

func (m *EchoOutbound) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*EchoInbound)(nil), "echo.EchoInbound")
	proto.RegisterType((*EchoOutbound)(nil), "echo.EchoOutbound")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EchoClient is the client API for Echo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EchoClient interface {
	Echo(ctx context.Context, in *EchoInbound, opts ...grpc.CallOption) (*EchoOutbound, error)
}

type echoClient struct {
	cc *grpc.ClientConn
}

func NewEchoClient(cc *grpc.ClientConn) EchoClient {
	return &echoClient{cc}
}

func (c *echoClient) Echo(ctx context.Context, in *EchoInbound, opts ...grpc.CallOption) (*EchoOutbound, error) {
	out := new(EchoOutbound)
	err := c.cc.Invoke(ctx, "/echo.Echo/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EchoServer is the server API for Echo service.
type EchoServer interface {
	Echo(context.Context, *EchoInbound) (*EchoOutbound, error)
}

func RegisterEchoServer(s *grpc.Server, srv EchoServer) {
	s.RegisterService(&_Echo_serviceDesc, srv)
}

func _Echo_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoInbound)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/echo.Echo/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServer).Echo(ctx, req.(*EchoInbound))
	}
	return interceptor(ctx, in, info, handler)
}

var _Echo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "echo.Echo",
	HandlerType: (*EchoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _Echo_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/echo/echo.proto",
}

func init() { proto.RegisterFile("proto/echo/echo.proto", fileDescriptor_echo_b76e9c412c4e5a98) }

var fileDescriptor_echo_b76e9c412c4e5a98 = []byte{
	// 116 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0x4d, 0xce, 0x80, 0x10, 0x7a, 0x60, 0xbe, 0x10, 0x0b, 0x88, 0xad, 0xa4, 0xce,
	0xc5, 0xed, 0x9a, 0x9c, 0x91, 0xef, 0x99, 0x97, 0x94, 0x5f, 0x9a, 0x97, 0x22, 0x24, 0xc1, 0xc5,
	0x9e, 0x9b, 0x5a, 0x5c, 0x9c, 0x98, 0x9e, 0x2a, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0xe3,
	0x2a, 0x69, 0x70, 0xf1, 0x80, 0x14, 0xfa, 0x97, 0x96, 0x10, 0x50, 0x69, 0x64, 0xca, 0xc5, 0x02,
	0x52, 0x29, 0xa4, 0x0b, 0xa5, 0x05, 0xf5, 0xc0, 0xb6, 0x22, 0x59, 0x23, 0x25, 0x84, 0x10, 0x82,
	0x19, 0x98, 0xc4, 0x06, 0x76, 0x96, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x5b, 0x74, 0x28, 0x1c,
	0xaf, 0x00, 0x00, 0x00,
}
