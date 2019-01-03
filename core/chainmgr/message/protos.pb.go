// Code generated by protoc-gen-go.
// source: protos.proto
// DO NOT EDIT!

/*
Package message is a generated protocol buffer package.

It is generated from these files:
	protos.proto

It has these top-level messages:
	CrossShardMsg
*/
package message

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import actor "github.com/ontio/ontology-eventbus/actor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CrossShardMsg struct {
	Version int32      `protobuf:"varint,1,opt,name=Version" json:"Version,omitempty"`
	Type    int32      `protobuf:"varint,2,opt,name=Type" json:"Type,omitempty"`
	Sender  *actor.PID `protobuf:"bytes,3,opt,name=Sender" json:"Sender,omitempty"`
	Data    []byte     `protobuf:"bytes,4,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (m *CrossShardMsg) Reset()                    { *m = CrossShardMsg{} }
func (m *CrossShardMsg) String() string            { return proto.CompactTextString(m) }
func (*CrossShardMsg) ProtoMessage()               {}
func (*CrossShardMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *CrossShardMsg) GetSender() *actor.PID {
	if m != nil {
		return m.Sender
	}
	return nil
}

func init() {
	proto.RegisterType((*CrossShardMsg)(nil), "message.CrossShardMsg")
}

func init() { proto.RegisterFile("protos.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 185 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x28, 0xca, 0x2f,
	0xc9, 0x2f, 0xd6, 0x03, 0x53, 0x42, 0xec, 0xb9, 0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9, 0x52, 0xa6,
	0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0xf9, 0x79, 0x25, 0x99, 0xf9,
	0x20, 0x32, 0x3f, 0x27, 0x3f, 0xbd, 0x52, 0x37, 0xb5, 0x2c, 0x35, 0xaf, 0x24, 0xa9, 0xb4, 0x58,
	0x3f, 0x31, 0xb9, 0x24, 0xbf, 0x48, 0x1f, 0x59, 0xbf, 0x52, 0x29, 0x17, 0xaf, 0x73, 0x51, 0x7e,
	0x71, 0x71, 0x70, 0x46, 0x62, 0x51, 0x8a, 0x6f, 0x71, 0xba, 0x90, 0x04, 0x17, 0x7b, 0x58, 0x6a,
	0x51, 0x71, 0x66, 0x7e, 0x9e, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x6b, 0x10, 0x8c, 0x2b, 0x24, 0xc4,
	0xc5, 0x12, 0x52, 0x59, 0x90, 0x2a, 0xc1, 0x04, 0x16, 0x06, 0xb3, 0x85, 0x94, 0xb8, 0xd8, 0x82,
	0x53, 0xf3, 0x52, 0x52, 0x8b, 0x24, 0x98, 0x15, 0x18, 0x35, 0xb8, 0x8d, 0xb8, 0xf4, 0xc0, 0x76,
	0xe8, 0x05, 0x78, 0xba, 0x04, 0x41, 0x65, 0x40, 0xfa, 0x5c, 0x12, 0x4b, 0x12, 0x25, 0x58, 0x14,
	0x18, 0x35, 0x78, 0x82, 0xc0, 0xec, 0x24, 0x36, 0xb0, 0xed, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x6b, 0x3d, 0x0d, 0x1f, 0xcd, 0x00, 0x00, 0x00,
}
