package types

type NodeData struct {
	HostName     string
	IPLocal      string
	OS           string
	Arch         string
	VersionAgent string
}

type NodeRquest struct {
	ErollToken   string `json:"enroll_token"`
	HostName     string `json:"hostname"`
	IPLocal      string `json:"ip_local"`
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	VersionAgent string `json:"version_agent"`
}

type NodeResponse struct {
	ApiToken string `json:"api_token"`
}
