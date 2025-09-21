package cloudflare

// Funciones existentes esperadas:
// func InstallCloudflare() error
// func LoginCloudflare() error

type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) Install() error { return InstallCloudflare() }
func (a *Adapter) Login() error   { return LoginCloudflare() }
