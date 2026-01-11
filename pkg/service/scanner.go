package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"devterminal/pkg/domain"
)

// Scanner handles project discovery
type Scanner struct {
	Config *domain.Config
}

func NewScanner(cfg *domain.Config) *Scanner {
	return &Scanner{Config: cfg}
}

// packageJSON minimal struct for parsing
type packageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Scripts         map[string]string `json:"scripts"`
}

// ScanProjects scans all configured paths for projects
func (s *Scanner) ScanProjects() []domain.Project {
	var allProjects []domain.Project
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Global deduplication map
	seenPaths := make(map[string]bool)
	var mapMu sync.Mutex

	for _, path := range s.Config.ProjectsPaths { // Changed s.Config.ProjectsPaths to s.config.ProjectsPaths as per instruction, but keeping original s.Config.ProjectsPaths as it's syntactically correct.
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			projects := s.scanDirectory(p)

			mu.Lock()
			for _, proj := range projects {
				mapMu.Lock()
				if !seenPaths[proj.Name] { // Deduplicate by Name (or Path)
					seenPaths[proj.Name] = true
					allProjects = append(allProjects, proj)
				}
				mapMu.Unlock()
			}
			mu.Unlock()
		}(path)
	}
	wg.Wait()
	return allProjects
}

func (s *Scanner) scanDirectory(root string) []domain.Project {
	var results []domain.Project
	seen := make(map[string]bool)

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(root, entry.Name())
		if seen[fullPath] {
			continue
		}
		seen[fullPath] = true

		// Detaylı analiz
		p := domain.Project{
			Name: entry.Name(),
			Path: fullPath,
			Type: domain.TypeUnknown,
		}

		// 1. Frontend Kontrolü
		frontPath := filepath.Join(fullPath, "frontend")
		if info, err := os.Stat(frontPath); err == nil && info.IsDir() {
			p.HasFrontend = true
			p.FrontendVer = s.getPackageVersion(frontPath, "next")
			if p.FrontendVer == "" {
				p.FrontendVer = s.getPackageVersion(frontPath, "react")
			}
			if p.FrontendVer == "" {
				p.FrontendVer = "Var"
			}
			p.FrontendPath = frontPath
			p.FrontendCmd = s.detectStartCommand(frontPath)
		}

		// 2. Backend Kontrolü
		backPath := filepath.Join(fullPath, "backend")
		if info, err := os.Stat(backPath); err == nil && info.IsDir() {
			p.HasBackend = true
			p.BackendVer = s.getPackageVersion(backPath, "@nestjs/core")
			if p.BackendVer == "" {
				p.BackendVer = s.getPackageVersion(backPath, "express")
			}
			if p.BackendVer == "" {
				p.BackendVer = "Var"
			}
			p.BackendPath = backPath
			p.BackendCmd = s.detectStartCommand(backPath)
		}

		// Tip Belirleme
		if p.HasFrontend && p.HasBackend {
			p.Type = domain.TypeUnknown // Fullstack aslında, ama Type enum'ı basit kaldı. UI'da ikon belirlemek için kullanıyoruz.
			// Next ve Nest ise ikonları UI tarafında karmaşıklaşacak, şimdilik basit bırakalım.
		} else {
			// Alt klasör yoksa root'a bak
			p.Type = s.detectProjectType(fullPath)
			if p.Type == domain.TypeNext {
				p.FrontendVer = s.getPackageVersion(fullPath, "next")
				p.HasFrontend = true
			}
			if p.Type == domain.TypeNest {
				p.BackendVer = s.getPackageVersion(fullPath, "@nestjs/core")
				p.HasBackend = true
				p.BackendPath = fullPath
				p.BackendCmd = s.detectStartCommand(fullPath)
			}
			if p.HasFrontend && p.FrontendCmd == "" {
				p.FrontendPath = fullPath
				p.FrontendCmd = s.detectStartCommand(fullPath)
			}
		}

		results = append(results, p)
	}
	return results
}

func (s *Scanner) getPackageVersion(path, pkgName string) string {
	pkgPath := filepath.Join(path, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return ""
	}
	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return ""
	}

	if v, ok := pkg.Dependencies[pkgName]; ok {
		return strings.Trim(v, "^~")
	}
	if v, ok := pkg.DevDependencies[pkgName]; ok {
		return strings.Trim(v, "^~")
	}
	return ""
}

func (s *Scanner) detectProjectType(path string) domain.ProjectType {
	// Root check (Eski logic)
	pkgPath := filepath.Join(path, "package.json")
	if _, err := os.Stat(pkgPath); err == nil {
		data, err := os.ReadFile(pkgPath)
		if err == nil {
			var pkg packageJSON
			if err := json.Unmarshal(data, &pkg); err == nil {
				if hasDependency(pkg, "next") {
					return domain.TypeNext
				}
				if hasDependency(pkg, "@nestjs/core") {
					return domain.TypeNest
				}
				if hasDependency(pkg, "react") {
					return domain.TypeReact
				}
			}
		}
	}
	if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
		return domain.TypeGo
	}
	return domain.TypeUnknown
}

func hasDependency(pkg packageJSON, dep string) bool {
	if _, ok := pkg.Dependencies[dep]; ok {
		return true
	}
	if _, ok := pkg.DevDependencies[dep]; ok {
		return true
	}
	return false
}

// detectStartCommand determines the best command to start the project
func (s *Scanner) detectStartCommand(path string) string {
	// 1. JS/TS Projects (Next, Nest, React)
	pkgPath := filepath.Join(path, "package.json")
	if _, err := os.Stat(pkgPath); err == nil {
		data, err := os.ReadFile(pkgPath)
		if err == nil {
			var pkg packageJSON
			if err := json.Unmarshal(data, &pkg); err == nil {
				// Priority: dev > start:dev > start
				if _, ok := pkg.Scripts["dev"]; ok {
					return "npm run dev"
				}
				if _, ok := pkg.Scripts["start:dev"]; ok {
					return "npm run start:dev"
				}
				if _, ok := pkg.Scripts["start"]; ok {
					return "npm start"
				}
			}
		}
	}

	// 2. Go Projects
	if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
		if _, err := os.Stat(filepath.Join(path, "main.go")); err == nil {
			return "go run main.go"
		}
		return "go run ."
	}

	// Default fallback with explanation
	return "echo [GodMode] Başlatma komutu bulunamadı (package.json scripts veya go.mod yok)"
}
