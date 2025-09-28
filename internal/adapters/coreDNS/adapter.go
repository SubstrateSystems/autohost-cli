package coredns

type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) InstallAndRun(tailscaleIP string) (string, error) {
	return InstallAndRunCoreDNS(tailscaleIP)
}

func (a *Adapter) UpdateCorefile(subdomain string, appIP string) error {
	return UpdateCorefile(subdomain, appIP)
}
