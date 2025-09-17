package app

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestAppStartCmd(t *testing.T) {
	// Test basic command structure
	cmd := appStartCmd()
	
	if cmd == nil {
		t.Fatal("appStartCmd returned nil")
	}
	
	if cmd.Use != "start [nombre]" {
		t.Errorf("Expected Use to be 'start [nombre]', got %s", cmd.Use)
	}
	
	if cmd.Short != "Inicia una aplicación" {
		t.Errorf("Expected Short description to match, got %s", cmd.Short)
	}
	
	// Test that it requires exactly one argument
	if cmd.Args == nil {
		t.Error("Expected Args to be set")
	}
}

func TestAppStartCmdWithInvalidArgs(t *testing.T) {
	cmd := appStartCmd()
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

func TestAppStartCmdValidArgument(t *testing.T) {
	cmd := appStartCmd()
	buf := &bytes.Buffer{}
	cmd.SetOutput(buf)
	cmd.SetErr(buf)

	// Test with valid single argument
	// Note: This will fail because we don't have actual docker/app setup
	// but we can test that the command structure is correct
	cmd.SetArgs([]string{"testapp"})
	err := cmd.ExecuteContext(context.Background())

	// We expect an error because the app doesn't exist, but the command structure should be valid
	output := buf.String()
	
	// The command should attempt to start the app and likely fail
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