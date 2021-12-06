package schemas

// StartTransactionRequest
type StartTransactionRequest struct {
	ConnectorId   int    `json:"connectorId"`
	IdTag         string `json:"idTag"`
	MeterStart    int    `json:"meterStart"`
	ReservationId int    `json:"reservationId,omitempty"`
	Timestamp     string `json:"timestamp"`
}
