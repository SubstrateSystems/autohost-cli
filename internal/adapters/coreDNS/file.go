package coredns

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func createTemplateCorefile(tailIP string) string {
	return fmt.Sprintf(`# CoreDNS (Docker) para AutoHost
.:53 {
    bind %s
    hosts {

        fallthrough
    }
    forward . 1.1.1.1
    log
    errors
}
`, tailIP)
}

func UpdateCorefile(subdomain, appIP string) error {
	home, _ := os.UserHomeDir()
	corefilePath := filepath.Join(home, ".autohost", "coredns")
	// Leer Corefile actual
	data, err := os.ReadFile(corefilePath)
	if err != nil {
		return fmt.Errorf("no pude leer Corefile: %w", err)
	}
	content := string(data)

	// Construir la nueva entrada
	newLine := fmt.Sprintf("    %s %s\n", appIP, subdomain)

	// Verificar si ya existe
	if strings.Contains(content, newLine) {
		fmt.Println("ℹ️ La entrada ya existe en el Corefile.")
		return nil
	}

	// Insertar dentro del bloque hosts antes de "fallthrough"
	updated := strings.Replace(
		content,
		"    fallthrough",
		newLine+"    fallthrough",
		1, // solo la primera ocurrencia
	)

	// Escribir de nuevo el Corefile actualizado
	if err := os.WriteFile(corefilePath, []byte(updated), 0644); err != nil {
		return fmt.Errorf("no pude escribir Corefile actualizado: %w", err)
	}

	fmt.Println("✅ Corefile actualizado con:", newLine)
	return nil
}
