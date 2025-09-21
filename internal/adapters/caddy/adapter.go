package caddy

type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) Install() error         { return InstallCaddy() }
func (a *Adapter) CreateCaddyfile() error { return CreateCaddyfile() }
func (a *Adapter) AddService(serviceHost string, servicePort int) error {
	return AddService(serviceHost, servicePort)
}
