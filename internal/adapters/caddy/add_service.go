package caddy

import (
	"autohost-cli/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func AddService(serviceHost string, servicePort int) error {
	homeDir := utils.GetAutohostDir()
	caddyfilePath := filepath.Join(homeDir, ".autohost", "caddy", "Caddyfile")

	block := fmt.Sprintf(`
	%s {
		reverse_proxy 127.0.0.1:%d
	}
	`, serviceHost, servicePort)

	contentBytes, err := os.ReadFile(caddyfilePath)
	if err != nil {
		fmt.Println("❌ No se pudo leer el archivo Caddyfile:", err)
		return err
	}
	content := string(contentBytes)

	if strings.Contains(content, serviceHost) {
		fmt.Printf("⚠️ Ya existe una entrada para %s en el Caddyfile.\n", serviceHost)
		return nil
	}

	file, err := os.OpenFile(caddyfilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("❌ No se pudo abrir el archivo Caddyfile:", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(block)
	if err != nil {
		fmt.Println("❌ No se pudo escribir en el archivo Caddyfile:", err)
		return err
	}
	return nil
}
