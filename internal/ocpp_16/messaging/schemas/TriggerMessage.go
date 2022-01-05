package schemas

// TriggerMessageRequest
type TriggerMessageRequest struct {
	ConnectorId      int    `json:"connectorId,omitempty"`
	RequestedMessage string `json:"requestedMessage"`
}
