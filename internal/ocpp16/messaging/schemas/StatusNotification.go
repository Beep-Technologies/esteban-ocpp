package schemas

// StatusNotificationRequest
type StatusNotificationRequest struct {
	ConnectorId     int    `json:"connectorId"`
	ErrorCode       string `json:"errorCode"`
	Info            string `json:"info,omitempty"`
	Status          string `json:"status"`
	Timestamp       string `json:"timestamp,omitempty"`
	VendorErrorCode string `json:"vendorErrorCode,omitempty"`
	VendorId        string `json:"vendorId,omitempty"`
}
