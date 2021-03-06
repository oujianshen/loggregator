// Code generated by protoc-gen-gogo.
// source: log_message.proto
// DO NOT EDIT!

package logMessage

import proto "code.google.com/p/gogoprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type LogMessage_MessageType int32

const (
	LogMessage_OUT LogMessage_MessageType = 1
	LogMessage_ERR LogMessage_MessageType = 2
)

var LogMessage_MessageType_name = map[int32]string{
	1: "OUT",
	2: "ERR",
}
var LogMessage_MessageType_value = map[string]int32{
	"OUT": 1,
	"ERR": 2,
}

func (x LogMessage_MessageType) Enum() *LogMessage_MessageType {
	p := new(LogMessage_MessageType)
	*p = x
	return p
}
func (x LogMessage_MessageType) String() string {
	return proto.EnumName(LogMessage_MessageType_name, int32(x))
}
func (x LogMessage_MessageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}
func (x *LogMessage_MessageType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(LogMessage_MessageType_value, data, "LogMessage_MessageType")
	if err != nil {
		return err
	}
	*x = LogMessage_MessageType(value)
	return nil
}

type LogMessage_SourceType int32

const (
	LogMessage_CLOUD_CONTROLLER LogMessage_SourceType = 1
	LogMessage_ROUTER           LogMessage_SourceType = 2
	LogMessage_UAA              LogMessage_SourceType = 3
	LogMessage_DEA              LogMessage_SourceType = 4
	LogMessage_WARDEN_CONTAINER LogMessage_SourceType = 5
)

var LogMessage_SourceType_name = map[int32]string{
	1: "CLOUD_CONTROLLER",
	2: "ROUTER",
	3: "UAA",
	4: "DEA",
	5: "WARDEN_CONTAINER",
}
var LogMessage_SourceType_value = map[string]int32{
	"CLOUD_CONTROLLER": 1,
	"ROUTER":           2,
	"UAA":              3,
	"DEA":              4,
	"WARDEN_CONTAINER": 5,
}

func (x LogMessage_SourceType) Enum() *LogMessage_SourceType {
	p := new(LogMessage_SourceType)
	*p = x
	return p
}
func (x LogMessage_SourceType) String() string {
	return proto.EnumName(LogMessage_SourceType_name, int32(x))
}
func (x LogMessage_SourceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}
func (x *LogMessage_SourceType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(LogMessage_SourceType_value, data, "LogMessage_SourceType")
	if err != nil {
		return err
	}
	*x = LogMessage_SourceType(value)
	return nil
}

type LogMessage struct {
	Message          []byte                  `protobuf:"bytes,1,req,name=message" json:"message,omitempty"`
	MessageType      *LogMessage_MessageType `protobuf:"varint,2,req,name=message_type,enum=logMessage.LogMessage_MessageType" json:"message_type,omitempty"`
	Timestamp        *int64                  `protobuf:"zigzag64,3,req,name=timestamp" json:"timestamp,omitempty"`
	AppId            *string                 `protobuf:"bytes,4,req,name=app_id" json:"app_id,omitempty"`
	SourceType       *LogMessage_SourceType  `protobuf:"varint,5,req,name=source_type,enum=logMessage.LogMessage_SourceType" json:"source_type,omitempty"`
	SourceId         *string                 `protobuf:"bytes,6,opt,name=source_id" json:"source_id,omitempty"`
	SpaceId          *string                 `protobuf:"bytes,7,req,name=space_id" json:"space_id,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (m *LogMessage) Reset()         { *m = LogMessage{} }
func (m *LogMessage) String() string { return proto.CompactTextString(m) }
func (*LogMessage) ProtoMessage()    {}

func (m *LogMessage) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *LogMessage) GetMessageType() LogMessage_MessageType {
	if m != nil && m.MessageType != nil {
		return *m.MessageType
	}
	return 0
}

func (m *LogMessage) GetTimestamp() int64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *LogMessage) GetAppId() string {
	if m != nil && m.AppId != nil {
		return *m.AppId
	}
	return ""
}

func (m *LogMessage) GetSourceType() LogMessage_SourceType {
	if m != nil && m.SourceType != nil {
		return *m.SourceType
	}
	return 0
}

func (m *LogMessage) GetSourceId() string {
	if m != nil && m.SourceId != nil {
		return *m.SourceId
	}
	return ""
}

func (m *LogMessage) GetSpaceId() string {
	if m != nil && m.SpaceId != nil {
		return *m.SpaceId
	}
	return ""
}

func init() {
	proto.RegisterEnum("logMessage.LogMessage_MessageType", LogMessage_MessageType_name, LogMessage_MessageType_value)
	proto.RegisterEnum("logMessage.LogMessage_SourceType", LogMessage_SourceType_name, LogMessage_SourceType_value)
}
