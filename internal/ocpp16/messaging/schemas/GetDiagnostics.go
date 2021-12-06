package schemas

// GetDiagnosticsRequest
type GetDiagnosticsRequest struct {
	Location      string `json:"location"`
	Retries       int    `json:"retries,omitempty"`
	RetryInterval int    `json:"retryInterval,omitempty"`
	StartTime     string `json:"startTime,omitempty"`
	StopTime      string `json:"stopTime,omitempty"`
}
