// Code generated by protoc-gen-go.
// source: protos/golang.conradwood.net/apis/goeasyops/goeasyops.proto
// DO NOT EDIT!

/*
Package goeasyops is a generated protocol buffer package.

It is generated from these files:
	protos/golang.conradwood.net/apis/goeasyops/goeasyops.proto

It has these top-level messages:
	InContext
	ImmutableContext
	MutableContext
	CTXRoutingTags
	GRPCErrorList
	GRPCError
*/
package goeasyops

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import auth "golang.conradwood.net/apis/auth"
import _ "golang.conradwood.net/apis/common"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// this is transported within the Context of a call. This proto is on-purpose very brief and static so that old versions of go-easyops can handle it
type InContext struct {
	ImCtx *ImmutableContext `protobuf:"bytes,1,opt,name=ImCtx" json:"ImCtx,omitempty"`
	MCtx  *MutableContext   `protobuf:"bytes,2,opt,name=MCtx" json:"MCtx,omitempty"`
}

func (m *InContext) Reset()                    { *m = InContext{} }
func (m *InContext) String() string            { return proto.CompactTextString(m) }
func (*InContext) ProtoMessage()               {}
func (*InContext) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *InContext) GetImCtx() *ImmutableContext {
	if m != nil {
		return m.ImCtx
	}
	return nil
}

func (m *InContext) GetMCtx() *MutableContext {
	if m != nil {
		return m.MCtx
	}
	return nil
}

// this must not be changed throughout its lifetime. furthermore, go-easyops _must_ transport this as-is (preserving unknown fields)
type ImmutableContext struct {
	RequestID      string              `protobuf:"bytes,1,opt,name=RequestID" json:"RequestID,omitempty"`
	CreatorService *auth.SignedUser    `protobuf:"bytes,2,opt,name=CreatorService" json:"CreatorService,omitempty"`
	User           *auth.SignedUser    `protobuf:"bytes,3,opt,name=User" json:"User,omitempty"`
	SudoUser       *auth.SignedUser    `protobuf:"bytes,4,opt,name=SudoUser" json:"SudoUser,omitempty"`
	Session        *auth.SignedSession `protobuf:"bytes,5,opt,name=Session" json:"Session,omitempty"`
}

func (m *ImmutableContext) Reset()                    { *m = ImmutableContext{} }
func (m *ImmutableContext) String() string            { return proto.CompactTextString(m) }
func (*ImmutableContext) ProtoMessage()               {}
func (*ImmutableContext) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ImmutableContext) GetRequestID() string {
	if m != nil {
		return m.RequestID
	}
	return ""
}

func (m *ImmutableContext) GetCreatorService() *auth.SignedUser {
	if m != nil {
		return m.CreatorService
	}
	return nil
}

func (m *ImmutableContext) GetUser() *auth.SignedUser {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *ImmutableContext) GetSudoUser() *auth.SignedUser {
	if m != nil {
		return m.SudoUser
	}
	return nil
}

func (m *ImmutableContext) GetSession() *auth.SignedSession {
	if m != nil {
		return m.Session
	}
	return nil
}

// this may change. fields are not guaranteed to be preserved
type MutableContext struct {
	CallingService *auth.SignedUser `protobuf:"bytes,1,opt,name=CallingService" json:"CallingService,omitempty"`
	Debug          bool             `protobuf:"varint,2,opt,name=Debug" json:"Debug,omitempty"`
	Trace          bool             `protobuf:"varint,3,opt,name=Trace" json:"Trace,omitempty"`
	Tags           *CTXRoutingTags  `protobuf:"bytes,4,opt,name=Tags" json:"Tags,omitempty"`
}

func (m *MutableContext) Reset()                    { *m = MutableContext{} }
func (m *MutableContext) String() string            { return proto.CompactTextString(m) }
func (*MutableContext) ProtoMessage()               {}
func (*MutableContext) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *MutableContext) GetCallingService() *auth.SignedUser {
	if m != nil {
		return m.CallingService
	}
	return nil
}

func (m *MutableContext) GetDebug() bool {
	if m != nil {
		return m.Debug
	}
	return false
}

func (m *MutableContext) GetTrace() bool {
	if m != nil {
		return m.Trace
	}
	return false
}

func (m *MutableContext) GetTags() *CTXRoutingTags {
	if m != nil {
		return m.Tags
	}
	return nil
}

//
// Routing tags are part of a context. Rules to match tags when looking for a suitable target service instance are as follows:
//
// If FallbackToPlain == false then:
//
// - if context has no tags - use any instance
//
// - if context has tags - only use instance that matches exactly (with all tags)
//
// If FallbackToPlain == true then:
//
// - if context has no tags - use any instance
//
// - if context has tags - if at least one instance matches exactly (with all tags), use only that. if none matches, but at least one instance has no tags, use that.
//
// Propagate: If it is true, the routing tags are kept as-is, otherwise the first target service will strip routing tags out
//
type CTXRoutingTags struct {
	Tags            map[string]string `protobuf:"bytes,1,rep,name=Tags" json:"Tags,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	FallbackToPlain bool              `protobuf:"varint,2,opt,name=FallbackToPlain" json:"FallbackToPlain,omitempty"`
	Propagate       bool              `protobuf:"varint,3,opt,name=Propagate" json:"Propagate,omitempty"`
}

func (m *CTXRoutingTags) Reset()                    { *m = CTXRoutingTags{} }
func (m *CTXRoutingTags) String() string            { return proto.CompactTextString(m) }
func (*CTXRoutingTags) ProtoMessage()               {}
func (*CTXRoutingTags) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CTXRoutingTags) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *CTXRoutingTags) GetFallbackToPlain() bool {
	if m != nil {
		return m.FallbackToPlain
	}
	return false
}

func (m *CTXRoutingTags) GetPropagate() bool {
	if m != nil {
		return m.Propagate
	}
	return false
}

// a single error can only hold a single proto
type GRPCErrorList struct {
	Errors []*GRPCError `protobuf:"bytes,1,rep,name=Errors" json:"Errors,omitempty"`
}

func (m *GRPCErrorList) Reset()                    { *m = GRPCErrorList{} }
func (m *GRPCErrorList) String() string            { return proto.CompactTextString(m) }
func (*GRPCErrorList) ProtoMessage()               {}
func (*GRPCErrorList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GRPCErrorList) GetErrors() []*GRPCError {
	if m != nil {
		return m.Errors
	}
	return nil
}

type GRPCError struct {
	UserMessage         string `protobuf:"bytes,1,opt,name=UserMessage" json:"UserMessage,omitempty"`
	LogMessage          string `protobuf:"bytes,2,opt,name=LogMessage" json:"LogMessage,omitempty"`
	MethodName          string `protobuf:"bytes,3,opt,name=MethodName" json:"MethodName,omitempty"`
	ServiceName         string `protobuf:"bytes,4,opt,name=ServiceName" json:"ServiceName,omitempty"`
	CallingServiceID    string `protobuf:"bytes,5,opt,name=CallingServiceID" json:"CallingServiceID,omitempty"`
	CallingServiceEmail string `protobuf:"bytes,6,opt,name=CallingServiceEmail" json:"CallingServiceEmail,omitempty"`
}

func (m *GRPCError) Reset()                    { *m = GRPCError{} }
func (m *GRPCError) String() string            { return proto.CompactTextString(m) }
func (*GRPCError) ProtoMessage()               {}
func (*GRPCError) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *GRPCError) GetUserMessage() string {
	if m != nil {
		return m.UserMessage
	}
	return ""
}

func (m *GRPCError) GetLogMessage() string {
	if m != nil {
		return m.LogMessage
	}
	return ""
}

func (m *GRPCError) GetMethodName() string {
	if m != nil {
		return m.MethodName
	}
	return ""
}

func (m *GRPCError) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *GRPCError) GetCallingServiceID() string {
	if m != nil {
		return m.CallingServiceID
	}
	return ""
}

func (m *GRPCError) GetCallingServiceEmail() string {
	if m != nil {
		return m.CallingServiceEmail
	}
	return ""
}

func init() {
	proto.RegisterType((*InContext)(nil), "goeasyops.InContext")
	proto.RegisterType((*ImmutableContext)(nil), "goeasyops.ImmutableContext")
	proto.RegisterType((*MutableContext)(nil), "goeasyops.MutableContext")
	proto.RegisterType((*CTXRoutingTags)(nil), "goeasyops.CTXRoutingTags")
	proto.RegisterType((*GRPCErrorList)(nil), "goeasyops.GRPCErrorList")
	proto.RegisterType((*GRPCError)(nil), "goeasyops.GRPCError")
}

func init() {
	proto.RegisterFile("protos/golang.conradwood.net/apis/goeasyops/goeasyops.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 556 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x94, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0xc7, 0x95, 0xf5, 0x83, 0xe6, 0x54, 0x94, 0xca, 0xdb, 0x45, 0x29, 0x08, 0x4a, 0xe1, 0xa2,
	0x9a, 0xb6, 0x0c, 0xc6, 0xc5, 0x26, 0x10, 0x37, 0xb4, 0x05, 0x55, 0x5a, 0x51, 0xe5, 0x16, 0x89,
	0x5b, 0xb7, 0xb5, 0xb2, 0x68, 0x89, 0x5d, 0x6c, 0x67, 0xac, 0x0f, 0xc4, 0xcb, 0xf0, 0x2e, 0xbc,
	0x03, 0xf2, 0x47, 0xd3, 0xa4, 0xac, 0xbb, 0x49, 0xe2, 0xff, 0xf9, 0xf9, 0xf8, 0xfc, 0x7d, 0x8e,
	0x02, 0x1f, 0x57, 0x82, 0x2b, 0x2e, 0xcf, 0x42, 0x1e, 0x13, 0x16, 0x06, 0x0b, 0xce, 0x04, 0x59,
	0xfe, 0xe2, 0x7c, 0x19, 0x30, 0xaa, 0xce, 0xc8, 0x2a, 0xd2, 0x21, 0x4a, 0xe4, 0x9a, 0xaf, 0x72,
	0x5f, 0x81, 0xd9, 0x85, 0xfc, 0x4c, 0x68, 0x1f, 0x3f, 0x90, 0x80, 0xa4, 0xea, 0xda, 0x3c, 0xec,
	0xb6, 0x76, 0xf0, 0x00, 0xbb, 0xe0, 0x49, 0xc2, 0x99, 0x7b, 0x59, 0xbe, 0x9b, 0x80, 0x3f, 0x62,
	0x7d, 0xce, 0x14, 0xbd, 0x53, 0xe8, 0x1d, 0x54, 0x46, 0x49, 0x5f, 0xdd, 0xb5, 0xbc, 0x8e, 0xd7,
	0xab, 0x9f, 0x3f, 0x0b, 0xb6, 0x45, 0x8d, 0x92, 0x24, 0x55, 0x64, 0x1e, 0x53, 0xc7, 0x62, 0x4b,
	0xa2, 0x53, 0x28, 0x8f, 0xf5, 0x8e, 0x03, 0xb3, 0xe3, 0x69, 0x6e, 0xc7, 0xb8, 0xc8, 0x1b, 0xac,
	0xfb, 0xd7, 0x83, 0xe6, 0x6e, 0x2a, 0xf4, 0x1c, 0x7c, 0x4c, 0x7f, 0xa6, 0x54, 0xaa, 0xd1, 0xc0,
	0x1c, 0xed, 0xe3, 0xad, 0x80, 0x2e, 0xa1, 0xd1, 0x17, 0x94, 0x28, 0x2e, 0xa6, 0x54, 0xdc, 0x46,
	0x0b, 0xea, 0xce, 0x6a, 0x06, 0xc6, 0xf6, 0x34, 0x0a, 0x19, 0x5d, 0x7e, 0x97, 0x54, 0xe0, 0x1d,
	0x0e, 0xbd, 0x81, 0xb2, 0xd6, 0x5b, 0xa5, 0x3d, 0xbc, 0x89, 0xa2, 0x13, 0xa8, 0x4d, 0xd3, 0x25,
	0x37, 0x64, 0x79, 0x0f, 0x99, 0x11, 0xe8, 0x14, 0x1e, 0x4d, 0xa9, 0x94, 0x11, 0x67, 0xad, 0x8a,
	0x81, 0x0f, 0xf3, 0xb0, 0x0b, 0xe1, 0x0d, 0xd3, 0xfd, 0xed, 0x41, 0xa3, 0x78, 0x11, 0xc6, 0x0f,
	0x89, 0xe3, 0x88, 0x85, 0x1b, 0x3f, 0xde, 0x5e, 0x3f, 0x05, 0x0e, 0x1d, 0x41, 0x65, 0x40, 0xe7,
	0x69, 0x68, 0x2e, 0xa0, 0x86, 0xed, 0x42, 0xab, 0x33, 0x41, 0x16, 0xd4, 0xd8, 0xac, 0x61, 0xbb,
	0xd0, 0x7d, 0x99, 0x91, 0x50, 0x3a, 0x47, 0xf9, 0xbe, 0xf4, 0x67, 0x3f, 0x30, 0x4f, 0x55, 0xc4,
	0x42, 0x0d, 0x60, 0x83, 0x75, 0xff, 0x78, 0xd0, 0x28, 0x06, 0xd0, 0x85, 0xcb, 0xe0, 0x75, 0x4a,
	0xbd, 0xfa, 0xf9, 0xeb, 0xbd, 0x19, 0x02, 0xfd, 0x18, 0x32, 0x25, 0xd6, 0x36, 0x17, 0xea, 0xc1,
	0x93, 0x2f, 0x24, 0x8e, 0xe7, 0x64, 0x71, 0x33, 0xe3, 0x93, 0x98, 0x44, 0xcc, 0x15, 0xbc, 0x2b,
	0xeb, 0xc6, 0x4f, 0x04, 0x5f, 0x91, 0x90, 0xa8, 0x4d, 0xf9, 0x5b, 0xa1, 0x7d, 0x01, 0x7e, 0x96,
	0x1a, 0x35, 0xa1, 0x74, 0x43, 0xd7, 0x6e, 0x3a, 0xf4, 0xa7, 0xf6, 0x7d, 0x4b, 0xe2, 0xd4, 0x8e,
	0x83, 0x8f, 0xed, 0xe2, 0xc3, 0xc1, 0xa5, 0xd7, 0xfd, 0x04, 0x8f, 0xbf, 0xe2, 0x49, 0x7f, 0x28,
	0x04, 0x17, 0x57, 0x91, 0x54, 0xe8, 0x04, 0xaa, 0x66, 0xb1, 0x31, 0x73, 0x94, 0x33, 0x93, 0x91,
	0xd8, 0x31, 0x7a, 0x46, 0xfd, 0x4c, 0x45, 0x1d, 0xa8, 0xeb, 0x66, 0x8c, 0xa9, 0x94, 0x24, 0xa4,
	0xae, 0x80, 0xbc, 0x84, 0x5e, 0x00, 0x5c, 0xf1, 0x70, 0x03, 0xd8, 0x6a, 0x72, 0x8a, 0x8e, 0x8f,
	0xa9, 0xba, 0xe6, 0xcb, 0x6f, 0x24, 0xb1, 0x36, 0x7d, 0x9c, 0x53, 0xf4, 0x09, 0xae, 0xc3, 0x06,
	0x28, 0xdb, 0x13, 0x72, 0x12, 0x3a, 0x86, 0x66, 0x71, 0x14, 0x46, 0x03, 0x33, 0x7d, 0x3e, 0xfe,
	0x4f, 0x47, 0x6f, 0xe1, 0xb0, 0xa8, 0x0d, 0x13, 0x12, 0xc5, 0xad, 0xaa, 0xc1, 0xef, 0x0b, 0x7d,
	0x7e, 0x05, 0x2f, 0xef, 0xfd, 0x69, 0x6c, 0x2f, 0x69, 0x5e, 0x35, 0x3f, 0x8b, 0xf7, 0xff, 0x02,
	0x00, 0x00, 0xff, 0xff, 0x96, 0x38, 0x23, 0x61, 0xd2, 0x04, 0x00, 0x00,
}
