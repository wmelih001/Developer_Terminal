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
	FrontendVer  string
	BackendVer   string
	FrontendCmd  string
	BackendCmd   string
	FrontendPath string
	BackendPath  string
}

// ProjectType teknoloji yığınını tanımlar
type ProjectType string

const (
	TypeReact   ProjectType = "React"
	TypeNext    ProjectType = "Next.js"
	TypeNest    ProjectType = "NestJS"
	TypeGo      ProjectType = "Go"
	TypeUnknown ProjectType = "Bilinmeyen"
)
