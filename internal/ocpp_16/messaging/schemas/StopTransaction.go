package schemas

// StopTransactionRequest
type StopTransactionRequest struct {
	IdTag           string                  `json:"idTag,omitempty"`
	MeterStop       int                     `json:"meterStop"`
	Reason          string                  `json:"reason,omitempty"`
	Timestamp       string                  `json:"timestamp"`
	TransactionData []*TransactionDataItems `json:"transactionData,omitempty"`
	TransactionId   int                     `json:"transactionId"`
}

// TransactionDataItems
type TransactionDataItems struct {
	SampledValue []*SampledValueItems `json:"sampledValue"`
	Timestamp    string               `json:"timestamp"`
}
