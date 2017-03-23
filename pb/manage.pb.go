// Code generated by protoc-gen-go.
// source: manage.proto
// DO NOT EDIT!

package Report

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Command_CommandType int32

const (
	Command_CMT_INVALID Command_CommandType = 0
	// das<->tms
	Command_CMT_REQ_LOGIN Command_CommandType = 1
	Command_CMT_REP_LOGIN Command_CommandType = 32769
	// das<->tms
	Command_CMT_REQ_SETTING Command_CommandType = 16
	Command_CMT_REP_SETTING Command_CommandType = 32784
	// das->tss
	Command_CMT_REQ_HEART Command_CommandType = 5
	// das<->tms
	Command_CMT_REQ_PRICE Command_CommandType = 3
	Command_CMT_REP_PRICE Command_CommandType = 32771
	// wechat <-> das
	Command_CMT_REQ_GET_GUN_STATUS Command_CommandType = 32774
	Command_CMT_REP_GET_GUN_STATUS Command_CommandType = 6
	// wechat <-> das
	Command_CMT_REQ_CHARGING Command_CommandType = 32775
	Command_CMT_REP_CHARGING Command_CommandType = 7
	// wechat <-> das
	Command_CMT_REQ_STOP_CHARGING Command_CommandType = 32782
	Command_CMT_REP_STOP_CHARGING Command_CommandType = 14
	// web<->das
	Command_CMT_REQ_NOTIFY_SET_PRICE Command_CommandType = 32783
	Command_CMT_REP_NOTIFY_SET_PRICE Command_CommandType = 15
	// das->wechat
	Command_CMT_REP_CHARGING_STARTED Command_CommandType = 8
	// das->wechat
	Command_CMT_REP_CHARGING_DATA_UPLOAD Command_CommandType = 9
	Command_CMT_REP_CHARGING_COST        Command_CommandType = 32777
	// das->wechat
	Command_CMT_REP_CHARGING_STOPPED Command_CommandType = 11
	// das <-> web
	Command_CMT_REQ_PIN Command_CommandType = 32780
	Command_CMT_REP_PIN Command_CommandType = 12
	// das -> ?
	Command_CMT_REP_OFFLINE_DATA Command_CommandType = 13
	// tss->wechat
	Command_CMT_NOTIFY_TRANSCATION Command_CommandType = 256
)

var Command_CommandType_name = map[int32]string{
	0:     "CMT_INVALID",
	1:     "CMT_REQ_LOGIN",
	32769: "CMT_REP_LOGIN",
	16:    "CMT_REQ_SETTING",
	32784: "CMT_REP_SETTING",
	5:     "CMT_REQ_HEART",
	3:     "CMT_REQ_PRICE",
	32771: "CMT_REP_PRICE",
	32774: "CMT_REQ_GET_GUN_STATUS",
	6:     "CMT_REP_GET_GUN_STATUS",
	32775: "CMT_REQ_CHARGING",
	7:     "CMT_REP_CHARGING",
	32782: "CMT_REQ_STOP_CHARGING",
	14:    "CMT_REP_STOP_CHARGING",
	32783: "CMT_REQ_NOTIFY_SET_PRICE",
	15:    "CMT_REP_NOTIFY_SET_PRICE",
	8:     "CMT_REP_CHARGING_STARTED",
	9:     "CMT_REP_CHARGING_DATA_UPLOAD",
	32777: "CMT_REP_CHARGING_COST",
	11:    "CMT_REP_CHARGING_STOPPED",
	32780: "CMT_REQ_PIN",
	12:    "CMT_REP_PIN",
	13:    "CMT_REP_OFFLINE_DATA",
	256:   "CMT_NOTIFY_TRANSCATION",
}
var Command_CommandType_value = map[string]int32{
	"CMT_INVALID":                  0,
	"CMT_REQ_LOGIN":                1,
	"CMT_REP_LOGIN":                32769,
	"CMT_REQ_SETTING":              16,
	"CMT_REP_SETTING":              32784,
	"CMT_REQ_HEART":                5,
	"CMT_REQ_PRICE":                3,
	"CMT_REP_PRICE":                32771,
	"CMT_REQ_GET_GUN_STATUS":       32774,
	"CMT_REP_GET_GUN_STATUS":       6,
	"CMT_REQ_CHARGING":             32775,
	"CMT_REP_CHARGING":             7,
	"CMT_REQ_STOP_CHARGING":        32782,
	"CMT_REP_STOP_CHARGING":        14,
	"CMT_REQ_NOTIFY_SET_PRICE":     32783,
	"CMT_REP_NOTIFY_SET_PRICE":     15,
	"CMT_REP_CHARGING_STARTED":     8,
	"CMT_REP_CHARGING_DATA_UPLOAD": 9,
	"CMT_REP_CHARGING_COST":        32777,
	"CMT_REP_CHARGING_STOPPED":     11,
	"CMT_REQ_PIN":                  32780,
	"CMT_REP_PIN":                  12,
	"CMT_REP_OFFLINE_DATA":         13,
	"CMT_NOTIFY_TRANSCATION":       256,
}

func (x Command_CommandType) String() string {
	return proto.EnumName(Command_CommandType_name, int32(x))
}
func (Command_CommandType) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0, 0} }

type Command struct {
	Type         Command_CommandType `protobuf:"varint,1,opt,name=type,enum=Report.Command_CommandType" json:"type,omitempty"`
	Uuid         string              `protobuf:"bytes,2,opt,name=uuid" json:"uuid,omitempty"`
	Tid          uint64              `protobuf:"varint,3,opt,name=tid" json:"tid,omitempty"`
	SerialNumber uint32              `protobuf:"varint,4,opt,name=serial_number" json:"serial_number,omitempty"`
	Paras        []*Param            `protobuf:"bytes,5,rep,name=paras" json:"paras,omitempty"`
}

func (m *Command) Reset()                    { *m = Command{} }
func (m *Command) String() string            { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()               {}
func (*Command) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *Command) GetType() Command_CommandType {
	if m != nil {
		return m.Type
	}
	return Command_CMT_INVALID
}

func (m *Command) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Command) GetTid() uint64 {
	if m != nil {
		return m.Tid
	}
	return 0
}

func (m *Command) GetSerialNumber() uint32 {
	if m != nil {
		return m.SerialNumber
	}
	return 0
}

func (m *Command) GetParas() []*Param {
	if m != nil {
		return m.Paras
	}
	return nil
}

func init() {
	proto.RegisterType((*Command)(nil), "Report.Command")
	proto.RegisterEnum("Report.Command_CommandType", Command_CommandType_name, Command_CommandType_value)
}

func init() { proto.RegisterFile("manage.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 458 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x92, 0xc1, 0x6e, 0x9b, 0x4e,
	0x10, 0xc6, 0xff, 0xac, 0xb1, 0xf3, 0xcf, 0x62, 0xe2, 0xcd, 0x24, 0x8e, 0x68, 0x6c, 0x55, 0x28,
	0x27, 0x7a, 0xf1, 0x21, 0x7d, 0x82, 0x15, 0x60, 0xb2, 0x92, 0xbb, 0x6c, 0x61, 0x5d, 0xa9, 0x27,
	0x44, 0x64, 0x54, 0x59, 0x2a, 0x36, 0x22, 0xf6, 0x21, 0x37, 0xaa, 0x4a, 0xad, 0x2a, 0x55, 0x6d,
	0x5f, 0xb1, 0xb7, 0x3e, 0x46, 0x05, 0x78, 0x13, 0x9a, 0x9c, 0x90, 0x7e, 0xdf, 0xcc, 0x37, 0xdf,
	0x30, 0x8b, 0x87, 0x79, 0xba, 0x49, 0x3f, 0x64, 0xb3, 0xa2, 0xdc, 0xee, 0xb6, 0x30, 0x88, 0xb2,
	0x62, 0x5b, 0xee, 0x2e, 0x8d, 0x22, 0x2d, 0xd3, 0xbc, 0x85, 0x57, 0x7f, 0xfa, 0xf8, 0xc8, 0xdd,
	0xe6, 0x79, 0xba, 0x59, 0xc1, 0x2b, 0xac, 0xef, 0xee, 0x8b, 0xcc, 0xd2, 0x6c, 0xcd, 0x39, 0xb9,
	0x9e, 0xcc, 0xda, 0xfa, 0xd9, 0x41, 0x56, 0x5f, 0x79, 0x5f, 0x64, 0x30, 0xc4, 0xfa, 0x7e, 0xbf,
	0x5e, 0x59, 0xc8, 0xd6, 0x9c, 0x63, 0x30, 0x70, 0x6f, 0xb7, 0x5e, 0x59, 0x3d, 0x5b, 0x73, 0x74,
	0x18, 0x63, 0xf3, 0x2e, 0x2b, 0xd7, 0xe9, 0xc7, 0x64, 0xb3, 0xcf, 0x6f, 0xb3, 0xd2, 0xd2, 0x6d,
	0xcd, 0x31, 0x61, 0x8a, 0xfb, 0xf5, 0xdc, 0x3b, 0xab, 0x6f, 0xf7, 0x1c, 0xe3, 0xda, 0x54, 0xee,
	0xa2, 0x0e, 0x73, 0xf5, 0x5b, 0xc7, 0x46, 0xd7, 0x7f, 0x84, 0x0d, 0xf7, 0x8d, 0x4c, 0x18, 0x7f,
	0x47, 0x17, 0xcc, 0x23, 0xff, 0xc1, 0x29, 0x36, 0x6b, 0x10, 0xf9, 0x6f, 0x93, 0x45, 0x18, 0x30,
	0x4e, 0x34, 0x38, 0x53, 0x48, 0x1c, 0xd0, 0xa7, 0x0a, 0xc1, 0x19, 0x1e, 0xa9, 0xba, 0xd8, 0x97,
	0x92, 0xf1, 0x80, 0x10, 0x18, 0x2b, 0x28, 0x1e, 0xe0, 0xaf, 0x0a, 0x75, 0x3d, 0x6f, 0x7c, 0x1a,
	0x49, 0xd2, 0xef, 0x22, 0x11, 0x31, 0xd7, 0x27, 0xbd, 0xee, 0x98, 0x16, 0x7d, 0xae, 0x10, 0x4c,
	0xf1, 0x85, 0xaa, 0x0b, 0x7c, 0x99, 0x04, 0x4b, 0x9e, 0xc4, 0x92, 0xca, 0x65, 0x4c, 0xbe, 0x54,
	0x08, 0x2e, 0x95, 0x2a, 0x9e, 0xaa, 0x03, 0xb8, 0xc0, 0x44, 0x75, 0xba, 0x37, 0x34, 0x0a, 0xea,
	0x30, 0x5f, 0x2b, 0x04, 0xe7, 0x8a, 0x8b, 0x47, 0x7e, 0x04, 0x13, 0x3c, 0x7e, 0x58, 0x47, 0x86,
	0x1d, 0xe9, 0x47, 0x85, 0xe0, 0x85, 0x12, 0xc5, 0x13, 0xf1, 0x04, 0x5e, 0x62, 0x4b, 0xf5, 0xf1,
	0x50, 0xb2, 0xf9, 0xfb, 0x7a, 0xf1, 0x43, 0xfe, 0x9f, 0x4d, 0x7e, 0x4b, 0xb5, 0x3e, 0xd3, 0x47,
	0x5d, 0x55, 0x79, 0xd6, 0x0b, 0x44, 0xd2, 0xf7, 0xc8, 0xff, 0x60, 0xe3, 0xe9, 0x33, 0xd5, 0xa3,
	0x92, 0x26, 0x4b, 0xb1, 0x08, 0xa9, 0x47, 0x8e, 0x1f, 0x53, 0x77, 0x2a, 0xdc, 0x30, 0x96, 0xe4,
	0xdb, 0xbf, 0xa3, 0x3b, 0xe6, 0xa1, 0x10, 0xbe, 0x47, 0x0c, 0x38, 0x6d, 0x0f, 0xdf, 0x1c, 0x80,
	0x71, 0xf2, 0xbd, 0x42, 0xea, 0x2d, 0x34, 0x07, 0x60, 0x9c, 0x0c, 0xc1, 0xc2, 0xe7, 0x0a, 0x84,
	0xf3, 0xf9, 0x82, 0x71, 0xbf, 0x99, 0x4f, 0x4c, 0x98, 0xb4, 0x3f, 0xfe, 0xb0, 0x92, 0x8c, 0x28,
	0x8f, 0x5d, 0x2a, 0x59, 0xc8, 0x49, 0x85, 0x6e, 0x07, 0xcd, 0x8b, 0x7f, 0xfd, 0x37, 0x00, 0x00,
	0xff, 0xff, 0xd6, 0x52, 0xe4, 0x8b, 0x16, 0x03, 0x00, 0x00,
}