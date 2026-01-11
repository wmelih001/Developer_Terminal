package config

import (
	"fmt"
	"os"
	"path/filepath"

	"devterminal/pkg/domain"

	"github.com/spf13/viper"
)

// LoadConfig reads configuration from ~/.godmode/config.yaml
func LoadConfig() (*domain.Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".devterminal")
	configName := "config"

	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".") // Search in current directory too
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")

	// Defaults
	viper.SetDefault("projects_paths", []string{})
	viper.SetDefault("ignored_files", []string{
		".git", "node_modules", "dist", ".next", ".idea", ".vscode",
	})
	// Default launch commands (optimistic defaults for Windows Terminal)
	viper.SetDefault("commands.launch_frontend", "wt.exe -w 0 new-tab -d \"%s\" cmd /k \"{{.FrontendCmd}}\"")
	viper.SetDefault("commands.launch_backend", "wt.exe -w 0 new-tab -d \"%s\" cmd /k \"{{.BackendCmd}}\"")
	viper.SetDefault("commands.launch_full", "wt.exe -w 0 new-tab -d \"%s\" cmd /k \"{{.FrontendCmd}}\" ; split-pane -d \"%s\" cmd /k \"{{.BackendCmd}}\"")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if we want to run with defaults
			// Or create it? For now, just return defaults.
			fmt.Println("Warning: Config file not found, using defaults.")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg domain.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}

// SaveConfig writes the current configuration to disk
func SaveConfig(cfg *domain.Config) error {
	viper.Set("projects_paths", cfg.ProjectsPaths)
	// Add other fields if necessary to sync back to viper before saving
	// For now, we mainly accept project paths updates

	// Ensure directory exists
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(home, ".devterminal")
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
	}

	configFile := filepath.Join(configPath, "config.yaml")
	return viper.WriteConfigAs(configFile)
}
