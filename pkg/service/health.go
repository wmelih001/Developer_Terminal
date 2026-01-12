package service

import (
	"os"
	"path/filepath"
	"strings"
)

type HealthIssue struct {
	Description string
	Points      int // How many points were lost
}

type HealthReport struct {
	Score       int
	MaxScore    int
	Issues      []HealthIssue
	PassedItems []string
}

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) CheckHealth(projectPath string) HealthReport {
	score := 0
	maxScore := 100
	var issues []HealthIssue
	var passed []string

	// Flags for detection
	hasGit := false
	hasDep := false
	hasReadme := false
	hasDocker := false
	hasEnv := false
	hasCICD := false
	hasLinter := false
	hasLicense := false

	// Tek seferlik recursive tarama
	_ = filepath.WalkDir(projectPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		name := d.Name()
		lowerName := strings.ToLower(name)

		// Gereksiz klasörleri atla
		if d.IsDir() && path != projectPath {
			// node_modules, .git (içeriği), vendor, dist, build atla
			if name == "node_modules" || name == ".git" || name == "vendor" || name == "dist" || name == "build" || name == ".next" {
				return filepath.SkipDir
			}
			// Gizli klasörleri atla (.config ve .github, .circleci hariç - bunlar aranıyor)
			if strings.HasPrefix(name, ".") {
				if name != ".config" && name != ".github" && name != ".circleci" {
					return filepath.SkipDir
				}
			}
		}

		// Derinlik kontrolü (Kökten en fazla 3 seviye aşağı in)
		rel, _ := filepath.Rel(projectPath, path)
		depth := strings.Count(rel, string(os.PathSeparator))
		if depth > 3 {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// --- 1. Git ---
		if d.IsDir() && name == ".git" {
			hasGit = true
		}

		// --- 2. Dependencies ---
		if !d.IsDir() {
			switch name {
			case "package.json", "go.mod", "requirements.txt", "pom.xml", "build.gradle", "Gemfile", "composer.json", "mix.exs", "Cargo.toml":
				hasDep = true
			}
		}

		// --- 3. Readme ---
		if !d.IsDir() && (strings.HasPrefix(lowerName, "readme")) {
			hasReadme = true
		}

		// --- 4. Docker ---
		if !d.IsDir() && (strings.HasPrefix(lowerName, "dockerfile") || strings.HasPrefix(lowerName, "docker-compose") || name == "Containerfile") {
			hasDocker = true
		}

		// --- 5. Env / Config ---
		if !d.IsDir() {
			if strings.HasPrefix(name, ".env") ||
				name == "config.yaml" ||
				name == "config.json" ||
				name == "config.js" {
				hasEnv = true
			}
		}

		// --- 6. CI/CD ---
		if d.IsDir() && (name == ".github" || name == ".circleci") {
			hasCICD = true
		}
		if !d.IsDir() && (name == ".gitlab-ci.yml" || name == "azure-pipelines.yml" || name == "Jenkinsfile") {
			hasCICD = true
		}

		// --- 7. Linter ---
		if !d.IsDir() {
			if strings.HasPrefix(name, ".eslintrc") ||
				strings.HasPrefix(name, ".prettierrc") ||
				name == "golangci.yml" ||
				name == ".pylintrc" ||
				name == "checkstyle.xml" ||
				name == "rubocop.yml" {
				hasLinter = true
			}
		}

		// --- 8. License ---
		if !d.IsDir() && (strings.HasPrefix(lowerName, "license") || strings.HasPrefix(lowerName, "copying")) {
			hasLicense = true
		}

		return nil
	})

	// Puanlama ve Raporlama
	if hasGit {
		score += 20
		passed = append(passed, "Git Repository (20p)")
	} else {
		issues = append(issues, HealthIssue{"Git başlatılmamış (.git yok)", 20})
	}

	if hasDep {
		score += 20
		passed = append(passed, "Bağımlılık Dosyası (20p)")
	} else {
		issues = append(issues, HealthIssue{"Bağımlılık dosyası bulunamadı", 20})
	}

	if hasReadme {
		score += 10
		passed = append(passed, "README (10p)")
	} else {
		issues = append(issues, HealthIssue{"README dosyası eksik", 10})
	}

	if hasDocker {
		score += 10
		passed = append(passed, "Konteyner Yapılandırması (10p)")
	} else {
		issues = append(issues, HealthIssue{"Docker/Konteyner yapılandırması yok", 10})
	}

	if hasEnv {
		score += 10
		passed = append(passed, "Ortam Değişkenleri (.env)/Config (10p)")
	} else {
		issues = append(issues, HealthIssue{"Konfigürasyon/Env dosyası yok", 10})
	}

	if hasCICD {
		score += 10
		passed = append(passed, "CI/CD Yapılandırması (10p)")
	} else {
		issues = append(issues, HealthIssue{"CI/CD yapılandırması bulunamadı", 10})
	}

	if hasLinter {
		score += 10
		passed = append(passed, "Linter/Formatter Ayarları (10p)")
	} else {
		issues = append(issues, HealthIssue{"Linter/Formatter ayarları eksik", 10})
	}

	if hasLicense {
		score += 10
		passed = append(passed, "Lisans Dosyası (10p)")
	} else {
		issues = append(issues, HealthIssue{"Lisans dosyası eksik", 10})
	}

	return HealthReport{
		Score:       score,
		MaxScore:    maxScore,
		Issues:      issues,
		PassedItems: passed,
	}
}
