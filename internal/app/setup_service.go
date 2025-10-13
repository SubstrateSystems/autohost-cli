package app

import (
	"autohost-cli/internal/ports"
	"autohost-cli/utils"
	"errors"
	"fmt"
)

type SetupService struct {
	Docker ports.Docker
}

func (s *SetupService) Setup() error {

	if !s.Docker.DockerInstalled() {
		if utils.Confirm("Docker no está instalado. ¿Quieres instalarlo ahora? [y/N]: ") {
			s.Docker.Install()
			s.Docker.CreateDockerNetwork()
			return errors.New("Docker instalado. Por favor, reinicia la terminal y ejecuta el comando de nuevo")

		} else {
			fmt.Println("🚫 Instalación cancelada. Instala Docker manualmente y vuelve a ejecutar el setup.")
		}

		if utils.Confirm("¿Deseas agregar tu usuario al grupo 'docker' para usar Docker sin sudo? [y/N]: ") {
			s.Docker.AddUserToDockerGroup()
		}

	} else {
		fmt.Println("✅ Docker ya está instalado.")
		s.Docker.CreateDockerNetwork()
	}
	return nil
}
