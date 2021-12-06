package schemas

// ChargingSchedule
type ChargingSchedule struct {
	ChargingRateUnit       string                         `json:"chargingRateUnit"`
	ChargingSchedulePeriod []*ChargingSchedulePeriodItems `json:"chargingSchedulePeriod"`
	Duration               int                            `json:"duration,omitempty"`
	MinChargingRate        float64                        `json:"minChargingRate,omitempty"`
	StartSchedule          string                         `json:"startSchedule,omitempty"`
}

// ChargingSchedulePeriodItems
type ChargingSchedulePeriodItems struct {
	Limit        float64 `json:"limit"`
	NumberPhases int     `json:"numberPhases,omitempty"`
	StartPeriod  int     `json:"startPeriod"`
}

// GetCompositeScheduleResponse
type GetCompositeScheduleResponse struct {
	ChargingSchedule *ChargingSchedule `json:"chargingSchedule,omitempty"`
	ConnectorId      int               `json:"connectorId,omitempty"`
	ScheduleStart    string            `json:"scheduleStart,omitempty"`
	Status           string            `json:"status"`
}
