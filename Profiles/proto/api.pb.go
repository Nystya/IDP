// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package api

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

type Void struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Void) Reset()         { *m = Void{} }
func (m *Void) String() string { return proto.CompactTextString(m) }
func (*Void) ProtoMessage()    {}
func (*Void) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *Void) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Void.Unmarshal(m, b)
}
func (m *Void) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Void.Marshal(b, m, deterministic)
}
func (m *Void) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Void.Merge(m, src)
}
func (m *Void) XXX_Size() int {
	return xxx_messageInfo_Void.Size(m)
}
func (m *Void) XXX_DiscardUnknown() {
	xxx_messageInfo_Void.DiscardUnknown(m)
}

var xxx_messageInfo_Void proto.InternalMessageInfo

type ErrorCode struct {
	Err                  int32    `protobuf:"varint,1,opt,name=err,proto3" json:"err,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ErrorCode) Reset()         { *m = ErrorCode{} }
func (m *ErrorCode) String() string { return proto.CompactTextString(m) }
func (*ErrorCode) ProtoMessage()    {}
func (*ErrorCode) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *ErrorCode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ErrorCode.Unmarshal(m, b)
}
func (m *ErrorCode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ErrorCode.Marshal(b, m, deterministic)
}
func (m *ErrorCode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorCode.Merge(m, src)
}
func (m *ErrorCode) XXX_Size() int {
	return xxx_messageInfo_ErrorCode.Size(m)
}
func (m *ErrorCode) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorCode.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorCode proto.InternalMessageInfo

func (m *ErrorCode) GetErr() int32 {
	if m != nil {
		return m.Err
	}
	return 0
}

func (m *ErrorCode) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type ID struct {
	ID                   string   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ID) Reset()         { *m = ID{} }
func (m *ID) String() string { return proto.CompactTextString(m) }
func (*ID) ProtoMessage()    {}
func (*ID) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}

func (m *ID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ID.Unmarshal(m, b)
}
func (m *ID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ID.Marshal(b, m, deterministic)
}
func (m *ID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ID.Merge(m, src)
}
func (m *ID) XXX_Size() int {
	return xxx_messageInfo_ID.Size(m)
}
func (m *ID) XXX_DiscardUnknown() {
	xxx_messageInfo_ID.DiscardUnknown(m)
}

var xxx_messageInfo_ID proto.InternalMessageInfo

func (m *ID) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

type ServiceCategory struct {
	ID                   *ID      `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Service              string   `protobuf:"bytes,2,opt,name=service,proto3" json:"service,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServiceCategory) Reset()         { *m = ServiceCategory{} }
func (m *ServiceCategory) String() string { return proto.CompactTextString(m) }
func (*ServiceCategory) ProtoMessage()    {}
func (*ServiceCategory) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}

func (m *ServiceCategory) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServiceCategory.Unmarshal(m, b)
}
func (m *ServiceCategory) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServiceCategory.Marshal(b, m, deterministic)
}
func (m *ServiceCategory) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceCategory.Merge(m, src)
}
func (m *ServiceCategory) XXX_Size() int {
	return xxx_messageInfo_ServiceCategory.Size(m)
}
func (m *ServiceCategory) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceCategory.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceCategory proto.InternalMessageInfo

func (m *ServiceCategory) GetID() *ID {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *ServiceCategory) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

type SkillCategory struct {
	ID                   *ID      `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Category             string   `protobuf:"bytes,2,opt,name=category,proto3" json:"category,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SkillCategory) Reset()         { *m = SkillCategory{} }
func (m *SkillCategory) String() string { return proto.CompactTextString(m) }
func (*SkillCategory) ProtoMessage()    {}
func (*SkillCategory) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{4}
}

func (m *SkillCategory) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SkillCategory.Unmarshal(m, b)
}
func (m *SkillCategory) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SkillCategory.Marshal(b, m, deterministic)
}
func (m *SkillCategory) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SkillCategory.Merge(m, src)
}
func (m *SkillCategory) XXX_Size() int {
	return xxx_messageInfo_SkillCategory.Size(m)
}
func (m *SkillCategory) XXX_DiscardUnknown() {
	xxx_messageInfo_SkillCategory.DiscardUnknown(m)
}

var xxx_messageInfo_SkillCategory proto.InternalMessageInfo

func (m *SkillCategory) GetID() *ID {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *SkillCategory) GetCategory() string {
	if m != nil {
		return m.Category
	}
	return ""
}

type Skill struct {
	ID                   *ID            `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Category             *SkillCategory `protobuf:"bytes,2,opt,name=category,proto3" json:"category,omitempty"`
	Skill                string         `protobuf:"bytes,3,opt,name=skill,proto3" json:"skill,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Skill) Reset()         { *m = Skill{} }
func (m *Skill) String() string { return proto.CompactTextString(m) }
func (*Skill) ProtoMessage()    {}
func (*Skill) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{5}
}

func (m *Skill) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Skill.Unmarshal(m, b)
}
func (m *Skill) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Skill.Marshal(b, m, deterministic)
}
func (m *Skill) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Skill.Merge(m, src)
}
func (m *Skill) XXX_Size() int {
	return xxx_messageInfo_Skill.Size(m)
}
func (m *Skill) XXX_DiscardUnknown() {
	xxx_messageInfo_Skill.DiscardUnknown(m)
}

var xxx_messageInfo_Skill proto.InternalMessageInfo

func (m *Skill) GetID() *ID {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *Skill) GetCategory() *SkillCategory {
	if m != nil {
		return m.Category
	}
	return nil
}

func (m *Skill) GetSkill() string {
	if m != nil {
		return m.Skill
	}
	return ""
}

type EmployerProfile struct {
	EUID                 *ID      `protobuf:"bytes,1,opt,name=EUID,proto3" json:"EUID,omitempty"`
	Phone                string   `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
	LastName             string   `protobuf:"bytes,3,opt,name=lastName,proto3" json:"lastName,omitempty"`
	FirstName            string   `protobuf:"bytes,4,opt,name=firstName,proto3" json:"firstName,omitempty"`
	Rating               float32  `protobuf:"fixed32,5,opt,name=rating,proto3" json:"rating,omitempty"`
	JobsPosted           int32    `protobuf:"varint,6,opt,name=jobsPosted,proto3" json:"jobsPosted,omitempty"`
	MoneySpent           float32  `protobuf:"fixed32,7,opt,name=moneySpent,proto3" json:"moneySpent,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmployerProfile) Reset()         { *m = EmployerProfile{} }
func (m *EmployerProfile) String() string { return proto.CompactTextString(m) }
func (*EmployerProfile) ProtoMessage()    {}
func (*EmployerProfile) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{6}
}

func (m *EmployerProfile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmployerProfile.Unmarshal(m, b)
}
func (m *EmployerProfile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmployerProfile.Marshal(b, m, deterministic)
}
func (m *EmployerProfile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmployerProfile.Merge(m, src)
}
func (m *EmployerProfile) XXX_Size() int {
	return xxx_messageInfo_EmployerProfile.Size(m)
}
func (m *EmployerProfile) XXX_DiscardUnknown() {
	xxx_messageInfo_EmployerProfile.DiscardUnknown(m)
}

var xxx_messageInfo_EmployerProfile proto.InternalMessageInfo

func (m *EmployerProfile) GetEUID() *ID {
	if m != nil {
		return m.EUID
	}
	return nil
}

func (m *EmployerProfile) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *EmployerProfile) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *EmployerProfile) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *EmployerProfile) GetRating() float32 {
	if m != nil {
		return m.Rating
	}
	return 0
}

func (m *EmployerProfile) GetJobsPosted() int32 {
	if m != nil {
		return m.JobsPosted
	}
	return 0
}

func (m *EmployerProfile) GetMoneySpent() float32 {
	if m != nil {
		return m.MoneySpent
	}
	return 0
}

type EditEmployerProfileRequest struct {
	EUID                 *ID      `protobuf:"bytes,1,opt,name=EUID,proto3" json:"EUID,omitempty"`
	Phone                string   `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
	LastName             string   `protobuf:"bytes,3,opt,name=lastName,proto3" json:"lastName,omitempty"`
	FirstName            string   `protobuf:"bytes,4,opt,name=firstName,proto3" json:"firstName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EditEmployerProfileRequest) Reset()         { *m = EditEmployerProfileRequest{} }
func (m *EditEmployerProfileRequest) String() string { return proto.CompactTextString(m) }
func (*EditEmployerProfileRequest) ProtoMessage()    {}
func (*EditEmployerProfileRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{7}
}

func (m *EditEmployerProfileRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EditEmployerProfileRequest.Unmarshal(m, b)
}
func (m *EditEmployerProfileRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EditEmployerProfileRequest.Marshal(b, m, deterministic)
}
func (m *EditEmployerProfileRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EditEmployerProfileRequest.Merge(m, src)
}
func (m *EditEmployerProfileRequest) XXX_Size() int {
	return xxx_messageInfo_EditEmployerProfileRequest.Size(m)
}
func (m *EditEmployerProfileRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EditEmployerProfileRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EditEmployerProfileRequest proto.InternalMessageInfo

func (m *EditEmployerProfileRequest) GetEUID() *ID {
	if m != nil {
		return m.EUID
	}
	return nil
}

func (m *EditEmployerProfileRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *EditEmployerProfileRequest) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *EditEmployerProfileRequest) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

type FreelancerProfile struct {
	FUID                 *ID              `protobuf:"bytes,1,opt,name=FUID,proto3" json:"FUID,omitempty"`
	Phone                string           `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
	LastName             string           `protobuf:"bytes,3,opt,name=lastName,proto3" json:"lastName,omitempty"`
	FirstName            string           `protobuf:"bytes,4,opt,name=firstName,proto3" json:"firstName,omitempty"`
	Rating               float32          `protobuf:"fixed32,5,opt,name=rating,proto3" json:"rating,omitempty"`
	Balance              float32          `protobuf:"fixed32,6,opt,name=balance,proto3" json:"balance,omitempty"`
	Description          string           `protobuf:"bytes,7,opt,name=description,proto3" json:"description,omitempty"`
	Photo                string           `protobuf:"bytes,8,opt,name=photo,proto3" json:"photo,omitempty"`
	SkillCategories      []*SkillCategory `protobuf:"bytes,9,rep,name=skillCategories,proto3" json:"skillCategories,omitempty"`
	Skills               []*Skill         `protobuf:"bytes,10,rep,name=skills,proto3" json:"skills,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *FreelancerProfile) Reset()         { *m = FreelancerProfile{} }
func (m *FreelancerProfile) String() string { return proto.CompactTextString(m) }
func (*FreelancerProfile) ProtoMessage()    {}
func (*FreelancerProfile) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{8}
}

func (m *FreelancerProfile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FreelancerProfile.Unmarshal(m, b)
}
func (m *FreelancerProfile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FreelancerProfile.Marshal(b, m, deterministic)
}
func (m *FreelancerProfile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FreelancerProfile.Merge(m, src)
}
func (m *FreelancerProfile) XXX_Size() int {
	return xxx_messageInfo_FreelancerProfile.Size(m)
}
func (m *FreelancerProfile) XXX_DiscardUnknown() {
	xxx_messageInfo_FreelancerProfile.DiscardUnknown(m)
}

var xxx_messageInfo_FreelancerProfile proto.InternalMessageInfo

func (m *FreelancerProfile) GetFUID() *ID {
	if m != nil {
		return m.FUID
	}
	return nil
}

func (m *FreelancerProfile) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *FreelancerProfile) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *FreelancerProfile) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *FreelancerProfile) GetRating() float32 {
	if m != nil {
		return m.Rating
	}
	return 0
}

func (m *FreelancerProfile) GetBalance() float32 {
	if m != nil {
		return m.Balance
	}
	return 0
}

func (m *FreelancerProfile) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *FreelancerProfile) GetPhoto() string {
	if m != nil {
		return m.Photo
	}
	return ""
}

func (m *FreelancerProfile) GetSkillCategories() []*SkillCategory {
	if m != nil {
		return m.SkillCategories
	}
	return nil
}

func (m *FreelancerProfile) GetSkills() []*Skill {
	if m != nil {
		return m.Skills
	}
	return nil
}

type EditFreelancerProfileRequest struct {
	FUID                 *ID              `protobuf:"bytes,1,opt,name=FUID,proto3" json:"FUID,omitempty"`
	Phone                string           `protobuf:"bytes,2,opt,name=phone,proto3" json:"phone,omitempty"`
	LastName             string           `protobuf:"bytes,3,opt,name=lastName,proto3" json:"lastName,omitempty"`
	FirstName            string           `protobuf:"bytes,4,opt,name=firstName,proto3" json:"firstName,omitempty"`
	Description          string           `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	SkillCategories      []*SkillCategory `protobuf:"bytes,6,rep,name=skillCategories,proto3" json:"skillCategories,omitempty"`
	Skills               []*Skill         `protobuf:"bytes,7,rep,name=skills,proto3" json:"skills,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *EditFreelancerProfileRequest) Reset()         { *m = EditFreelancerProfileRequest{} }
func (m *EditFreelancerProfileRequest) String() string { return proto.CompactTextString(m) }
func (*EditFreelancerProfileRequest) ProtoMessage()    {}
func (*EditFreelancerProfileRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{9}
}

func (m *EditFreelancerProfileRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EditFreelancerProfileRequest.Unmarshal(m, b)
}
func (m *EditFreelancerProfileRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EditFreelancerProfileRequest.Marshal(b, m, deterministic)
}
func (m *EditFreelancerProfileRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EditFreelancerProfileRequest.Merge(m, src)
}
func (m *EditFreelancerProfileRequest) XXX_Size() int {
	return xxx_messageInfo_EditFreelancerProfileRequest.Size(m)
}
func (m *EditFreelancerProfileRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EditFreelancerProfileRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EditFreelancerProfileRequest proto.InternalMessageInfo

func (m *EditFreelancerProfileRequest) GetFUID() *ID {
	if m != nil {
		return m.FUID
	}
	return nil
}

func (m *EditFreelancerProfileRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *EditFreelancerProfileRequest) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *EditFreelancerProfileRequest) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *EditFreelancerProfileRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *EditFreelancerProfileRequest) GetSkillCategories() []*SkillCategory {
	if m != nil {
		return m.SkillCategories
	}
	return nil
}

func (m *EditFreelancerProfileRequest) GetSkills() []*Skill {
	if m != nil {
		return m.Skills
	}
	return nil
}

func init() {
	proto.RegisterType((*Void)(nil), "api.Void")
	proto.RegisterType((*ErrorCode)(nil), "api.ErrorCode")
	proto.RegisterType((*ID)(nil), "api.ID")
	proto.RegisterType((*ServiceCategory)(nil), "api.ServiceCategory")
	proto.RegisterType((*SkillCategory)(nil), "api.SkillCategory")
	proto.RegisterType((*Skill)(nil), "api.Skill")
	proto.RegisterType((*EmployerProfile)(nil), "api.EmployerProfile")
	proto.RegisterType((*EditEmployerProfileRequest)(nil), "api.EditEmployerProfileRequest")
	proto.RegisterType((*FreelancerProfile)(nil), "api.FreelancerProfile")
	proto.RegisterType((*EditFreelancerProfileRequest)(nil), "api.EditFreelancerProfileRequest")
}

func init() {
	proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c)
}

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 525 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x55, 0xcd, 0x8e, 0xd3, 0x30,
	0x10, 0x56, 0xd2, 0x36, 0x6d, 0xa6, 0x62, 0x0b, 0xa6, 0x2c, 0x56, 0x59, 0x41, 0xf1, 0xa9, 0xa7,
	0x22, 0x75, 0x0f, 0x5c, 0xb8, 0x6d, 0xda, 0x55, 0x39, 0xa0, 0x55, 0x2a, 0xb8, 0xa7, 0x8d, 0x5b,
	0x0c, 0x69, 0x1c, 0x6c, 0x83, 0xd4, 0x47, 0x40, 0x42, 0xe2, 0xc5, 0x38, 0xf3, 0x3c, 0xc8, 0x4e,
	0x9a, 0x66, 0xf3, 0x23, 0xa4, 0x3d, 0xf4, 0x12, 0x79, 0x66, 0xbe, 0xf9, 0x66, 0xe6, 0x1b, 0x27,
	0x01, 0x37, 0x48, 0xd8, 0x34, 0x11, 0x5c, 0x71, 0xd4, 0x0a, 0x12, 0x46, 0x1c, 0x68, 0x7f, 0xe2,
	0x2c, 0x24, 0x6f, 0xc0, 0x9d, 0x0b, 0xc1, 0xc5, 0x0d, 0x0f, 0x29, 0x7a, 0x0c, 0x2d, 0x2a, 0x04,
	0xb6, 0xc6, 0xd6, 0xa4, 0xe3, 0xeb, 0xa3, 0xf6, 0xec, 0xe5, 0x0e, 0xdb, 0x63, 0x6b, 0xe2, 0xfa,
	0xfa, 0x48, 0x86, 0x60, 0x2f, 0x3d, 0x74, 0xa1, 0x9f, 0x06, 0xe8, 0xfa, 0xf6, 0xd2, 0x23, 0x1e,
	0x0c, 0x56, 0x54, 0xfc, 0x60, 0x1b, 0x7a, 0x13, 0x28, 0xba, 0xe3, 0xe2, 0x80, 0x9e, 0xe7, 0x90,
	0xfe, 0xac, 0x3b, 0xd5, 0xe5, 0x97, 0x9e, 0xc6, 0x22, 0x0c, 0x5d, 0x99, 0x62, 0x33, 0xde, 0xa3,
	0x49, 0x3c, 0x78, 0xb4, 0xfa, 0xca, 0xa2, 0xe8, 0xff, 0x1c, 0x23, 0xe8, 0x6d, 0x32, 0x50, 0x46,
	0x92, 0xdb, 0x64, 0x0b, 0x1d, 0xc3, 0xd2, 0x9c, 0x3d, 0x2d, 0x65, 0xf7, 0x67, 0xc8, 0x84, 0xef,
	0x15, 0x3f, 0x31, 0xa2, 0x21, 0x74, 0xa4, 0x0e, 0xe1, 0x96, 0x29, 0x95, 0x1a, 0xe4, 0xaf, 0x05,
	0x83, 0xf9, 0x3e, 0x89, 0xf8, 0x81, 0x8a, 0x3b, 0xc1, 0xb7, 0x2c, 0xa2, 0xe8, 0x05, 0xb4, 0xe7,
	0x1f, 0xab, 0x45, 0x8d, 0x53, 0xd3, 0x24, 0x9f, 0x79, 0x7c, 0x1c, 0x3b, 0x35, 0xf4, 0x28, 0x51,
	0x20, 0xd5, 0x87, 0x60, 0x4f, 0x33, 0xfe, 0xdc, 0x46, 0x57, 0xe0, 0x6e, 0x99, 0xc8, 0x82, 0x6d,
	0x13, 0x3c, 0x39, 0xd0, 0x25, 0x38, 0x22, 0x50, 0x2c, 0xde, 0xe1, 0xce, 0xd8, 0x9a, 0xd8, 0x7e,
	0x66, 0xa1, 0x97, 0x00, 0x5f, 0xf8, 0x5a, 0xde, 0x71, 0xa9, 0x68, 0x88, 0x1d, 0xb3, 0xcd, 0x82,
	0x47, 0xc7, 0xf7, 0x3c, 0xa6, 0x87, 0x55, 0x42, 0x63, 0x85, 0xbb, 0x26, 0xb7, 0xe0, 0x21, 0x3f,
	0x2d, 0x18, 0xcd, 0x43, 0xa6, 0x4a, 0xc3, 0xf9, 0xf4, 0xdb, 0x77, 0x2a, 0xd5, 0x59, 0x67, 0x24,
	0x7f, 0x6c, 0x78, 0xb2, 0x10, 0x94, 0x46, 0x41, 0xbc, 0xb9, 0x27, 0xf3, 0xa2, 0xae, 0x85, 0xc5,
	0x39, 0x65, 0xc6, 0xd0, 0x5d, 0x07, 0xa6, 0x2f, 0xa3, 0xb1, 0xed, 0x1f, 0x4d, 0x34, 0x86, 0x7e,
	0x48, 0xe5, 0x46, 0xb0, 0x44, 0x31, 0x1e, 0x1b, 0x85, 0x5d, 0xbf, 0xe8, 0xca, 0x7a, 0x54, 0x1c,
	0xf7, 0xf2, 0x1e, 0x15, 0x47, 0xef, 0x60, 0x20, 0x0b, 0x57, 0x90, 0x51, 0x89, 0xdd, 0x71, 0xab,
	0xe1, 0x7a, 0x96, 0xa1, 0x88, 0x80, 0x63, 0x5c, 0x12, 0x83, 0x49, 0x82, 0x53, 0x92, 0x9f, 0x45,
	0xc8, 0x6f, 0x1b, 0xae, 0xf4, 0x6a, 0x2b, 0x92, 0x16, 0x96, 0x7b, 0x3e, 0x65, 0x4b, 0x3a, 0x75,
	0xaa, 0x3a, 0xd5, 0x28, 0xe2, 0x3c, 0x44, 0x91, 0x6e, 0x93, 0x22, 0xb3, 0x5f, 0x36, 0xf4, 0x32,
	0x0d, 0x24, 0x5a, 0xc0, 0xd3, 0x9a, 0x8b, 0x8f, 0x5e, 0x99, 0xbc, 0xe6, 0x57, 0x62, 0x74, 0x91,
	0x02, 0xf2, 0x0f, 0xe9, 0x7b, 0x78, 0x56, 0xab, 0x32, 0x7a, 0x9d, 0x33, 0x35, 0x6d, 0xa0, 0xc2,
	0x75, 0x0d, 0xe8, 0x96, 0x56, 0x5a, 0x3a, 0x6e, 0x66, 0x34, 0x4c, 0xe1, 0xa5, 0xf0, 0x5b, 0x18,
	0xde, 0xd2, 0x9a, 0xfa, 0x79, 0xda, 0xa5, 0x39, 0x54, 0x00, 0x6b, 0xc7, 0xfc, 0x23, 0xae, 0xff,
	0x05, 0x00, 0x00, 0xff, 0xff, 0x83, 0xda, 0xde, 0x9e, 0x30, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ProfilesClient is the client API for Profiles service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProfilesClient interface {
	EditEmployerProfile(ctx context.Context, in *EditEmployerProfileRequest, opts ...grpc.CallOption) (*ErrorCode, error)
	EditFreelancerProfile(ctx context.Context, in *EditFreelancerProfileRequest, opts ...grpc.CallOption) (*ErrorCode, error)
	GetEmployerProfile(ctx context.Context, in *ID, opts ...grpc.CallOption) (*EmployerProfile, error)
	GetFreelancerProfile(ctx context.Context, in *ID, opts ...grpc.CallOption) (*FreelancerProfile, error)
}

type profilesClient struct {
	cc grpc.ClientConnInterface
}

func NewProfilesClient(cc grpc.ClientConnInterface) ProfilesClient {
	return &profilesClient{cc}
}

func (c *profilesClient) EditEmployerProfile(ctx context.Context, in *EditEmployerProfileRequest, opts ...grpc.CallOption) (*ErrorCode, error) {
	out := new(ErrorCode)
	err := c.cc.Invoke(ctx, "/api.Profiles/EditEmployerProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profilesClient) EditFreelancerProfile(ctx context.Context, in *EditFreelancerProfileRequest, opts ...grpc.CallOption) (*ErrorCode, error) {
	out := new(ErrorCode)
	err := c.cc.Invoke(ctx, "/api.Profiles/EditFreelancerProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profilesClient) GetEmployerProfile(ctx context.Context, in *ID, opts ...grpc.CallOption) (*EmployerProfile, error) {
	out := new(EmployerProfile)
	err := c.cc.Invoke(ctx, "/api.Profiles/GetEmployerProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profilesClient) GetFreelancerProfile(ctx context.Context, in *ID, opts ...grpc.CallOption) (*FreelancerProfile, error) {
	out := new(FreelancerProfile)
	err := c.cc.Invoke(ctx, "/api.Profiles/GetFreelancerProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfilesServer is the server API for Profiles service.
type ProfilesServer interface {
	EditEmployerProfile(context.Context, *EditEmployerProfileRequest) (*ErrorCode, error)
	EditFreelancerProfile(context.Context, *EditFreelancerProfileRequest) (*ErrorCode, error)
	GetEmployerProfile(context.Context, *ID) (*EmployerProfile, error)
	GetFreelancerProfile(context.Context, *ID) (*FreelancerProfile, error)
}

// UnimplementedProfilesServer can be embedded to have forward compatible implementations.
type UnimplementedProfilesServer struct {
}

func (*UnimplementedProfilesServer) EditEmployerProfile(ctx context.Context, req *EditEmployerProfileRequest) (*ErrorCode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditEmployerProfile not implemented")
}
func (*UnimplementedProfilesServer) EditFreelancerProfile(ctx context.Context, req *EditFreelancerProfileRequest) (*ErrorCode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditFreelancerProfile not implemented")
}
func (*UnimplementedProfilesServer) GetEmployerProfile(ctx context.Context, req *ID) (*EmployerProfile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEmployerProfile not implemented")
}
func (*UnimplementedProfilesServer) GetFreelancerProfile(ctx context.Context, req *ID) (*FreelancerProfile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFreelancerProfile not implemented")
}

func RegisterProfilesServer(s *grpc.Server, srv ProfilesServer) {
	s.RegisterService(&_Profiles_serviceDesc, srv)
}

func _Profiles_EditEmployerProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditEmployerProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfilesServer).EditEmployerProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Profiles/EditEmployerProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfilesServer).EditEmployerProfile(ctx, req.(*EditEmployerProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profiles_EditFreelancerProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditFreelancerProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfilesServer).EditFreelancerProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Profiles/EditFreelancerProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfilesServer).EditFreelancerProfile(ctx, req.(*EditFreelancerProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profiles_GetEmployerProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfilesServer).GetEmployerProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Profiles/GetEmployerProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfilesServer).GetEmployerProfile(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Profiles_GetFreelancerProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfilesServer).GetFreelancerProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Profiles/GetFreelancerProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfilesServer).GetFreelancerProfile(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}

var _Profiles_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Profiles",
	HandlerType: (*ProfilesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EditEmployerProfile",
			Handler:    _Profiles_EditEmployerProfile_Handler,
		},
		{
			MethodName: "EditFreelancerProfile",
			Handler:    _Profiles_EditFreelancerProfile_Handler,
		},
		{
			MethodName: "GetEmployerProfile",
			Handler:    _Profiles_GetEmployerProfile_Handler,
		},
		{
			MethodName: "GetFreelancerProfile",
			Handler:    _Profiles_GetFreelancerProfile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}