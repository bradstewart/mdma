package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if len(config.Servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(config.Servers))
	}

	if config.Servers[0].Name != "local" {
		t.Errorf("Expected server name 'local', got '%s'", config.Servers[0].Name)
	}

	if config.Audio.SampleRate != 44100 {
		t.Errorf("Expected sample rate 44100, got %d", config.Audio.SampleRate)
	}

	if config.UI.Theme != "dark" {
		t.Errorf("Expected theme 'dark', got '%s'", config.UI.Theme)
	}
}

func TestConfigSaveLoad(t *testing.T) {
	// Create temporary config directory
	tempDir, err := os.MkdirTemp("", "mdma-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override config path for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create and save config
	config := DefaultConfig()
	config.Servers[0].Username = "testuser"
	config.Audio.BufferSize = 1024

	err = config.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config and verify
	loadedConfig, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedConfig.Servers[0].Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", loadedConfig.Servers[0].Username)
	}

	if loadedConfig.Audio.BufferSize != 1024 {
		t.Errorf("Expected buffer size 1024, got %d", loadedConfig.Audio.BufferSize)
	}
}