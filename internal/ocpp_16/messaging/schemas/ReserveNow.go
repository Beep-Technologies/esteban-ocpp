package schemas

// ReserveNowRequest
type ReserveNowRequest struct {
	ConnectorId   int    `json:"connectorId"`
	ExpiryDate    string `json:"expiryDate"`
	IdTag         string `json:"idTag"`
	ParentIdTag   string `json:"parentIdTag,omitempty"`
	ReservationId int    `json:"reservationId"`
}
