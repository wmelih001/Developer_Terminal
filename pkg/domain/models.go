package domain

// Config uygulama konfigürasyonunu tutar
type Config struct {
	ProjectsPaths []string `mapstructure:"projects_paths"`
	Commands      Commands `mapstructure:"commands"`
	IgnoredFiles  []string `mapstructure:"ignored_files"`
	NgrokPath     string   `mapstructure:"ngrok_path"`
}

// Commands özel komut şablonlarını tutar
type Commands struct {
	LaunchFrontend string `mapstructure:"launch_frontend"`
	LaunchBackend  string `mapstructure:"launch_backend"`
	LaunchFull     string `mapstructure:"launch_full"`
}

// Project diskteki bir geliştirici projesini temsil eder
type Project struct {
	Name         string
	Path         string
	Type         ProjectType
	Tags         []string
	HasFrontend  bool
	HasBackend   bool
	HasDocker    bool // Docker desteği var mı
	FrontendVer  string
	BackendVer   string
	FrontendType ProjectType // Algılanan frontend teknolojisi (Next.js, React, Vue vb.)
	BackendType  ProjectType // Algılanan backend teknolojisi (NestJS, Go, Django vb.)
	FrontendCmd  string
	BackendCmd   string
	FrontendPath string
	BackendPath  string
}

// ProjectType teknoloji yığınını tanımlar
type ProjectType string

const (
	// Frontend
	TypeReact       ProjectType = "React"
	TypeNext        ProjectType = "Next.js"
	TypeVue         ProjectType = "Vue"
	TypeVite        ProjectType = "Vite"
	TypeReactNative ProjectType = "React Native"
	TypeMobile      ProjectType = "Mobile"
	TypeHTML        ProjectType = "HTML"
	TypeTypeScript  ProjectType = "TypeScript"

	// Backend
	TypeNest    ProjectType = "NestJS"
	TypeExpress ProjectType = "Express"
	TypeGo      ProjectType = "Go"
	TypeDjango  ProjectType = "Django"
	TypeFlask   ProjectType = "Flask"
	TypeLaravel ProjectType = "Laravel"
	TypeSpring  ProjectType = "Spring"
	TypePHP     ProjectType = "PHP"

	// Infrastructure
	TypeDocker ProjectType = "Docker"

	// Unknown
	TypeUnknown ProjectType = "Bilinmeyen"
)
