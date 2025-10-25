package infra

import "fmt"

// RunStep imprime estado y envuelve el error con contexto uniforme
func RunStep(nombre string, fn func() error) error {
	fmt.Printf("🔄 %s...\n", nombre)
	if err := fn(); err != nil {
		return fmt.Errorf("❌ %s: %w", nombre, err)
	}
	fmt.Printf("✅ %s completado.\n", nombre)
	return nil
}
