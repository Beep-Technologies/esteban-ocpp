package schemas

// ChangeAvailabilityRequest
type ChangeAvailabilityRequest struct {
	ConnectorId int    `json:"connectorId"`
	Type        string `json:"type"`
}
