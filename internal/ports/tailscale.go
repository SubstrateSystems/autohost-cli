package ports

type Tailscale interface {
	Install() error
	Login() error
	IP() (string, error)
}

type SplitDNS interface {
	Ensure(domain string, nameservers []string) error
}
