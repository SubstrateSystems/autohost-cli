package tailscale

type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) Install() error      { return InstallTailscale() }
func (a *Adapter) Login() error        { return LoginTailscale() }
func (a *Adapter) IP() (string, error) { return TailscaleIP() }
