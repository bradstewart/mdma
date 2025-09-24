package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Servers []Server `yaml:"servers"`
	Audio   Audio    `yaml:"audio"`
	UI      UI       `yaml:"ui"`
}

type Server struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Audio struct {
	BufferSize   int    `yaml:"buffer_size"`
	SampleRate   int    `yaml:"sample_rate"`
	OutputDevice string `yaml:"output_device"`
}

type UI struct {
	Theme          string `yaml:"theme"`
	WaveformDetail string `yaml:"waveform_detail"`
}

func DefaultConfig() *Config {
	return &Config{
		Servers: []Server{
			{
				Name:     "local",
				URL:      "http://localhost:4533",
				Username: "",
				Password: "",
			},
		},
		Audio: Audio{
			BufferSize:   512,
			SampleRate:   44100,
			OutputDevice: "default",
		},
		UI: UI{
			Theme:          "dark",
			WaveformDetail: "medium",
		},
	}
}

func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config if it doesn't exist
		config := DefaultConfig()
		if err := config.Save(); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %v", err)
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "mdma", "config.yaml"), nil
}