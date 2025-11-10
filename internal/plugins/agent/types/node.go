package types

type NodeData struct {
	HostName     string
	IPLocal      string
	OS           string
	Arch         string
	VersionAgent string
}

type NodeRquest struct {
	HostName     string `json:"hostname"`
	IPLocal      string `json:"ip_local"`
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	VersionAgent string `json:"version_agent"`
}

type NodeResponse struct {
	ID string `json:"id"`
}
