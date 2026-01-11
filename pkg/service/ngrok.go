package service

import (
	"devterminal/pkg/domain"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

type NgrokService struct {
	Config *domain.Config
}

func NewNgrokService(cfg *domain.Config) *NgrokService {
	return &NgrokService{Config: cfg}
}

// CheckCommonPaths looks for ngrok in standard locations without scanning full PATH
func (n *NgrokService) CheckCommonPaths() string {
	// 1. Check Config Cache first
	if n.Config.NgrokPath != "" && n.ValidatePath(n.Config.NgrokPath) {
		return n.Config.NgrokPath
	}

	// 2. Check Common Windows Paths
	localAppData := os.Getenv("LOCALAPPDATA")
	candidates := []string{
		filepath.Join(localAppData, "Microsoft", "WinGet", "Links", "ngrok.exe"),
		"C:\\ProgramData\\chocolatey\\bin\\ngrok.exe",
		"C:\\Program Files\\ngrok\\ngrok.exe",
	}

	for _, path := range candidates {
		if n.ValidatePath(path) {
			// Found it! Cache it.
			n.SavePath(path)
			return path
		}
	}
	return ""
}

// ValidatePath checks if the file exists
func (n *NgrokService) ValidatePath(path string) bool {
	// Remove quotes if user added them
	cleanPath := strings.Trim(path, "\"'")
	info, err := os.Stat(cleanPath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// SavePath persists the path to config
func (n *NgrokService) SavePath(path string) error {
	cleanPath := strings.Trim(path, "\"'")
	n.Config.NgrokPath = cleanPath
	viper.Set("ngrok_path", cleanPath)
	return viper.WriteConfig()
}

// GetExecutable resolves the path or returns empty
func (n *NgrokService) GetExecutable() string {
	if n.Config.NgrokPath != "" {
		return n.Config.NgrokPath
	}
	return n.CheckCommonPaths()
}

// HasAuthToken checks if default config file contains an authtoken
func (n *NgrokService) HasAuthToken() bool {
	exe := n.GetExecutable()
	if exe == "" {
		return false
	}

	cmd := exec.Command(exe, "config", "check")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	output := string(out)

	re := regexp.MustCompile(`at\s+(.*ngrok\.yml)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return false
	}
	configPath := strings.TrimSpace(matches[1])

	data, err := os.ReadFile(configPath)
	if err != nil {
		return false
	}

	return strings.Contains(string(data), "authtoken:")
}

// SetAuthToken executes command to add token
func (n *NgrokService) SetAuthToken(token string) error {
	exe := n.GetExecutable()
	if exe == "" {
		return os.ErrNotExist
	}
	cmd := exec.Command(exe, "config", "add-authtoken", token)
	return cmd.Run()
}
