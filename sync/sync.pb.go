// Code generated by protoc-gen-go.
// source: sync.proto
// DO NOT EDIT!

/*
Package sync is a generated protocol buffer package.

It is generated from these files:
	sync.proto

It has these top-level messages:
	CacheEntry
	Cache
	Entry
	Change
	Conflict
	Problem
*/
package sync

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EntryKind int32

const (
	EntryKind_Directory EntryKind = 0
	EntryKind_File      EntryKind = 1
)

var EntryKind_name = map[int32]string{
	0: "Directory",
	1: "File",
}
var EntryKind_value = map[string]int32{
	"Directory": 0,
	"File":      1,
}

func (x EntryKind) String() string {
	return proto.EnumName(EntryKind_name, int32(x))
}
func (EntryKind) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type CacheEntry struct {
	Mode             uint32                     `protobuf:"varint,1,opt,name=mode" json:"mode,omitempty"`
	ModificationTime *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=modificationTime" json:"modificationTime,omitempty"`
	Size             uint64                     `protobuf:"varint,3,opt,name=size" json:"size,omitempty"`
	Digest           []byte                     `protobuf:"bytes,4,opt,name=digest,proto3" json:"digest,omitempty"`
}

func (m *CacheEntry) Reset()                    { *m = CacheEntry{} }
func (m *CacheEntry) String() string            { return proto.CompactTextString(m) }
func (*CacheEntry) ProtoMessage()               {}
func (*CacheEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *CacheEntry) GetModificationTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.ModificationTime
	}
	return nil
}

type Cache struct {
	Entries map[string]*CacheEntry `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Cache) Reset()                    { *m = Cache{} }
func (m *Cache) String() string            { return proto.CompactTextString(m) }
func (*Cache) ProtoMessage()               {}
func (*Cache) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Cache) GetEntries() map[string]*CacheEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

type Entry struct {
	Kind       EntryKind         `protobuf:"varint,1,opt,name=kind,enum=sync.EntryKind" json:"kind,omitempty"`
	Executable bool              `protobuf:"varint,2,opt,name=executable" json:"executable,omitempty"`
	Digest     []byte            `protobuf:"bytes,3,opt,name=digest,proto3" json:"digest,omitempty"`
	Contents   map[string]*Entry `protobuf:"bytes,4,rep,name=contents" json:"contents,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Entry) Reset()                    { *m = Entry{} }
func (m *Entry) String() string            { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()               {}
func (*Entry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Entry) GetContents() map[string]*Entry {
	if m != nil {
		return m.Contents
	}
	return nil
}

type Change struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Old  *Entry `protobuf:"bytes,2,opt,name=old" json:"old,omitempty"`
	New  *Entry `protobuf:"bytes,3,opt,name=new" json:"new,omitempty"`
}

func (m *Change) Reset()                    { *m = Change{} }
func (m *Change) String() string            { return proto.CompactTextString(m) }
func (*Change) ProtoMessage()               {}
func (*Change) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Change) GetOld() *Entry {
	if m != nil {
		return m.Old
	}
	return nil
}

func (m *Change) GetNew() *Entry {
	if m != nil {
		return m.New
	}
	return nil
}

type Conflict struct {
	Path         string    `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	AlphaChanges []*Change `protobuf:"bytes,2,rep,name=alphaChanges" json:"alphaChanges,omitempty"`
	BetaChanges  []*Change `protobuf:"bytes,3,rep,name=betaChanges" json:"betaChanges,omitempty"`
}

func (m *Conflict) Reset()                    { *m = Conflict{} }
func (m *Conflict) String() string            { return proto.CompactTextString(m) }
func (*Conflict) ProtoMessage()               {}
func (*Conflict) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Conflict) GetAlphaChanges() []*Change {
	if m != nil {
		return m.AlphaChanges
	}
	return nil
}

func (m *Conflict) GetBetaChanges() []*Change {
	if m != nil {
		return m.BetaChanges
	}
	return nil
}

type Problem struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	// TODO: Should we switch this to an enumeration? See how many error
	// conditions we run into while implementing transition methods.
	Error string `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
}

func (m *Problem) Reset()                    { *m = Problem{} }
func (m *Problem) String() string            { return proto.CompactTextString(m) }
func (*Problem) ProtoMessage()               {}
func (*Problem) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func init() {
	proto.RegisterType((*CacheEntry)(nil), "sync.CacheEntry")
	proto.RegisterType((*Cache)(nil), "sync.Cache")
	proto.RegisterType((*Entry)(nil), "sync.Entry")
	proto.RegisterType((*Change)(nil), "sync.Change")
	proto.RegisterType((*Conflict)(nil), "sync.Conflict")
	proto.RegisterType((*Problem)(nil), "sync.Problem")
	proto.RegisterEnum("sync.EntryKind", EntryKind_name, EntryKind_value)
}

func init() { proto.RegisterFile("sync.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 467 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x92, 0xcf, 0xaf, 0xd2, 0x40,
	0x10, 0xc7, 0x5d, 0x5a, 0x78, 0x30, 0x80, 0x36, 0x1b, 0x63, 0x2a, 0x89, 0x5a, 0xab, 0x31, 0x8d,
	0x87, 0x3e, 0xc3, 0x8b, 0x89, 0xf1, 0x8a, 0xef, 0xc5, 0x44, 0x0f, 0x66, 0xf3, 0x4e, 0xde, 0xfa,
	0x63, 0x80, 0xcd, 0x6b, 0x77, 0x49, 0xbb, 0x88, 0x78, 0xf2, 0xea, 0xdd, 0x7f, 0xd1, 0xff, 0xc3,
	0x74, 0x16, 0x48, 0xd1, 0xde, 0x66, 0xe7, 0x3b, 0xf3, 0xe1, 0xfb, 0x65, 0x0a, 0x50, 0xef, 0x55,
	0x16, 0x6f, 0x2a, 0x6d, 0x34, 0x77, 0x9b, 0x7a, 0xf6, 0x6c, 0xa5, 0xf5, 0xaa, 0xc0, 0x4b, 0xea,
	0xa5, 0xdb, 0xe5, 0xa5, 0x91, 0x25, 0xd6, 0x26, 0x29, 0x37, 0x76, 0x2c, 0xfc, 0xcd, 0x00, 0x16,
	0x49, 0xb6, 0xc6, 0x6b, 0x65, 0xaa, 0x3d, 0xe7, 0xe0, 0x96, 0x3a, 0x47, 0x9f, 0x05, 0x2c, 0x9a,
	0x0a, 0xaa, 0xf9, 0x0d, 0x78, 0xa5, 0xce, 0xe5, 0x52, 0x66, 0x89, 0x91, 0x5a, 0xdd, 0xca, 0x12,
	0xfd, 0x5e, 0xc0, 0xa2, 0xf1, 0x7c, 0x16, 0x5b, 0x7c, 0x7c, 0xc4, 0xc7, 0xb7, 0x47, 0xbc, 0xf8,
	0x6f, 0xa7, 0x61, 0xd7, 0xf2, 0x07, 0xfa, 0x4e, 0xc0, 0x22, 0x57, 0x50, 0xcd, 0x1f, 0xc1, 0x20,
	0x97, 0x2b, 0xac, 0x8d, 0xef, 0x06, 0x2c, 0x9a, 0x88, 0xc3, 0x2b, 0xfc, 0xc5, 0xa0, 0x4f, 0xb6,
	0xf8, 0x1c, 0x2e, 0x50, 0x99, 0x4a, 0x62, 0xed, 0xb3, 0xc0, 0x89, 0xc6, 0x73, 0x3f, 0xa6, 0x94,
	0xa4, 0xc6, 0xd7, 0x56, 0x22, 0xf3, 0xe2, 0x38, 0x38, 0xfb, 0x0c, 0x93, 0xb6, 0xc0, 0x3d, 0x70,
	0xee, 0x70, 0x4f, 0xa1, 0x46, 0xa2, 0x29, 0xf9, 0x2b, 0xe8, 0x7f, 0x4b, 0x8a, 0xed, 0x31, 0x88,
	0xd7, 0x62, 0x5a, 0x96, 0x95, 0xdf, 0xf7, 0xde, 0xb1, 0xf0, 0x0f, 0x83, 0xbe, 0xe5, 0xbc, 0x00,
	0xf7, 0x4e, 0xaa, 0x9c, 0x40, 0xf7, 0xe7, 0x0f, 0xec, 0x12, 0x49, 0x9f, 0xa4, 0xca, 0x05, 0x89,
	0xfc, 0x29, 0x00, 0x7e, 0xc7, 0x6c, 0x6b, 0x92, 0xb4, 0xb0, 0xfc, 0xa1, 0x68, 0x75, 0x5a, 0x91,
	0x9d, 0x76, 0x64, 0xfe, 0x16, 0x86, 0x99, 0x56, 0x06, 0x95, 0xa9, 0x7d, 0x97, 0x92, 0x3e, 0x6e,
	0xfd, 0x40, 0xbc, 0x38, 0x68, 0xd6, 0xde, 0x69, 0x74, 0xf6, 0x11, 0xa6, 0x67, 0x52, 0x47, 0xd8,
	0xe7, 0xe7, 0x61, 0xc7, 0x2d, 0x6c, 0x3b, 0xe7, 0x57, 0x18, 0x2c, 0xd6, 0x89, 0x5a, 0xd1, 0xa5,
	0x36, 0x89, 0x59, 0x1f, 0x18, 0x54, 0xf3, 0x27, 0xe0, 0xe8, 0x22, 0xef, 0x42, 0x34, 0xfd, 0x46,
	0x56, 0xb8, 0xa3, 0x48, 0xff, 0xca, 0x0a, 0x77, 0xe1, 0x4f, 0x06, 0xc3, 0x85, 0x56, 0xcb, 0x42,
	0x66, 0xa6, 0x13, 0xff, 0x06, 0x26, 0x49, 0xb1, 0x59, 0x27, 0xd6, 0x41, 0xed, 0xf7, 0xe8, 0x1f,
	0x98, 0x1c, 0xee, 0x42, 0x4d, 0x71, 0x36, 0xc1, 0x63, 0x18, 0xa7, 0x68, 0x4e, 0x0b, 0x4e, 0xc7,
	0x42, 0x7b, 0x20, 0xbc, 0x82, 0x8b, 0x2f, 0x95, 0x4e, 0x0b, 0x2c, 0x3b, 0x0d, 0x3c, 0x84, 0x3e,
	0x56, 0x95, 0xae, 0x28, 0xe1, 0x48, 0xd8, 0xc7, 0xeb, 0x97, 0x30, 0x3a, 0xdd, 0x97, 0x4f, 0x61,
	0xf4, 0x41, 0x56, 0x98, 0x19, 0x5d, 0xed, 0xbd, 0x7b, 0x7c, 0x08, 0xee, 0x8d, 0x2c, 0xd0, 0x63,
	0xe9, 0x80, 0xbe, 0xff, 0xab, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x91, 0x60, 0x3d, 0xab, 0x80,
	0x03, 0x00, 0x00,
}
