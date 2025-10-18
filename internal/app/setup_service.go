package app

import (
	"autohost-cli/internal/adapters/infra"
	"autohost-cli/internal/ports"
	"autohost-cli/utils"
	"fmt"
)

type SetupService struct {
	Docker ports.Docker
}

func (s *SetupService) Setup() error {
	// 1) Verificar Docker
	if !s.Docker.DockerInstalled() {
		// Preguntar instalación
		if !utils.Confirm("Docker no está instalado. ¿Quieres instalarlo ahora? [y/N]: ") {
			return fmt.Errorf("🚫 instalación cancelada por el usuario: Docker es requerido para continuar")
		}

		// Instalar Docker
		if err := infra.RunStep("Instalación de Docker", s.Docker.Install); err != nil {
			return err // ya viene envuelto con contexto y emoji
		}
		// Ofrecer agregar al grupo docker
		if utils.Confirm("¿Deseas agregar tu usuario al grupo 'docker' para usar Docker sin sudo? [y/N]: ") {
			if err := infra.RunStep("Agregar usuario al grupo 'docker'", s.Docker.AddUserToDockerGroup); err != nil {
				return err
			}
			// Nota: newgrp solo afecta a shells interactivos; aquí mejor avisar
			fmt.Println("ℹ️  Cierra sesión/reinicia la terminal para aplicar los cambios del grupo 'docker'.")
		}

		// Crear red de Docker
		if err := infra.RunStep("Creación de red de Docker", s.Docker.CreateDockerNetwork); err != nil {
			return err
		}

		// Todo bien
		fmt.Println("✅ Docker instalado y configurado.")
		return nil
	}

	// 2) Docker ya instalado: crear/red validar red
	fmt.Println("✅ Docker ya está instalado.")
	if err := infra.RunStep("Creación de red de Docker", s.Docker.CreateDockerNetwork); err != nil {
		return err
	}

	return nil
}
