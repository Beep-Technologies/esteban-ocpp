package schemas

// UpdateFirmwareRequest
type UpdateFirmwareRequest struct {
	Location      string `json:"location"`
	Retries       int    `json:"retries,omitempty"`
	RetrieveDate  string `json:"retrieveDate"`
	RetryInterval int    `json:"retryInterval,omitempty"`
}
