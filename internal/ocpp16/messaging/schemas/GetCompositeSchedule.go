package schemas

// GetCompositeScheduleRequest
type GetCompositeScheduleRequest struct {
	ChargingRateUnit string `json:"chargingRateUnit,omitempty"`
	ConnectorId      int    `json:"connectorId"`
	Duration         int    `json:"duration"`
}
