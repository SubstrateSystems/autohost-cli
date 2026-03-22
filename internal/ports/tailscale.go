package ports

type Tailscale interface {
	Installed() bool
	Install() error
	Login() error
	IP() (string, error)
	GetMachineName() (string, error)
}

type SplitDNS interface {
	Ensure(domain string, nameservers []string) error
}
