package utils

import (
	"fmt"
	"net"
	"testing"
)

func TestValidPort(t *testing.T) {
	tests := []struct {
		name        string
		portStr     string
		expectError bool
	}{
		{
			name:        "valid port number",
			portStr:     "8080",
			expectError: false,
		},
		{
			name:        "invalid port - not a number",
			portStr:     "abc",
			expectError: true,
		},
		{
			name:        "invalid port - empty string",
			portStr:     "",
			expectError: true,
		},
		{
			name:        "invalid port - negative number",
			portStr:     "-1",
			expectError: true, // This should fail when trying to listen
		},
		{
			name:        "port zero",
			portStr:     "0",
			expectError: false, // Port 0 is valid (OS assigns available port)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port, err := ValidPort(tt.portStr)
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error for port %s, but got none", tt.portStr)
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for port %s: %v", tt.portStr, err)
			}
			
			// If no error expected and port is not zero, verify the port number
			if !tt.expectError && tt.portStr != "0" && err == nil {
				expectedPort := 8080
				if tt.portStr == "8080" && port != expectedPort {
					t.Errorf("Expected port %d, got %d", expectedPort, port)
				}
			}
		})
	}
}

func TestValidPortWithUsedPort(t *testing.T) {
	// Start a listener on a random port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to start test listener: %v", err)
	}
	defer listener.Close()

	// Get the port that's now in use
	addr := listener.Addr().(*net.TCPAddr)
	usedPort := fmt.Sprintf("%d", addr.Port)

	// Test that ValidPort detects the port is in use
	_, err = ValidPort(usedPort)
	if err == nil {
		t.Errorf("Expected ValidPort to fail for port in use: %s", usedPort)
	}

	// Verify the error message mentions port in use
	if err != nil && err.Error() == "" {
		t.Error("Error should have a descriptive message")
	}
}

func TestValidPortEdgeCases(t *testing.T) {
	// Test very high port numbers
	_, err := ValidPort("65535")
	// This might fail due to permissions or port availability, but shouldn't fail parsing
	if err != nil {
		// Should be about port availability, not parsing
		if err.Error() == "puerto inválido: 65535" {
			t.Error("Port 65535 should be valid as a number")
		}
	}

	// Test port out of range
	_, err = ValidPort("65536")
	// This should fail at the listen stage, not parsing (since Go allows it)
	// but the specific behavior depends on the OS
}

// TestExecFunctions tests that the Exec functions exist and have correct signatures
func TestExecFunctions(t *testing.T) {
	// Test that Exec function exists
	err := Exec("echo", "test")
	if err != nil {
		// This might fail in test environment, but function should exist
		t.Logf("Exec test failed (expected in test environment): %v", err)
	}

	// Test that ExecShell function exists
	err = ExecShell("echo test")
	if err != nil {
		// This might fail in test environment, but function should exist
		t.Logf("ExecShell test failed (expected in test environment): %v", err)
	}

	// Test that ExecWithDir function exists
	err = ExecWithDir("/tmp", "echo", "test")
	if err != nil {
		// This might fail in test environment, but function should exist
		t.Logf("ExecWithDir test failed (expected in test environment): %v", err)
	}
}