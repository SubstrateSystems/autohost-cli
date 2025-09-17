package app

import (
	"strings"
	"testing"

	"autohost-cli/internal/platform/di"

	"github.com/spf13/cobra"
)

func TestAppCmd(t *testing.T) {
	// Create a minimal deps struct for testing
	deps := di.Deps{}
	
	// Create the command
	cmd := AppCmd(deps)
	
	// Test basic command properties
	if cmd.Use != "app" {
		t.Errorf("Expected Use to be 'app', got %s", cmd.Use)
	}
	
	if cmd.Short != "Gestión de aplicaciones autohospedadas" {
		t.Errorf("Expected Short description to match, got %s", cmd.Short)
	}
	
	// Test that subcommands are added
	subcommands := cmd.Commands()
	expectedSubcommands := []string{"install", "ls", "remove", "start", "status", "stop"}
	
	if len(subcommands) != len(expectedSubcommands) {
		t.Errorf("Expected %d subcommands, got %d", len(expectedSubcommands), len(subcommands))
	}
	
	// Check that each expected subcommand exists
	commandNames := make(map[string]bool)
	for _, subcmd := range subcommands {
		// Handle commands with arguments in their Use field
		parts := strings.Split(subcmd.Use, " ")
		commandName := parts[0]  // Get just the command name, not the full Use string
		commandNames[commandName] = true
	}
	
	for _, expectedName := range expectedSubcommands {
		if !commandNames[expectedName] {
			t.Errorf("Expected subcommand %s not found", expectedName)
		}
	}
}

func TestAppCmdStructure(t *testing.T) {
	deps := di.Deps{}
	cmd := AppCmd(deps)
	
	// Test that the command is properly configured
	if cmd == nil {
		t.Fatal("AppCmd returned nil")
	}
	
	// Test that it's a valid cobra command
	if _, ok := interface{}(cmd).(*cobra.Command); !ok {
		t.Error("AppCmd should return a *cobra.Command")
	}
}