package domain

// NodeData holds information gathered from the local machine before enrollment.
type NodeData struct {
	HostName     string
	IPLocal      string
	OS           string
	Arch         string
	VersionAgent string
}

// NodeRequest is the payload sent to the AutoHost API to enroll this node.
type NodeRequest struct {
	EnrollToken  string `json:"enroll_token"`
	HostName     string `json:"hostname"`
	IPLocal      string `json:"ip_local"`
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	VersionAgent string `json:"version_agent"`
}

// NodeResponse is the API response after a successful enrollment.
type NodeResponse struct {
	NodeID       string `json:"node_id"`
	ApiToken     string `json:"api_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
