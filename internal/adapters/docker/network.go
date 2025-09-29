package docker

import (
	"fmt"
	"os/exec"
)

func CreateDockerNetwork() error {
	// Verificar si la red ya existe
	cmd := exec.Command("docker", "network", "inspect", "autohost_net")
	if err := cmd.Run(); err == nil {
		fmt.Println("✅ La red 'autohost_net' ya existe.")
		return nil
	}
	// Crear la red
	cmd = exec.Command("docker", "network", "create", "autohost_net")
	return cmd.Run()
}
