package ports

type CoreDNS interface {
	InstallAndRun(bindIP string) (corefilePath string, err error)
	UpdateCorefile(subdomain string, appIP string) error
}
