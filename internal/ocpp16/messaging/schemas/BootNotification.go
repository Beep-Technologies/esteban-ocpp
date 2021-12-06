package schemas

// BootNotificationRequest
type BootNotificationRequest struct {
	ChargeBoxSerialNumber   string `json:"chargeBoxSerialNumber,omitempty"`
	ChargePointModel        string `json:"chargePointModel"`
	ChargePointSerialNumber string `json:"chargePointSerialNumber,omitempty"`
	ChargePointVendor       string `json:"chargePointVendor"`
	FirmwareVersion         string `json:"firmwareVersion,omitempty"`
	Iccid                   string `json:"iccid,omitempty"`
	Imsi                    string `json:"imsi,omitempty"`
	MeterSerialNumber       string `json:"meterSerialNumber,omitempty"`
	MeterType               string `json:"meterType,omitempty"`
}
