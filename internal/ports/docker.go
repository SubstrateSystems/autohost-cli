package ports

type Docker interface {
	Install() error
	StopApp(app string) error
	StartApp(app string) error
	RemoveApp(app string) error
	GetAppStatus(app string) (string, error)
}
