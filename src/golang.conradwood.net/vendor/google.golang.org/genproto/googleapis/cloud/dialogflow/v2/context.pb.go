// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/dialogflow/v2/context.proto

package dialogflow

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	_struct "github.com/golang/protobuf/ptypes/struct"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// Represents a context.
type Context struct {
	// Required. The unique identifier of the context. Format:
	// `projects/<Project ID>/agent/sessions/<Session ID>/contexts/<Context ID>`.
	//
	// The `Context ID` is always converted to lowercase, may only contain
	// characters in [a-zA-Z0-9_-%] and may be at most 250 bytes long.
	//
	// The following context names are reserved for internal use by Dialogflow.
	// You should not use these contexts or create contexts with these names:
	//
	// * `__system_counters__`
	// * `*_id_dialog_context`
	// * `*_dialog_params_size`
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Optional. The number of conversational query requests after which the
	// context expires. If set to `0` (the default) the context expires
	// immediately. Contexts expire automatically after 20 minutes if there
	// are no matching queries.
	LifespanCount int32 `protobuf:"varint,2,opt,name=lifespan_count,json=lifespanCount,proto3" json:"lifespan_count,omitempty"`
	// Optional. The collection of parameters associated with this context.
	// Refer to [this
	// doc](https://cloud.google.com/dialogflow/docs/intents-actions-parameters)
	// for syntax.
	Parameters           *_struct.Struct `protobuf:"bytes,3,opt,name=parameters,proto3" json:"parameters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Context) Reset()         { *m = Context{} }
func (m *Context) String() string { return proto.CompactTextString(m) }
func (*Context) ProtoMessage()    {}
func (*Context) Descriptor() ([]byte, []int) {
	return fileDescriptor_e7e2e3bf8515c3b3, []int{0}
}

func (m *Context) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Context.Unmarshal(m, b)
}
func (m *Context) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Context.Marshal(b, m, deterministic)
}
func (m *Context) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Context.Merge(m, src)
}
func (m *Context) XXX_Size() int {
	return xxx_messageInfo_Context.Size(m)
}
func (m *Context) XXX_DiscardUnknown() {
	xxx_messageInfo_Context.DiscardUnknown(m)
}

var xxx_messageInfo_Context proto.InternalMessageInfo

func (m *Context) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Context) GetLifespanCount() int32 {
	if m != nil {
		return m.LifespanCount
	}
	return 0
}

func (m *Context) GetParameters() *_struct.Struct {
	if m != nil {
		return m.Parameters
	}
	return nil
}

// The request message for [Contexts.ListContexts][google.cloud.dialogflow.v2.Contexts.ListContexts].
type ListContextsRequest struct {
	// Required. The session to list all contexts from.
	// Format: `projects/<Project ID>/agent/sessions/<Session ID>`.
	Parent string `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
	// Optional. The maximum number of items to return in a single page. By
	// default 100 and at most 1000.
	PageSize int32 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// Optional. The next_page_token value returned from a previous list request.
	PageToken            string   `protobuf:"bytes,3,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListContextsRequest) Reset()         { *m = ListContextsRequest{} }
func (m *ListContextsRequest) String() string { return proto.CompactTextString(m) }
func (*ListContextsRequest) ProtoMessage()    {}
func (*ListContextsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e7e2e3bf8515c3b3, []int{1}
}

func (m *ListContextsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListContextsRequest.Unmarshal(m, b)
}
func (m *ListContextsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListContextsRequest.Marshal(b, m, deterministic)
}
func (m *ListContextsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListContextsRequest.Merge(m, src)
}
func (m *ListContextsRequest) XXX_Size() int {
	return xxx_messageInfo_ListContextsRequest.Size(m)
}
func (m *ListContextsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListContextsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListContextsRequest proto.InternalMessageInfo

func (m *ListContextsRequest) GetParent() string {
	if m != nil {
		return m.Parent
	}
	return ""
}

func (m *ListContextsRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ListContextsRequest) GetPageToken() string {
	if m != nil {
		return m.PageToken
	}
	return ""
}

// The response message for [Contexts.ListContexts][google.cloud.dialogflow.v2.Contexts.ListContexts].
type ListContextsResponse struct {
	// The list of contexts. There will be a maximum number of items
	// returned based on the page_size field in the request.
	Contexts []*Context `protobuf:"bytes,1,rep,name=contexts,proto3" json:"contexts,omitempty"`
	// Token to retrieve the next page of results, or empty if there are no
	// more results in the list.
	NextPageToken        string   `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListContextsResponse) Reset()         { *m = ListContextsResponse{} }
func (m *ListContextsResponse) String() string { return proto.CompactTextString(m) }
func (*ListContextsResponse) ProtoMessage()    {}
func (*ListContextsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e7e2e3bf8515c3b3, []int{2}
}

func (m *ListContextsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListContextsResponse.Unmarshal(m, b)
}
func (m *ListContextsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListContextsResponse.Marshal(b, m, deterministic)
}
func (m *ListContextsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListContextsResponse.Merge(m, src)
}
func (m *ListContextsResponse) XXX_Size() int {
	return xxx_messageInfo_ListContextsResponse.Size(m)
}
func (m *ListContextsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListContextsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListContextsResponse proto.InternalMessageInfo

func (m *ListContextsResponse) GetContexts() []*Context {
	if m != nil {
		return m.Contexts
	}
	return nil
}

func (m *ListContextsResponse) GetNextPageToken() string {
	if m != nil {
		return m.NextPageToken
	}
	return ""
}

// The request message for [Contexts.GetContext][google.cloud.dialogflow.v2.Contexts.GetContext].
type GetContextRequest struct {
	// Required. The name of the context. Format:
	// `projects/<Project ID>/agent/sessions/<Session ID>/contexts/<Context ID>`.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetContextRequest) Reset()         { *m = GetContextRequest{} }
func (m *GetContextRequest) String() string { return proto.CompactTextString(m) }
func (*GetContextRequest) ProtoMessage()    {}
func (*GetContextRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e7e2e3bf8515c3b3, []int{3}
}

func (m *GetContextRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetContextRequest.Unmarshal(m, b)
}
func (m *GetContextRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetContextRequest.Marshal(b, m, deterministic)
}
func (m *GetContextRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetContextRequest.Merge(m, src)
}
func (m *GetContextRequest) XXX_Size() int {
	return xxx_messageInfo_GetContextRequest.Size(m)
}
func (m *GetContextRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetContextRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetContextRequest proto.InternalMessageInfo

func (m *GetContextRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// The request message for [Contexts.CreateContext][google.cloud.dialogflow.v2.Contexts.CreateContext].
type CreateContextRequest struct {
	// Required. The session to create a context for.
	// Format: `projects/<Project ID>/agent/sessions/<Session ID>`.
	Parent string `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
	// Required. The context to create.
	Context              *Context `protobuf:"bytes,2,opt,name=context,proto3" json:"context,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateContextRequest) Reset()         { *m = CreateContextRequest{} }
func (m *CreateContextRequest) String() string { return proto.CompactTextString(m) }
func (*CreateContextRequest) ProtoMessage()    {}
func (*CreateContextRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e7e2e3bf8515c3b3, []int{4}
}

func (m *CreateContextRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateContextRequest.Unmarshal(m, b)
}
func (m *CreateContextRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateContextRequest.Marshal(b, m, deterministic)
}
func (m *CreateContextRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateContextRequest.Merge(m, src)
}
func (m *CreateContextRequest) XXX_Size() int {
	return xxx_messageInfo_CreateContextRequest.Size(m)
}
func (m *CreateContextRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateContextRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateContextRequest proto.InternalMessageInfo

func (m *CreateContextRequest) GetParent() string {
	if m != nil {
		return m.Parent
	}
	return ""
}

func (m *CreateContextRequest) GetContext() *Context {
	if m != nil {
		return m.Context
	}
	return nil
}

// The request message for [Contexts.UpdateContext][google.cloud.dialogflow.v2.Contexts.UpdateContext].
type UpdateContextRequest struct {
	// Required. The context to update.
	Context *Context `protobuf:"bytes,1,opt,name=context,proto3" json:"context,omitempty"`
	// Optional. The mask to control which fields get updated.
	UpdateMask           *field_mask.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *UpdateContextRequest) Reset()         { *m = UpdateContextRequest{} }
func (m *UpdateContextRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateContextRequest) ProtoMessage()    {}
func (*UpdateContextRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e7e2e3bf8515c3b3, []int{5}
}

func (m *UpdateContextRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateContextRequest.Unmarshal(m, b)
}
func (m *UpdateContextRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateContextRequest.Marshal(b, m, deterministic)
}
func (m *UpdateContextRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateContextRequest.Merge(m, src)
}
func (m *UpdateContextRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateContextRequest.Size(m)
}
func (m *UpdateContextRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateContextRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateContextRequest proto.InternalMessageInfo

func (m *UpdateContextRequest) GetContext() *Context {
	if m != nil {
		return m.Context
	}
	return nil
}

func (m *UpdateContextRequest) GetUpdateMask() *field_mask.FieldMask {
	if m != nil {
		return m.UpdateMask
	}
	return nil
}

// The request message for [Contexts.DeleteContext][google.cloud.dialogflow.v2.Contexts.DeleteContext].
type DeleteContextRequest struct {
	// Required. The name of the context to delete. Format:
	// `projects/<Project ID>/agent/sessions/<Session ID>/contexts/<Context ID>`.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteContextRequest) Reset()         { *m = DeleteContextRequest{} }
func (m *DeleteContextRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteContextRequest) ProtoMessage()    {}
func (*DeleteContextRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e7e2e3bf8515c3b3, []int{6}
}

func (m *DeleteContextRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteContextRequest.Unmarshal(m, b)
}
func (m *DeleteContextRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteContextRequest.Marshal(b, m, deterministic)
}
func (m *DeleteContextRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteContextRequest.Merge(m, src)
}
func (m *DeleteContextRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteContextRequest.Size(m)
}
func (m *DeleteContextRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteContextRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteContextRequest proto.InternalMessageInfo

func (m *DeleteContextRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// The request message for [Contexts.DeleteAllContexts][google.cloud.dialogflow.v2.Contexts.DeleteAllContexts].
type DeleteAllContextsRequest struct {
	// Required. The name of the session to delete all contexts from. Format:
	// `projects/<Project ID>/agent/sessions/<Session ID>`.
	Parent               string   `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteAllContextsRequest) Reset()         { *m = DeleteAllContextsRequest{} }
func (m *DeleteAllContextsRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteAllContextsRequest) ProtoMessage()    {}
func (*DeleteAllContextsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e7e2e3bf8515c3b3, []int{7}
}

func (m *DeleteAllContextsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteAllContextsRequest.Unmarshal(m, b)
}
func (m *DeleteAllContextsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteAllContextsRequest.Marshal(b, m, deterministic)
}
func (m *DeleteAllContextsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteAllContextsRequest.Merge(m, src)
}
func (m *DeleteAllContextsRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteAllContextsRequest.Size(m)
}
func (m *DeleteAllContextsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteAllContextsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteAllContextsRequest proto.InternalMessageInfo

func (m *DeleteAllContextsRequest) GetParent() string {
	if m != nil {
		return m.Parent
	}
	return ""
}

func init() {
	proto.RegisterType((*Context)(nil), "google.cloud.dialogflow.v2.Context")
	proto.RegisterType((*ListContextsRequest)(nil), "google.cloud.dialogflow.v2.ListContextsRequest")
	proto.RegisterType((*ListContextsResponse)(nil), "google.cloud.dialogflow.v2.ListContextsResponse")
	proto.RegisterType((*GetContextRequest)(nil), "google.cloud.dialogflow.v2.GetContextRequest")
	proto.RegisterType((*CreateContextRequest)(nil), "google.cloud.dialogflow.v2.CreateContextRequest")
	proto.RegisterType((*UpdateContextRequest)(nil), "google.cloud.dialogflow.v2.UpdateContextRequest")
	proto.RegisterType((*DeleteContextRequest)(nil), "google.cloud.dialogflow.v2.DeleteContextRequest")
	proto.RegisterType((*DeleteAllContextsRequest)(nil), "google.cloud.dialogflow.v2.DeleteAllContextsRequest")
}

func init() {
	proto.RegisterFile("google/cloud/dialogflow/v2/context.proto", fileDescriptor_e7e2e3bf8515c3b3)
}

var fileDescriptor_e7e2e3bf8515c3b3 = []byte{
	// 871 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0xcf, 0x6f, 0xe3, 0x44,
	0x14, 0xd6, 0xa4, 0xb0, 0x9b, 0x4c, 0x37, 0xa0, 0x1d, 0x22, 0x9a, 0xf5, 0xae, 0x20, 0x78, 0x25,
	0x14, 0x22, 0xd6, 0x03, 0x59, 0x24, 0xb4, 0x65, 0x17, 0x70, 0x12, 0x5a, 0x09, 0x81, 0x54, 0xa5,
	0x14, 0x09, 0x2e, 0xd1, 0xd4, 0x99, 0xb8, 0xa6, 0x8e, 0xc7, 0x78, 0xc6, 0x4d, 0x68, 0x14, 0x21,
	0xf1, 0x2f, 0xc0, 0xad, 0x3d, 0x71, 0xec, 0x11, 0x89, 0x7f, 0x82, 0x23, 0xdc, 0x38, 0xf5, 0xc0,
	0x89, 0x23, 0xe2, 0x80, 0x38, 0x21, 0xdb, 0xe3, 0xd8, 0x89, 0x53, 0x37, 0x29, 0x37, 0xfb, 0xfd,
	0xfc, 0xe6, 0x7b, 0xef, 0x1b, 0x1b, 0xd6, 0x4d, 0xc6, 0x4c, 0x9b, 0x62, 0xc3, 0x66, 0x7e, 0x1f,
	0xf7, 0x2d, 0x62, 0x33, 0x73, 0x60, 0xb3, 0x11, 0x3e, 0x69, 0x62, 0x83, 0x39, 0x82, 0x8e, 0x85,
	0xe6, 0x7a, 0x4c, 0x30, 0xa4, 0x44, 0x91, 0x5a, 0x18, 0xa9, 0x25, 0x91, 0xda, 0x49, 0x53, 0x79,
	0x20, 0xab, 0x10, 0xd7, 0xc2, 0xc4, 0x71, 0x98, 0x20, 0xc2, 0x62, 0x0e, 0x8f, 0x32, 0x95, 0xad,
	0x94, 0xd7, 0xb0, 0x2d, 0xea, 0xc8, 0x92, 0xca, 0xab, 0x29, 0xc7, 0xc0, 0xa2, 0x76, 0xbf, 0x77,
	0x48, 0x8f, 0xc8, 0x89, 0xc5, 0x3c, 0x19, 0x70, 0x2f, 0x15, 0xe0, 0x51, 0xce, 0x7c, 0xcf, 0xa0,
	0xd2, 0x75, 0x5f, 0xba, 0xc2, 0xb7, 0x43, 0x7f, 0x80, 0xe9, 0xd0, 0x15, 0xdf, 0x48, 0x67, 0x6d,
	0xd1, 0x19, 0x55, 0x1f, 0x12, 0x7e, 0x2c, 0x23, 0x1e, 0x2c, 0x46, 0x70, 0xe1, 0xf9, 0x86, 0x04,
	0xa6, 0xfe, 0x0d, 0xe0, 0xed, 0x76, 0x74, 0x7a, 0xb4, 0x05, 0x9f, 0x73, 0xc8, 0x90, 0x56, 0x41,
	0x0d, 0xd4, 0x4b, 0xad, 0x8d, 0x4b, 0xbd, 0xd0, 0x0d, 0x0d, 0xa8, 0x01, 0x5f, 0xb0, 0xad, 0x01,
	0xe5, 0x2e, 0x71, 0x7a, 0x06, 0xf3, 0x1d, 0x51, 0x2d, 0xd4, 0x40, 0xfd, 0xf9, 0x20, 0x04, 0x74,
	0xcb, 0xb1, 0xab, 0x1d, 0x78, 0xd0, 0x53, 0x08, 0x5d, 0xe2, 0x91, 0x21, 0x15, 0xd4, 0xe3, 0xd5,
	0x8d, 0x1a, 0xa8, 0x6f, 0x36, 0xb7, 0x34, 0xc9, 0x68, 0x8c, 0x41, 0xdb, 0x0f, 0x31, 0x44, 0x05,
	0x52, 0xf1, 0xdb, 0x83, 0x3f, 0x75, 0x03, 0xbe, 0x96, 0xa2, 0x3c, 0x4a, 0x24, 0xae, 0xc5, 0x35,
	0x83, 0x0d, 0x71, 0x0c, 0xf5, 0x7d, 0xd7, 0x63, 0x5f, 0x51, 0x43, 0x70, 0x3c, 0x91, 0x4f, 0x53,
	0x4c, 0x4c, 0xea, 0x08, 0xcc, 0x29, 0xe7, 0xc1, 0x44, 0xf0, 0x44, 0x3e, 0x4d, 0xe3, 0xe9, 0x72,
	0x3c, 0x91, 0x4f, 0x53, 0x75, 0x0c, 0x5f, 0xfa, 0xc4, 0xe2, 0x42, 0x96, 0xe3, 0x5d, 0xfa, 0xb5,
	0x4f, 0xb9, 0x40, 0xf7, 0xe1, 0x2d, 0x97, 0x78, 0xd4, 0x11, 0x69, 0x0e, 0xa4, 0x09, 0xd5, 0x60,
	0xc9, 0x25, 0x26, 0xed, 0x71, 0xeb, 0x94, 0xa6, 0x09, 0x28, 0x06, 0xd6, 0x7d, 0xeb, 0x94, 0x22,
	0x35, 0x38, 0xbb, 0x49, 0x7b, 0x82, 0x1d, 0x53, 0x27, 0x3c, 0x7b, 0x29, 0x0a, 0x09, 0x13, 0x3f,
	0x0b, 0xac, 0xea, 0xb7, 0xb0, 0x32, 0xdf, 0x99, 0xbb, 0xcc, 0xe1, 0x14, 0x7d, 0x00, 0x8b, 0x31,
	0xce, 0x2a, 0xa8, 0x6d, 0xd4, 0x37, 0x9b, 0x0f, 0xb5, 0xab, 0xf7, 0x50, 0x93, 0xf9, 0xdd, 0x59,
	0x12, 0x7a, 0x1d, 0xbe, 0xe8, 0xd0, 0xb1, 0xe8, 0xa5, 0x10, 0x04, 0x20, 0x4b, 0xdd, 0x72, 0x60,
	0xde, 0x9b, 0x01, 0xe8, 0xc2, 0xbb, 0xbb, 0x34, 0xee, 0x1f, 0x1f, 0xfc, 0xd9, 0xdc, 0xe8, 0xdf,
	0xb8, 0xd4, 0x0b, 0xff, 0xea, 0x0f, 0x57, 0x18, 0x44, 0xb4, 0x20, 0xea, 0x08, 0x56, 0xda, 0x1e,
	0x25, 0x82, 0x2e, 0x94, 0xcd, 0xe5, 0xb3, 0x05, 0x6f, 0x4b, 0xf0, 0x21, 0xd0, 0xd5, 0x0e, 0x1c,
	0x95, 0x88, 0x13, 0xd5, 0x73, 0x00, 0x2b, 0x07, 0x6e, 0x3f, 0xdb, 0x39, 0x55, 0x1c, 0xdc, 0xb0,
	0x38, 0xfa, 0x10, 0x6e, 0xfa, 0x61, 0xed, 0x50, 0x4e, 0x12, 0xa4, 0x92, 0xd9, 0xe5, 0x9d, 0x40,
	0x71, 0x9f, 0x12, 0x7e, 0x2c, 0xd7, 0x39, 0xca, 0x09, 0x0c, 0xea, 0x01, 0xac, 0x74, 0xa8, 0x4d,
	0x33, 0xe8, 0xfe, 0x27, 0xdd, 0xef, 0xc2, 0x6a, 0x54, 0x56, 0xb7, 0xed, 0x75, 0x56, 0xb8, 0xf9,
	0x57, 0x11, 0x16, 0xe3, 0x04, 0xf4, 0x33, 0x80, 0x77, 0xd2, 0xab, 0x88, 0x70, 0x1e, 0x45, 0x4b,
	0xe4, 0xa2, 0xbc, 0xb5, 0x7a, 0x42, 0xb4, 0xe5, 0x6a, 0xeb, 0x77, 0x5d, 0x62, 0xf9, 0xee, 0xb7,
	0x3f, 0xbe, 0x2f, 0x3c, 0x46, 0x6f, 0x07, 0xd7, 0xef, 0x24, 0x32, 0x3d, 0x9b, 0x89, 0xba, 0xb1,
	0x28, 0xe6, 0x46, 0xa2, 0x62, 0x74, 0x06, 0x20, 0x4c, 0x36, 0x18, 0x3d, 0xca, 0x03, 0x91, 0xd9,
	0x74, 0x65, 0x95, 0x3d, 0x50, 0x9f, 0xcc, 0xa1, 0x0b, 0x38, 0xcf, 0xc3, 0x96, 0x5c, 0x30, 0x8d,
	0x29, 0xba, 0x00, 0xb0, 0x3c, 0xa7, 0x05, 0x94, 0xcb, 0xd2, 0x32, 0xd9, 0xac, 0x86, 0xb1, 0x15,
	0x62, 0x7c, 0xaa, 0xae, 0xcf, 0xe0, 0xf6, 0x6c, 0xc3, 0x7f, 0x02, 0xb0, 0x3c, 0x27, 0x9f, 0x7c,
	0xb0, 0xcb, 0x94, 0xb6, 0x1a, 0xd8, 0x8f, 0x43, 0xb0, 0x9d, 0xe6, 0x93, 0x10, 0x6c, 0xfc, 0xb9,
	0x5d, 0x87, 0xd8, 0x04, 0xf4, 0x0f, 0x00, 0x96, 0xe7, 0x54, 0x95, 0x0f, 0x7a, 0x99, 0x00, 0x95,
	0x97, 0x33, 0x2a, 0xfe, 0x28, 0xf8, 0xa8, 0xc6, 0x83, 0x6f, 0xdc, 0x60, 0xf0, 0xe7, 0x00, 0xde,
	0xcd, 0xa8, 0x12, 0xbd, 0x73, 0x3d, 0xb4, 0xac, 0x88, 0x57, 0x84, 0xb7, 0xce, 0xcc, 0x95, 0xf1,
	0x2f, 0xfa, 0xbd, 0x2b, 0xef, 0x97, 0x5f, 0xf5, 0x2f, 0x8e, 0x84, 0x70, 0xf9, 0x36, 0xc6, 0xa3,
	0x51, 0xe6, 0xf2, 0x21, 0xbe, 0x38, 0x8a, 0x7e, 0x9b, 0x1e, 0xb9, 0x36, 0x11, 0x03, 0xe6, 0x0d,
	0xdf, 0xbc, 0x2e, 0x3c, 0x69, 0xd5, 0x3a, 0x03, 0xf0, 0x15, 0x83, 0x0d, 0x73, 0x88, 0x68, 0xdd,
	0x91, 0x04, 0xec, 0x05, 0xc7, 0xdd, 0x03, 0x5f, 0x76, 0x64, 0xac, 0xc9, 0x6c, 0xe2, 0x98, 0x1a,
	0xf3, 0x4c, 0x6c, 0x52, 0x27, 0x24, 0x03, 0x27, 0xdd, 0x96, 0xfd, 0xca, 0xbd, 0x97, 0xbc, 0xfd,
	0x03, 0xc0, 0x8f, 0x85, 0x42, 0x67, 0xe7, 0xa2, 0xa0, 0xec, 0x46, 0xe5, 0xda, 0x61, 0xeb, 0x4e,
	0xd2, 0xfa, 0xf3, 0xe6, 0xe1, 0xad, 0xb0, 0xea, 0xe3, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0xe2,
	0xc0, 0x75, 0xc6, 0x1f, 0x0a, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ContextsClient is the client API for Contexts service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ContextsClient interface {
	// Returns the list of all contexts in the specified session.
	ListContexts(ctx context.Context, in *ListContextsRequest, opts ...grpc.CallOption) (*ListContextsResponse, error)
	// Retrieves the specified context.
	GetContext(ctx context.Context, in *GetContextRequest, opts ...grpc.CallOption) (*Context, error)
	// Creates a context.
	//
	// If the specified context already exists, overrides the context.
	CreateContext(ctx context.Context, in *CreateContextRequest, opts ...grpc.CallOption) (*Context, error)
	// Updates the specified context.
	UpdateContext(ctx context.Context, in *UpdateContextRequest, opts ...grpc.CallOption) (*Context, error)
	// Deletes the specified context.
	DeleteContext(ctx context.Context, in *DeleteContextRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Deletes all active contexts in the specified session.
	DeleteAllContexts(ctx context.Context, in *DeleteAllContextsRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type contextsClient struct {
	cc *grpc.ClientConn
}

func NewContextsClient(cc *grpc.ClientConn) ContextsClient {
	return &contextsClient{cc}
}

func (c *contextsClient) ListContexts(ctx context.Context, in *ListContextsRequest, opts ...grpc.CallOption) (*ListContextsResponse, error) {
	out := new(ListContextsResponse)
	err := c.cc.Invoke(ctx, "/google.cloud.dialogflow.v2.Contexts/ListContexts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextsClient) GetContext(ctx context.Context, in *GetContextRequest, opts ...grpc.CallOption) (*Context, error) {
	out := new(Context)
	err := c.cc.Invoke(ctx, "/google.cloud.dialogflow.v2.Contexts/GetContext", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextsClient) CreateContext(ctx context.Context, in *CreateContextRequest, opts ...grpc.CallOption) (*Context, error) {
	out := new(Context)
	err := c.cc.Invoke(ctx, "/google.cloud.dialogflow.v2.Contexts/CreateContext", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextsClient) UpdateContext(ctx context.Context, in *UpdateContextRequest, opts ...grpc.CallOption) (*Context, error) {
	out := new(Context)
	err := c.cc.Invoke(ctx, "/google.cloud.dialogflow.v2.Contexts/UpdateContext", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextsClient) DeleteContext(ctx context.Context, in *DeleteContextRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/google.cloud.dialogflow.v2.Contexts/DeleteContext", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextsClient) DeleteAllContexts(ctx context.Context, in *DeleteAllContextsRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/google.cloud.dialogflow.v2.Contexts/DeleteAllContexts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ContextsServer is the server API for Contexts service.
type ContextsServer interface {
	// Returns the list of all contexts in the specified session.
	ListContexts(context.Context, *ListContextsRequest) (*ListContextsResponse, error)
	// Retrieves the specified context.
	GetContext(context.Context, *GetContextRequest) (*Context, error)
	// Creates a context.
	//
	// If the specified context already exists, overrides the context.
	CreateContext(context.Context, *CreateContextRequest) (*Context, error)
	// Updates the specified context.
	UpdateContext(context.Context, *UpdateContextRequest) (*Context, error)
	// Deletes the specified context.
	DeleteContext(context.Context, *DeleteContextRequest) (*empty.Empty, error)
	// Deletes all active contexts in the specified session.
	DeleteAllContexts(context.Context, *DeleteAllContextsRequest) (*empty.Empty, error)
}

// UnimplementedContextsServer can be embedded to have forward compatible implementations.
type UnimplementedContextsServer struct {
}

func (*UnimplementedContextsServer) ListContexts(ctx context.Context, req *ListContextsRequest) (*ListContextsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListContexts not implemented")
}
func (*UnimplementedContextsServer) GetContext(ctx context.Context, req *GetContextRequest) (*Context, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetContext not implemented")
}
func (*UnimplementedContextsServer) CreateContext(ctx context.Context, req *CreateContextRequest) (*Context, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateContext not implemented")
}
func (*UnimplementedContextsServer) UpdateContext(ctx context.Context, req *UpdateContextRequest) (*Context, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateContext not implemented")
}
func (*UnimplementedContextsServer) DeleteContext(ctx context.Context, req *DeleteContextRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteContext not implemented")
}
func (*UnimplementedContextsServer) DeleteAllContexts(ctx context.Context, req *DeleteAllContextsRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAllContexts not implemented")
}

func RegisterContextsServer(s *grpc.Server, srv ContextsServer) {
	s.RegisterService(&_Contexts_serviceDesc, srv)
}

func _Contexts_ListContexts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListContextsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextsServer).ListContexts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.cloud.dialogflow.v2.Contexts/ListContexts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextsServer).ListContexts(ctx, req.(*ListContextsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Contexts_GetContext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetContextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextsServer).GetContext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.cloud.dialogflow.v2.Contexts/GetContext",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextsServer).GetContext(ctx, req.(*GetContextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Contexts_CreateContext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateContextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextsServer).CreateContext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.cloud.dialogflow.v2.Contexts/CreateContext",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextsServer).CreateContext(ctx, req.(*CreateContextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Contexts_UpdateContext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateContextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextsServer).UpdateContext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.cloud.dialogflow.v2.Contexts/UpdateContext",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextsServer).UpdateContext(ctx, req.(*UpdateContextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Contexts_DeleteContext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteContextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextsServer).DeleteContext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.cloud.dialogflow.v2.Contexts/DeleteContext",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextsServer).DeleteContext(ctx, req.(*DeleteContextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Contexts_DeleteAllContexts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAllContextsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextsServer).DeleteAllContexts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.cloud.dialogflow.v2.Contexts/DeleteAllContexts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextsServer).DeleteAllContexts(ctx, req.(*DeleteAllContextsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Contexts_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.cloud.dialogflow.v2.Contexts",
	HandlerType: (*ContextsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListContexts",
			Handler:    _Contexts_ListContexts_Handler,
		},
		{
			MethodName: "GetContext",
			Handler:    _Contexts_GetContext_Handler,
		},
		{
			MethodName: "CreateContext",
			Handler:    _Contexts_CreateContext_Handler,
		},
		{
			MethodName: "UpdateContext",
			Handler:    _Contexts_UpdateContext_Handler,
		},
		{
			MethodName: "DeleteContext",
			Handler:    _Contexts_DeleteContext_Handler,
		},
		{
			MethodName: "DeleteAllContexts",
			Handler:    _Contexts_DeleteAllContexts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/cloud/dialogflow/v2/context.proto",
}
