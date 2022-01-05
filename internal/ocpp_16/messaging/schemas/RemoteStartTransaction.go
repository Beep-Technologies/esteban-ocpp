package schemas

// ChargingProfile
type ChargingProfile struct {
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

// RemoteStartTransactionRequest
type RemoteStartTransactionRequest struct {
	ChargingProfile *ChargingProfile `json:"chargingProfile,omitempty"`
	ConnectorId     int              `json:"connectorId,omitempty"`
	IdTag           string           `json:"idTag"`
}
