package app

import (
	"testing"

	"autohost-cli/internal/platform/di"
)

func TestAppInstallCmd(t *testing.T) {
	// Create mock dependencies
	deps := di.Deps{}

	// Test basic command structure
	cmd := appInstallCmd(deps)
	
	if cmd == nil {
		t.Fatal("appInstallCmd returned nil")
	}
	
	if cmd.Use != "install [nombre]" {
		t.Errorf("Expected Use to be 'install [nombre]', got %s", cmd.Use)
	}
	
	if cmd.Short != "Instala una aplicación (por ejemplo: nextcloud, bookstack, etc.)" {
		t.Errorf("Expected Short description to match, got %s", cmd.Short)
	}
	
	if cmd.RunE == nil {
		t.Error("Expected RunE to be set")
	}
}

func TestAppInstallCmdFlags(t *testing.T) {
	deps := di.Deps{}
	cmd := appInstallCmd(deps)
	
	// Check that the list flag exists
	listFlag := cmd.Flags().Lookup("list")
	if listFlag == nil {
		t.Error("Expected --list flag to exist")
	}
	
	if listFlag.Shorthand != "l" {
		t.Errorf("Expected shorthand to be 'l', got %s", listFlag.Shorthand)
	}

	if listFlag.Usage != "Mostrar catálogo e ignorar instalación" {
		t.Errorf("Expected specific usage text, got %s", listFlag.Usage)
	}
}