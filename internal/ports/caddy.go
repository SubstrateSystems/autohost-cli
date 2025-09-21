package ports

type Caddy interface {
	Install() error
	CreateCaddyfile() error
	AddService(serviceHost string, servicePort int) error
}
