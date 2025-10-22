package cloudflare

type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) Install() error { return Install() }
func (a *Adapter) Login() error   { return Login() }
func (a *Adapter) Tunnel() error  { return Tunnel() }
