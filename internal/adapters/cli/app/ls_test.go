package app

import (
	"testing"

	"autohost-cli/internal/platform/di"
)

func TestAppLsCmd(t *testing.T) {
	// Test basic command structure
	deps := di.Deps{}
	cmd := appLsCmd(deps)
	
	if cmd == nil {
		t.Fatal("appLsCmd returned nil")
	}
	
	if cmd.Use != "ls" {
		t.Errorf("Expected Use to be 'ls', got %s", cmd.Use)
	}

	if cmd.Short != "Lista las apps instaladas" {
		t.Errorf("Expected Short description to match, got %s", cmd.Short)
	}

	if cmd.RunE == nil {
		t.Error("Expected RunE to be set")
	}
}

func TestAppLsCmdBasicStructure(t *testing.T) {
	deps := di.Deps{}
	cmd := appLsCmd(deps)
	
	if cmd == nil {
		t.Fatal("appLsCmd returned nil")
	}
	
	if cmd.Use != "ls" {
		t.Errorf("Expected Use to be 'ls', got %s", cmd.Use)
	}
	
	if cmd.RunE == nil {
		t.Error("Expected RunE to be set")
	}
}