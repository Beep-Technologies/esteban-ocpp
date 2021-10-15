package main

import (
	"fmt"
	"time"

	"github.com/Beep-Technologies/beepbeep3-ocpp/pkg/util"
)

type TransactionRPC struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id"`
	ChargePointId        int32    `protobuf:"varint,2,opt,name=charge_point_id,json=chargePointId,proto3" json:"charge_point_id"`
	ConnectorId          int32    `protobuf:"varint,3,opt,name=connector_id,json=connectorId,proto3" json:"connector_id"`
	IdTag                string   `protobuf:"bytes,4,opt,name=id_tag,json=idTag,proto3" json:"id_tag"`
	State                string   `protobuf:"bytes,5,opt,name=state,proto3" json:"state"`
	RemoteInitiated      bool     `protobuf:"varint,6,opt,name=remote_initiated,json=remoteInitiated,proto3" json:"remote_initiated"`
	StartTimestamp       string   `protobuf:"bytes,7,opt,name=start_timestamp,json=startTimestamp,proto3" json:"start_timestamp"`
	StopTimestamp        string   `protobuf:"bytes,8,opt,name=stop_timestamp,json=stopTimestamp,proto3" json:"stop_timestamp"`
	StartMeterValue      int32    `protobuf:"varint,9,opt,name=start_meter_value,json=startMeterValue,proto3" json:"start_meter_value"`
	StopMeterValue       int32    `protobuf:"varint,10,opt,name=stop_meter_value,json=stopMeterValue,proto3" json:"stop_meter_value"`
	StopReason           string   `protobuf:"bytes,11,opt,name=stop_reason,json=stopReason,proto3" json:"stop_reason"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type TransactionModel struct {
	ID              int32     `gorm:"column:id"`
	ChargePointID   int32     `gorm:"column:charge_point_id"`
	ConnectorID     int32     `gorm:"column:connector_id"`
	IDTag           string    `gorm:"column:id_tag"`
	State           string    `gorm:"column:state"`
	RemoteInitiated bool      `gorm:"column:remote_initiated"`
	StartTimestamp  time.Time `gorm:"column:start_timestamp"`
	StopTimestamp   time.Time `gorm:"column:stop_timestamp"`
	StartMeterValue int32     `gorm:"column:start_meter_value"`
	StopMeterValue  int32     `gorm:"column:stop_meter_value"`
	StopReason      string    `gorm:"column:stop_reason"`
}

func main() {
	a := TransactionRPC{
		Id:              21,
		ChargePointId:   21,
		ConnectorId:     21,
		IdTag:           "savage",
		State:           "INITIATED",
		RemoteInitiated: true,
		StartTimestamp:  "2006-01-02T15:04:05.000+07:00",
		StartMeterValue: 21,
		StopMeterValue:  21,
		StopReason:      "twentyone",
	}

	b := TransactionModel{}

	converters := make(map[string]util.ConverterFunc)
	converters["StartTimeStamp"] = util.ConvertRFC3339MilliToTime

	err := util.ConvertCopyStruct(&b, &a, converters)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n%+v\n", a, b)
}
