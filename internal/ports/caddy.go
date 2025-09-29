package ports

import "context"

type Caddy interface {
	Install() error
	CreateCaddyfile() error
	AddService(serviceHost string, servicePort int) error

	EnsureCaddySnippetsSetup(ctx context.Context) error
}
