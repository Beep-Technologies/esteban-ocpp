package schemas

// ClearChargingProfileRequest
type ClearChargingProfileRequest struct {
	ChargingProfilePurpose string `json:"chargingProfilePurpose,omitempty"`
	ConnectorId            int    `json:"connectorId,omitempty"`
	Id                     int    `json:"id,omitempty"`
	StackLevel             int    `json:"stackLevel,omitempty"`
}
