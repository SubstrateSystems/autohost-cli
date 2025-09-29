package caddy

import (
	"autohost-cli/utils"
	"fmt"
	"os"
	"path/filepath"
)

func AddService(serviceHost string, servicePort int) error {
	// 1) Ruta del snippet
	dir := filepath.Join("/", "etc", "caddy", "autohost")
	if err := os.MkdirAll(dir, 0o775); err != nil {
		return fmt.Errorf("crear dir snippets: %w", err)
	}
	snippet := filepath.Join(dir, fmt.Sprintf("%s.caddy", serviceHost))

	// 2) Bloque Caddy por host
	block := fmt.Sprintf(`%s {
    reverse_proxy 127.0.0.1:%d
}
`, serviceHost, servicePort)

	// 3) Si ya existe, no duplicar
	if _, err := os.Stat(snippet); err == nil {
		fmt.Printf("⚠️ Ya existe una entrada para %s (%s).\n", serviceHost, snippet)
	} else {
		if err := os.WriteFile(snippet, []byte(block), 0o664); err != nil {
			return fmt.Errorf("escribir snippet: %w", err)
		}
		fmt.Printf("✅ Snippet creado: %s\n", snippet)
	}
	restartCaddy()
	fmt.Println("🔄 Caddy recargado.")

	return nil
}

func restartCaddy() error {
	fmt.Println("🔄 Reiniciando Caddy...")
	err := utils.Exec("sudo", "systemctl", "restart", "caddy")
	return err
}
