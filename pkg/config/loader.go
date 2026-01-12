package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	viper.SetDefault("commands.launch_frontend", "wt.exe -w 0 new-tab -d \"{{.FrontendPath}}\" cmd /k \"{{.FrontendCmd}}\"")
	viper.SetDefault("commands.launch_backend", "wt.exe -w 0 new-tab -d \"{{.BackendPath}}\" cmd /k \"{{.BackendCmd}}\"")
	viper.SetDefault("commands.launch_full", "wt.exe -w 0 new-tab -d \"{{.FrontendPath}}\" cmd /k \"{{.FrontendCmd}}\" ; split-pane -d \"{{.BackendPath}}\" cmd /k \"{{.BackendCmd}}\"")

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

	// Auto-repair bad templates from old config files
	if strings.Contains(cfg.Commands.LaunchFrontend, "%s") {
		cfg.Commands.LaunchFrontend = strings.ReplaceAll(cfg.Commands.LaunchFrontend, "%s", "{{.FrontendPath}}")
	}
	if strings.Contains(cfg.Commands.LaunchBackend, "%s") {
		cfg.Commands.LaunchBackend = strings.ReplaceAll(cfg.Commands.LaunchBackend, "%s", "{{.BackendPath}}")
	}
	if strings.Contains(cfg.Commands.LaunchFull, "%s") {
		// LaunchFull genellikle iki tane %s içerir (FrontendPath ve BackendPath)
		// Basit ReplaceAll ikisine de FrontendPath'i basar, bu HATALI olur.
		// Doğru strateji:
		// 1. Eğer içinde %s varsa ve sayıca 2 ise:
		if strings.Count(cfg.Commands.LaunchFull, "%s") == 2 {
			cfg.Commands.LaunchFull = strings.Replace(cfg.Commands.LaunchFull, "%s", "{{.FrontendPath}}", 1)
			cfg.Commands.LaunchFull = strings.Replace(cfg.Commands.LaunchFull, "%s", "{{.BackendPath}}", 1)
		} else {
			// Bilinmeyen format, güvenli moda geç ve varsayılanı kullan
			// Kullanıcı manuel olarak saçma bir şey yazmış olabilir, sıfırlamak en iyisi.
			cfg.Commands.LaunchFull = "wt.exe -w 0 new-tab -d \"{{.FrontendPath}}\" cmd /k \"{{.FrontendCmd}}\" ; split-pane -d \"{{.BackendPath}}\" cmd /k \"{{.BackendCmd}}\""
		}
	}

	// ---------------------------------------------------------
	// DEDUPLICATION & NORMALIZATION LOGIC
	// ---------------------------------------------------------
	// Bu kısım, Windows'un case-insensitive yapısından kaynaklanan
	// (örn: M:\Projeler vs m:\projeler) çift kayıtları temizler.
	if len(cfg.ProjectOverrides) > 0 {
		normalizedOverrides := make(map[string]domain.ProjectOverride)

		for path, override := range cfg.ProjectOverrides {
			// Windows yollarını normalize et (küçük harf)
			normalizedPath := strings.ToLower(path)

			// Eğer zaten varsa, boş olmayan değerleri koru
			if existing, exists := normalizedOverrides[normalizedPath]; exists {
				if override.Frontend != "" {
					existing.Frontend = override.Frontend
				}
				if override.Backend != "" {
					existing.Backend = override.Backend
				}
				normalizedOverrides[normalizedPath] = existing
			} else {
				normalizedOverrides[normalizedPath] = override
			}
		}

		// Temizlenmiş listeyi geri yükle
		cfg.ProjectOverrides = normalizedOverrides
	}

	return &cfg, nil
}

// SaveConfig writes the current configuration to disk
func SaveConfig(cfg *domain.Config) error {
	viper.Set("projects_paths", cfg.ProjectsPaths)
	viper.Set("project_overrides", cfg.ProjectOverrides)
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
