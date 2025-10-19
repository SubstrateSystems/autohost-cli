package cloudflare

type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) Install() error { return InstallCloudflare() }
func (a *Adapter) Login() error   { return LoginCloudflare() }
