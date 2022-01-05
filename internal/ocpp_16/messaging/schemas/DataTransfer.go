package schemas

// DataTransferRequest
type DataTransferRequest struct {
	Data      string `json:"data,omitempty"`
	MessageId string `json:"messageId,omitempty"`
	VendorId  string `json:"vendorId"`
}
