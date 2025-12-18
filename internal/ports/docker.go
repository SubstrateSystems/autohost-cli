package ports

import "autohost-cli/internal/domain"

type Docker interface {
	Install() error
	StopApp(app string) error
	StartApp(app string) error
	RemoveApp(app domain.AppName) error
	GetAppStatus(app string) (string, error)
	DockerInstalled() bool
	CreateDockerNetwork() error
	AddUserToDockerGroup() error
}
