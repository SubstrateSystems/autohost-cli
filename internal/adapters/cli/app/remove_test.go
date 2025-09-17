package app

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"autohost-cli/internal/platform/di"
)

func TestAppRemoveCmd(t *testing.T) {
	// Create mock dependencies
	deps := di.Deps{}

	// Test basic command structure
	cmd := appRemoveCmd(deps)
	
	if cmd == nil {
		t.Fatal("appRemoveCmd returned nil")
	}
	
	if cmd.Use != "remove [nombre]" {
		t.Errorf("Expected Use to be 'remove [nombre]', got %s", cmd.Use)
	}
	
	if cmd.Short != "Elimina una aplicación" {
		t.Errorf("Expected Short description to match, got %s", cmd.Short)
	}
	
	// Test that it requires exactly one argument
	if cmd.Args == nil {
		t.Error("Expected Args to be set")
	}
}

func TestAppRemoveCmdWithInvalidArgs(t *testing.T) {
	deps := di.Deps{}
	cmd := appRemoveCmd(deps)
	buf := &bytes.Buffer{}
	cmd.SetOutput(buf)
	cmd.SetErr(buf)

	// Test with no arguments
	cmd.SetArgs([]string{})
	err := cmd.ExecuteContext(context.Background())
	
	if err == nil {
		t.Error("Expected error when no arguments provided")
	}

	// Test with multiple arguments  
	cmd.SetArgs([]string{"app1", "app2"})
	err = cmd.ExecuteContext(context.Background())
	
	if err == nil {
		t.Error("Expected error when multiple arguments provided")
	}
}

func TestAppRemoveCmdValidArgument(t *testing.T) {
	deps := di.Deps{}
	cmd := appRemoveCmd(deps)
	buf := &bytes.Buffer{}
	cmd.SetOutput(buf)
	cmd.SetErr(buf)

	// Test with valid single argument
	cmd.SetArgs([]string{"testapp"})
	err := cmd.ExecuteContext(context.Background())

	// We expect an error because the app doesn't exist, but the command structure should be valid
	output := buf.String()
	
	// The command should attempt to remove the app and likely fail
	// but it should not fail due to argument validation
	if strings.Contains(output, "accepts 1 arg") {
		t.Error("Command failed due to argument validation, which shouldn't happen with valid args")
	}

	// The error should be related to app not existing or docker not available
	// not argument parsing
	if err != nil && strings.Contains(err.Error(), "accepts") {
		t.Errorf("Error should not be about argument count: %v", err)
	}
}