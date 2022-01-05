package schemas

// MeterValueItems
type MeterValueItems struct {
	SampledValue []*SampledValueItems `json:"sampledValue"`
	Timestamp    string               `json:"timestamp"`
}

// MeterValuesRequest
type MeterValuesRequest struct {
	ConnectorId   int                `json:"connectorId"`
	MeterValue    []*MeterValueItems `json:"meterValue"`
	TransactionId int                `json:"transactionId,omitempty"`
}

// SampledValueItems
type SampledValueItems struct {
	Context   string `json:"context,omitempty"`
	Format    string `json:"format,omitempty"`
	Location  string `json:"location,omitempty"`
	Measurand string `json:"measurand,omitempty"`
	Phase     string `json:"phase,omitempty"`
	Unit      string `json:"unit,omitempty"`
	Value     string `json:"value"`
}
