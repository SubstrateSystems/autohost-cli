package docker

type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) Install() error             { return Install() }
func (a *Adapter) StopApp(app string) error   { return StopApp(app) }
func (a *Adapter) StartApp(app string) error  { return StartApp(app) }
func (a *Adapter) RemoveApp(app string) error { return RemoveApp(app) }
func (a *Adapter) GetAppStatus(app string) (string, error) {
	return GetAppStatus(app)
}
