// Code generated by protoc-gen-go.
// source: charge_pile_status.proto
// DO NOT EDIT!

/*
Package Report is a generated protocol buffer package.

It is generated from these files:
	charge_pile_status.proto
	charge_station_status.proto

It has these top-level messages:
	ChargingPileStatus
	ChargingStationStatus
*/
package Report

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ChargingPileStatus_ChargingPileStatusType int32

const (
	ChargingPileStatus_IDLE         ChargingPileStatus_ChargingPileStatusType = 0
	ChargingPileStatus_CHARGING     ChargingPileStatus_ChargingPileStatusType = 1
	ChargingPileStatus_TOBECHARGING ChargingPileStatus_ChargingPileStatusType = 2
	ChargingPileStatus_FULL         ChargingPileStatus_ChargingPileStatusType = 3
	ChargingPileStatus_MAINTAINACE  ChargingPileStatus_ChargingPileStatusType = 4
	ChargingPileStatus_OFFLINE      ChargingPileStatus_ChargingPileStatusType = 5
)

var ChargingPileStatus_ChargingPileStatusType_name = map[int32]string{
	0: "IDLE",
	1: "CHARGING",
	2: "TOBECHARGING",
	3: "FULL",
	4: "MAINTAINACE",
	5: "OFFLINE",
}
var ChargingPileStatus_ChargingPileStatusType_value = map[string]int32{
	"IDLE":         0,
	"CHARGING":     1,
	"TOBECHARGING": 2,
	"FULL":         3,
	"MAINTAINACE":  4,
	"OFFLINE":      5,
}

func (x ChargingPileStatus_ChargingPileStatusType) String() string {
	return proto.EnumName(ChargingPileStatus_ChargingPileStatusType_name, int32(x))
}

type ChargingPileStatus struct {
	DasUuid          string                                    `protobuf:"bytes,1,opt,name=das_uuid" json:"das_uuid,omitempty"`
	Cpid             uint64                                    `protobuf:"varint,2,opt,name=cpid" json:"cpid,omitempty"`
	Status           ChargingPileStatus_ChargingPileStatusType `protobuf:"varint,3,opt,name=status,enum=Report.ChargingPileStatus_ChargingPileStatusType" json:"status,omitempty"`
	ChargingDuration uint32                                    `protobuf:"varint,4,opt,name=ChargingDuration" json:"ChargingDuration,omitempty"`
	ChargingCapacity uint32                                    `protobuf:"varint,5,opt,name=ChargingCapacity" json:"ChargingCapacity,omitempty"`
	ChargingPrice    uint32                                    `protobuf:"varint,6,opt,name=ChargingPrice" json:"ChargingPrice,omitempty"`
	// id
	Id uint32 `protobuf:"varint,7,opt,name=id" json:"id,omitempty"`
	// 所属充电站id
	StationId uint32 `protobuf:"varint,8,opt,name=stationId" json:"stationId,omitempty"`
	// 终端类型id
	TerminalTypeId uint32 `protobuf:"varint,9,opt,name=terminalTypeId" json:"terminalTypeId,omitempty"`
	// 额定功率 单位：KW
	RatedPower float64 `protobuf:"fixed64,10,opt,name=ratedPower" json:"ratedPower,omitempty"`
	// 电流类型 0：交流，1：直流
	ElectricCurrentType uint32 `protobuf:"varint,11,opt,name=electricCurrentType" json:"electricCurrentType,omitempty"`
	// 输入电压 单位：伏特
	VoltageInput uint32 `protobuf:"varint,12,opt,name=voltageInput" json:"voltageInput,omitempty"`
	// 输出电压 单位：伏特
	VoltageOutput uint32 `protobuf:"varint,13,opt,name=voltageOutput" json:"voltageOutput,omitempty"`
	// 输出电流 单位：安培
	ElectricCurrentOutput uint32 `protobuf:"varint,14,opt,name=electricCurrentOutput" json:"electricCurrentOutput,omitempty"`
	// 枪个数
	GunNumber uint32 `protobuf:"varint,15,opt,name=gunNumber" json:"gunNumber,omitempty"`
	// 电表读数
	AmmeterNumber float64 `protobuf:"fixed64,16,opt,name=ammeterNumber" json:"ammeterNumber,omitempty"`
	// 充电桩编码 车位号
	Code uint32 `protobuf:"varint,17,opt,name=code" json:"code,omitempty"`
	// 当前订单的订单号 30位（14位日期+16为cpid）
	CurrentOrderNumber string `protobuf:"bytes,18,opt,name=currentOrderNumber" json:"currentOrderNumber,omitempty"`
	// 接口 0:RS232,1:RS485,2:CAN,3:USB,4:RJ45,5:RS232(DEBUG)
	InterfaceType uint32 `protobuf:"varint,19,opt,name=interfaceType" json:"interfaceType,omitempty"`
	// 波特率 0:9600,1:14400,2:19200,3:38400,4:576005,5:115200
	BaudRate  uint32 `protobuf:"varint,20,opt,name=baudRate" json:"baudRate,omitempty"`
	Timestamp uint64 `protobuf:"varint,21,opt,name=Timestamp" json:"Timestamp,omitempty"`
}

func (m *ChargingPileStatus) Reset()         { *m = ChargingPileStatus{} }
func (m *ChargingPileStatus) String() string { return proto.CompactTextString(m) }
func (*ChargingPileStatus) ProtoMessage()    {}

func init() {
	proto.RegisterEnum("Report.ChargingPileStatus_ChargingPileStatusType", ChargingPileStatus_ChargingPileStatusType_name, ChargingPileStatus_ChargingPileStatusType_value)
}
