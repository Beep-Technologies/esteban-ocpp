package schemas

// ConfigurationKeyItems
type ConfigurationKeyItems struct {
	Key      string `json:"key"`
	Readonly bool   `json:"readonly"`
	Value    string `json:"value,omitempty"`
}

// GetConfigurationResponse
type GetConfigurationResponse struct {
	ConfigurationKey []*ConfigurationKeyItems `json:"configurationKey,omitempty"`
	UnknownKey       []string                 `json:"unknownKey,omitempty"`
}
