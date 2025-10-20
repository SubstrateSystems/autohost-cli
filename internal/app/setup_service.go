package app

import (
	"autohost-cli/internal/adapters/infra"
	"autohost-cli/internal/ports"
	"autohost-cli/utils"
	"fmt"
	"os"
	"strings"
)

type SetupService struct {
	Docker ports.Docker
}

func (s *SetupService) Setup() error {
	// helper para respuestas automáticas
	autoYes := os.Getenv("CI") == "true" || strings.EqualFold(os.Getenv("ASSUME_YES"), "-y")

	// 1) Verificar Docker
	if !s.Docker.DockerInstalled() {
		confirm := autoYes || utils.Confirm("Docker no está instalado. ¿Quieres instalarlo ahora? [y/N]: ")
		if !confirm {
			return fmt.Errorf("🚫 instalación cancelada por el usuario: Docker es requerido para continuar")
		}

		if err := infra.RunStep("Instalación de Docker", s.Docker.Install); err != nil {
			return err
		}

		addGroup := autoYes || utils.Confirm("¿Deseas agregar tu usuario al grupo 'docker'? [y/N]: ")
		if addGroup {
			if err := infra.RunStep("Agregar usuario al grupo 'docker'", s.Docker.AddUserToDockerGroup); err != nil {
				return err
			}
			fmt.Println("ℹ️ Cierra sesión/reinicia la terminal para aplicar los cambios del grupo 'docker'.")
		}

		if err := infra.RunStep("Creación de red de Docker", s.Docker.CreateDockerNetwork); err != nil {
			return err
		}

		fmt.Println("✅ Docker instalado y configurado.")
		return nil
	}

	// 2) Docker ya instalado: solo validar red
	fmt.Println("✅ Docker ya está instalado.")
	if err := infra.RunStep("Creación de red de Docker", s.Docker.CreateDockerNetwork); err != nil {
		return err
	}
	return nil
}
