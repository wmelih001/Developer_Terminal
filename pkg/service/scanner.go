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

// composerJSON for PHP projects
type composerJSON struct {
	Require    map[string]string `json:"require"`
	RequireDev map[string]string `json:"require-dev"`
}

// TechSignature represents a technology detection signature
type TechSignature struct {
	Type       domain.ProjectType
	IsFrontend bool                             // true = Frontend, false = Backend
	CheckFunc  func(path string) (bool, string) // Returns (found, version)
}

// ScanProjects scans all configured paths for projects
func (s *Scanner) ScanProjects() []domain.Project {
	var allProjects []domain.Project
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Global deduplication map
	seenPaths := make(map[string]bool)
	var mapMu sync.Mutex

	for _, path := range s.Config.ProjectsPaths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			projects := s.scanDirectory(p)

			mu.Lock()
			for _, proj := range projects {
				mapMu.Lock()
				if !seenPaths[proj.Name] {
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

// getFrontendSignatures returns all frontend technology signatures
func (s *Scanner) getFrontendSignatures() []TechSignature {
	return []TechSignature{
		// Next.js
		{Type: domain.TypeNext, IsFrontend: true, CheckFunc: func(path string) (bool, string) {
			ver := s.getPackageVersion(path, "next")
			return ver != "", ver
		}},
		// React
		{Type: domain.TypeReact, IsFrontend: true, CheckFunc: func(path string) (bool, string) {
			// React ama Next değilse
			if s.getPackageVersion(path, "next") != "" {
				return false, ""
			}
			ver := s.getPackageVersion(path, "react")
			return ver != "", ver
		}},
		// Vue
		{Type: domain.TypeVue, IsFrontend: true, CheckFunc: func(path string) (bool, string) {
			ver := s.getPackageVersion(path, "vue")
			return ver != "", ver
		}},
		// Vite (config file check)
		{Type: domain.TypeVite, IsFrontend: true, CheckFunc: func(path string) (bool, string) {
			if _, err := os.Stat(filepath.Join(path, "vite.config.ts")); err == nil {
				return true, "Var"
			}
			if _, err := os.Stat(filepath.Join(path, "vite.config.js")); err == nil {
				return true, "Var"
			}
			return false, ""
		}},
		// React Native
		{Type: domain.TypeReactNative, IsFrontend: true, CheckFunc: func(path string) (bool, string) {
			ver := s.getPackageVersion(path, "react-native")
			return ver != "", ver
		}},
		// Mobile (Native - Android/iOS folders)
		{Type: domain.TypeMobile, IsFrontend: true, CheckFunc: func(path string) (bool, string) {
			androidPath := filepath.Join(path, "android")
			iosPath := filepath.Join(path, "ios")
			hasAndroid := false
			hasIOS := false
			if info, err := os.Stat(androidPath); err == nil && info.IsDir() {
				hasAndroid = true
			}
			if info, err := os.Stat(iosPath); err == nil && info.IsDir() {
				hasIOS = true
			}
			if hasAndroid && hasIOS {
				return true, "iOS & Android"
			} else if hasAndroid {
				return true, "Android"
			} else if hasIOS {
				return true, "iOS"
			}
			return false, ""
		}},
		// HTML (Static Website)
		{Type: domain.TypeHTML, IsFrontend: true, CheckFunc: func(path string) (bool, string) {
			// Check for index.html - indicates static HTML project
			if _, err := os.Stat(filepath.Join(path, "index.html")); err == nil {
				// Make sure it's not a framework project (no package.json with frameworks)
				if s.getPackageVersion(path, "next") != "" ||
					s.getPackageVersion(path, "react") != "" ||
					s.getPackageVersion(path, "vue") != "" {
					return false, ""
				}
				return true, "Var"
			}
			return false, ""
		}},
		// TypeScript (Standalone TS project)
		{Type: domain.TypeTypeScript, IsFrontend: true, CheckFunc: func(path string) (bool, string) {
			// Check for tsconfig.json - indicates TypeScript project
			if _, err := os.Stat(filepath.Join(path, "tsconfig.json")); err == nil {
				// Make sure it's not a framework project
				if s.getPackageVersion(path, "next") != "" ||
					s.getPackageVersion(path, "react") != "" ||
					s.getPackageVersion(path, "vue") != "" ||
					s.getPackageVersion(path, "@nestjs/core") != "" {
					return false, ""
				}
				return true, "Var"
			}
			return false, ""
		}},
	}
}

// getBackendSignatures returns all backend technology signatures
func (s *Scanner) getBackendSignatures() []TechSignature {
	return []TechSignature{
		// NestJS
		{Type: domain.TypeNest, IsFrontend: false, CheckFunc: func(path string) (bool, string) {
			ver := s.getPackageVersion(path, "@nestjs/core")
			return ver != "", ver
		}},
		// Express
		{Type: domain.TypeExpress, IsFrontend: false, CheckFunc: func(path string) (bool, string) {
			// Express ama Nest değilse
			if s.getPackageVersion(path, "@nestjs/core") != "" {
				return false, ""
			}
			ver := s.getPackageVersion(path, "express")
			return ver != "", ver
		}},
		// Go
		{Type: domain.TypeGo, IsFrontend: false, CheckFunc: func(path string) (bool, string) {
			if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
				return true, "Var"
			}
			return false, ""
		}},
		// Django
		{Type: domain.TypeDjango, IsFrontend: false, CheckFunc: func(path string) (bool, string) {
			if _, err := os.Stat(filepath.Join(path, "manage.py")); err == nil {
				return true, "Var"
			}
			return false, ""
		}},
		// Flask
		{Type: domain.TypeFlask, IsFrontend: false, CheckFunc: func(path string) (bool, string) {
			// Check for app.py or wsgi.py
			hasAppPy := false
			if _, err := os.Stat(filepath.Join(path, "app.py")); err == nil {
				hasAppPy = true
			}
			if _, err := os.Stat(filepath.Join(path, "wsgi.py")); err == nil {
				hasAppPy = true
			}
			if !hasAppPy {
				return false, ""
			}
			// Check requirements.txt for flask
			reqPath := filepath.Join(path, "requirements.txt")
			if data, err := os.ReadFile(reqPath); err == nil {
				if strings.Contains(strings.ToLower(string(data)), "flask") {
					return true, "Var"
				}
			}
			return false, ""
		}},
		// Laravel
		{Type: domain.TypeLaravel, IsFrontend: false, CheckFunc: func(path string) (bool, string) {
			if _, err := os.Stat(filepath.Join(path, "artisan")); err == nil {
				return true, "Var"
			}
			return false, ""
		}},
		// PHP (Generic)
		{Type: domain.TypePHP, IsFrontend: false, CheckFunc: func(path string) (bool, string) {
			// Laravel değilse ve composer.json varsa
			if _, err := os.Stat(filepath.Join(path, "artisan")); err == nil {
				return false, "" // Laravel olarak algılansın
			}
			if _, err := os.Stat(filepath.Join(path, "composer.json")); err == nil {
				return true, "Var"
			}
			return false, ""
		}},
		// Spring (Java)
		{Type: domain.TypeSpring, IsFrontend: false, CheckFunc: func(path string) (bool, string) {
			// Check pom.xml or build.gradle for spring-boot
			pomPath := filepath.Join(path, "pom.xml")
			gradlePath := filepath.Join(path, "build.gradle")
			if data, err := os.ReadFile(pomPath); err == nil {
				if strings.Contains(string(data), "spring-boot") {
					return true, "Var"
				}
			}
			if data, err := os.ReadFile(gradlePath); err == nil {
				if strings.Contains(string(data), "spring-boot") {
					return true, "Var"
				}
			}
			return false, ""
		}},
	}
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

		p := domain.Project{
			Name: entry.Name(),
			Path: fullPath,
			Type: domain.TypeUnknown,
		}

		// ========================================
		// ADIM 1: Alt klasörleri tara (Signature-Based)
		// ========================================
		s.scanSubdirectories(fullPath, &p)

		// ========================================
		// ADIM 2: Root dizini tara (Monorepo olmayan projeler)
		// ========================================
		if !p.HasFrontend && !p.HasBackend {
			s.scanRootDirectory(fullPath, &p)
		}

		// ========================================
		// ADIM 3: Tip Belirleme
		// ========================================
		s.determineProjectType(&p)

		// ========================================
		// ADIM 4: Sadece proje olarak algılananları ekle
		// ========================================
		isProject := p.HasFrontend || p.HasBackend || p.HasDocker ||
			p.Type != domain.TypeUnknown ||
			p.FrontendType != "" || p.BackendType != ""

		if isProject {
			results = append(results, p)
		}
	}
	return results
}

// scanSubdirectories scans all immediate subdirectories for tech signatures
func (s *Scanner) scanSubdirectories(projectPath string, p *domain.Project) {
	entries, err := os.ReadDir(projectPath)
	if err != nil {
		return
	}

	frontendSigs := s.getFrontendSignatures()
	backendSigs := s.getBackendSignatures()

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		subPath := filepath.Join(projectPath, entry.Name())

		// Skip common non-project directories
		name := strings.ToLower(entry.Name())
		if name == "node_modules" || name == ".git" || name == "dist" || name == "build" || name == ".next" {
			continue
		}

		// Check for Frontend signatures
		if !p.HasFrontend {
			for _, sig := range frontendSigs {
				if found, ver := sig.CheckFunc(subPath); found {
					p.HasFrontend = true
					p.FrontendPath = subPath
					p.FrontendVer = ver
					p.FrontendType = sig.Type
					if p.FrontendVer == "" {
						p.FrontendVer = "Var"
					}
					p.FrontendCmd = s.detectStartCommand(subPath)
					break
				}
			}
		}

		// Check for Backend signatures
		if !p.HasBackend {
			for _, sig := range backendSigs {
				if found, ver := sig.CheckFunc(subPath); found {
					p.HasBackend = true
					p.BackendPath = subPath
					p.BackendVer = ver
					p.BackendType = sig.Type
					if p.BackendVer == "" {
						p.BackendVer = "Var"
					}
					p.BackendCmd = s.detectStartCommand(subPath)
					break
				}
			}
		}

		// Docker check
		if _, err := os.Stat(filepath.Join(subPath, "Dockerfile")); err == nil {
			p.HasDocker = true
		}
		if _, err := os.Stat(filepath.Join(subPath, "docker-compose.yml")); err == nil {
			p.HasDocker = true
		}
		if _, err := os.Stat(filepath.Join(subPath, "docker-compose.yaml")); err == nil {
			p.HasDocker = true
		}
	}
}

// scanRootDirectory checks the project root for tech signatures (single-folder projects)
func (s *Scanner) scanRootDirectory(projectPath string, p *domain.Project) {
	frontendSigs := s.getFrontendSignatures()
	backendSigs := s.getBackendSignatures()

	// Check Frontend first
	for _, sig := range frontendSigs {
		if found, ver := sig.CheckFunc(projectPath); found {
			p.HasFrontend = true
			p.FrontendPath = projectPath
			p.FrontendVer = ver
			p.FrontendType = sig.Type
			if p.FrontendVer == "" {
				p.FrontendVer = "Var"
			}
			p.FrontendCmd = s.detectStartCommand(projectPath)
			p.Type = sig.Type
			break
		}
	}

	// Check Backend
	for _, sig := range backendSigs {
		if found, ver := sig.CheckFunc(projectPath); found {
			p.HasBackend = true
			p.BackendPath = projectPath
			p.BackendVer = ver
			p.BackendType = sig.Type
			if p.BackendVer == "" {
				p.BackendVer = "Var"
			}
			p.BackendCmd = s.detectStartCommand(projectPath)
			if p.Type == domain.TypeUnknown {
				p.Type = sig.Type
			}
			break
		}
	}

	// Docker check for root
	if _, err := os.Stat(filepath.Join(projectPath, "Dockerfile")); err == nil {
		p.HasDocker = true
	}
	if _, err := os.Stat(filepath.Join(projectPath, "docker-compose.yml")); err == nil {
		p.HasDocker = true
	}
	if _, err := os.Stat(filepath.Join(projectPath, "docker-compose.yaml")); err == nil {
		p.HasDocker = true
	}
}

// determineProjectType sets the main project type based on detected components
func (s *Scanner) determineProjectType(p *domain.Project) {
	if p.Type != domain.TypeUnknown {
		return // Already set
	}

	// If both frontend and backend exist, keep as Unknown (Fullstack)
	if p.HasFrontend && p.HasBackend {
		p.Type = domain.TypeUnknown
		return
	}

	// Detect from root if still unknown
	if !p.HasFrontend && !p.HasBackend {
		p.Type = domain.TypeUnknown
	}
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

// detectStartCommand determines the best command to start the project
func (s *Scanner) detectStartCommand(path string) string {
	// 1. JS/TS Projects (Next, Nest, React, Vue)
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

	// 3. Python/Django
	if _, err := os.Stat(filepath.Join(path, "manage.py")); err == nil {
		return "python manage.py runserver"
	}

	// 4. Python/Flask
	if _, err := os.Stat(filepath.Join(path, "app.py")); err == nil {
		return "flask run"
	}

	// 5. PHP/Laravel
	if _, err := os.Stat(filepath.Join(path, "artisan")); err == nil {
		return "php artisan serve"
	}

	// 6. Java/Spring (Maven)
	if _, err := os.Stat(filepath.Join(path, "pom.xml")); err == nil {
		return "mvn spring-boot:run"
	}

	// 7. Java/Spring (Gradle)
	if _, err := os.Stat(filepath.Join(path, "build.gradle")); err == nil {
		return "./gradlew bootRun"
	}

	// Default fallback
	return "echo [DevTerminal] Başlatma komutu bulunamadı"
}
