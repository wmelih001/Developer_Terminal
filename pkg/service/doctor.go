package service

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"devterminal/pkg/domain"
)

type Doctor struct {
	Config *domain.Config
}

func NewDoctor(cfg *domain.Config) *Doctor {
	return &Doctor{Config: cfg}
}

type NpmOutdatedResult map[string]struct {
	Current string `json:"current"`
	Wanted  string `json:"wanted"`
	Latest  string `json:"latest"`
}

// CheckDependencies runs npm outdated in the project path with comprehensive validation
func (d *Doctor) CheckDependencies(path string) (NpmOutdatedResult, error) {
	// 1. Check if npm is installed
	npmCheck := exec.Command("npm", "--version")
	if err := npmCheck.Run(); err != nil {
		return nil, fmt.Errorf("npm bulunamadı. Lütfen Node.js ve npm'in kurulu olduğundan emin olun")
	}

	// 2. Check if package.json exists
	packageJSONPath := filepath.Join(path, "package.json")
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("package.json bulunamadı. Bu proje bir Node.js projesi değil")
	}

	// 3. Run npm outdated
	cmd := exec.Command("npm", "outdated", "--json")
	cmd.Dir = path

	// Use CombinedOutput to capture both stdout and stderr
	output, err := cmd.CombinedOutput()

	// npm outdated returns exit code 1 if there are outdated packages
	// This is NOT an error for us - we want to parse the output
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Exit code 1 means outdated packages exist - this is what we want!
			if exitErr.ExitCode() == 1 {
				// Continue to parse output below
			} else {
				// Other exit codes indicate real errors
				return nil, fmt.Errorf("npm outdated komutu başarısız (exit %d): %s", exitErr.ExitCode(), string(output))
			}
		} else {
			return nil, fmt.Errorf("npm komutu çalıştırılamadı: %v", err)
		}
	}

	// Empty output means all packages are up to date (exit code was 0)
	if len(output) == 0 {
		return nil, nil
	}

	// Parse JSON result
	var res NpmOutdatedResult
	if err := json.Unmarshal(output, &res); err != nil {
		return nil, fmt.Errorf("npm çıktısı parse edilemedi: %v\nOutput: %s", err, string(output))
	}

	return res, nil
}
