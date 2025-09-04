package docker

import (
	"autohost-cli/utils"
	"fmt"
	"os"
	"os/user"
)

func AddUserToDockerGroup() {
	// Si eres root en servidor, agrega al usuario “real” si existe.
	// En contenedor o siendo root sin usuario objetivo, omite.
	if runningInContainer() {
		fmt.Println("⚠️  En contenedor no modifico grupos. Omite este paso.")
		return
	}
	current, _ := user.Current()
	uid0 := current != nil && current.Uid == "0"

	// Detecta usuario adecuado:
	u := os.Getenv("SUDO_USER")
	if u == "" && !uid0 && current != nil {
		u = current.Username
	}
	if u == "" || u == "root" {
		fmt.Println("ℹ️  Saltando: no hay usuario no-root claro para agregar a 'docker'.")
		return
	}

	// Crea grupo si falta y agrega usuario
	if err := utils.ExecShell(`getent group docker >/dev/null 2>&1 || sudo groupadd docker`); err != nil {
		fmt.Println("⚠️  No pude crear/verificar grupo docker:", err)
	}
	if err := utils.Exec("sudo", "usermod", "-aG", "docker", u); err != nil {
		fmt.Printf("⚠️  No pude agregar el usuario '%s' al grupo docker: %v\n", u, err)
		return
	}
	fmt.Printf("✅ Usuario '%s' agregado al grupo 'docker'. Cierra sesión y vuelve a entrar para aplicar cambios.\n", u)
}
