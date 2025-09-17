package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAutohostDir(t *testing.T) {
	dir := GetAutohostDir()
	
	// Should not be empty
	if dir == "" {
		t.Error("GetAutohostDir returned empty string")
	}
	
	// Should contain .autohost
	if !filepath.IsAbs(dir) {
		t.Error("GetAutohostDir should return absolute path")
	}
	
	// Should end with .autohost
	if filepath.Base(dir) != ".autohost" {
		t.Errorf("Expected directory to end with .autohost, got %s", dir)
	}
}

func TestGetSubdir(t *testing.T) {
	tests := []struct {
		subdir   string
		expected string
	}{
		{"config", "config"},
		{"apps", "apps"},
		{"logs", "logs"},
		{"", ""}, // edge case
	}

	for _, tt := range tests {
		t.Run(tt.subdir, func(t *testing.T) {
			result := GetSubdir(tt.subdir)
			
			// Should be absolute path
			if !filepath.IsAbs(result) {
				t.Error("GetSubdir should return absolute path")
			}
			
			// Should end with expected subdir
			if tt.subdir != "" && filepath.Base(result) != tt.expected {
				t.Errorf("Expected path to end with %s, got %s", tt.expected, result)
			}
			
			// Should contain .autohost
			if !filepath.HasPrefix(result, GetAutohostDir()) {
				t.Errorf("Expected path to be under autohost dir, got %s", result)
			}
		})
	}
}

func TestIsInitialized(t *testing.T) {
	// Test when directory doesn't exist
	// We can't easily test this without potentially affecting the actual directory
	// So we'll test the function exists and returns a boolean
	result := IsInitialized()
	
	// Should return a boolean
	if _, ok := interface{}(result).(bool); !ok {
		t.Error("IsInitialized should return a boolean")
	}
	
	// Test logic: if GetAutohostDir exists, should return true
	dir := GetAutohostDir()
	if _, err := os.Stat(dir); err == nil {
		if !result {
			t.Error("Expected IsInitialized to return true when directory exists")
		}
	} else if os.IsNotExist(err) {
		if result {
			t.Error("Expected IsInitialized to return false when directory doesn't exist")
		}
	}
}