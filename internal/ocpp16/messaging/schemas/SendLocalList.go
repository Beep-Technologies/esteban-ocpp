package schemas

// LocalAuthorizationListItems
type LocalAuthorizationListItems struct {
	IdTag     string     `json:"idTag"`
	IdTagInfo *IdTagInfo `json:"idTagInfo,omitempty"`
}

// SendLocalListRequest
type SendLocalListRequest struct {
	ListVersion            int                            `json:"listVersion"`
	LocalAuthorizationList []*LocalAuthorizationListItems `json:"localAuthorizationList,omitempty"`
	UpdateType             string                         `json:"updateType"`
}
