package schemas

// AuthorizeResponse
type AuthorizeResponse struct {
	IdTagInfo *IdTagInfo `json:"idTagInfo"`
}

// IdTagInfo
type IdTagInfo struct {
	ExpiryDate  string `json:"expiryDate,omitempty"`
	ParentIdTag string `json:"parentIdTag,omitempty"`
	Status      string `json:"status"`
}
