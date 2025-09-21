package ports

type CoreDNS interface {
	InstallAndRun(bindIP string) (corefilePath string, err error)
}
