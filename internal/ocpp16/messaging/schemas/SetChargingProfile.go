package schemas

// CsChargingProfiles
type CsChargingProfiles struct {
	ChargingProfileId      int               `json:"chargingProfileId"`
	ChargingProfileKind    string            `json:"chargingProfileKind"`
	ChargingProfilePurpose string            `json:"chargingProfilePurpose"`
	ChargingSchedule       *ChargingSchedule `json:"chargingSchedule"`
	RecurrencyKind         string            `json:"recurrencyKind,omitempty"`
	StackLevel             int               `json:"stackLevel"`
	TransactionId          int               `json:"transactionId,omitempty"`
	ValidFrom              string            `json:"validFrom,omitempty"`
	ValidTo                string            `json:"validTo,omitempty"`
}

// SetChargingProfileRequest
type SetChargingProfileRequest struct {
	ConnectorId        int                 `json:"connectorId"`
	CsChargingProfiles *CsChargingProfiles `json:"csChargingProfiles"`
}
