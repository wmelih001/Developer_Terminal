package service

import (
	"bytes"
	"fmt"
	"os/exec"
	"text/template"

	"devterminal/pkg/domain"
)

type Launcher struct {
	Config *domain.Config
}

func NewLauncher(cfg *domain.Config) *Launcher {
	return &Launcher{Config: cfg}
}

// LaunchProject opens the project in Windows Terminal using the configured template
func (l *Launcher) LaunchProject(p *domain.Project, mode string) error {
	var cmdTmpl string
	switch mode {
	case "frontend":
		cmdTmpl = l.Config.Commands.LaunchFrontend
	case "backend":
		cmdTmpl = l.Config.Commands.LaunchBackend
	case "full":
		cmdTmpl = l.Config.Commands.LaunchFull
	default:
		return fmt.Errorf("unknown mode: %s", mode)
	}

	// Parse template
	tmpl, err := template.New("cmd").Parse(cmdTmpl)
	if err != nil {
		return err
	}

	var cmdStr bytes.Buffer
	if err := tmpl.Execute(&cmdStr, p); err != nil {
		return err
	}

	// Parse the command string into executable and args
	// We need to support quoted arguments.
	expandedCmd := cmdStr.String()
	args := parseArgs(expandedCmd)

	if len(args) == 0 {
		return fmt.Errorf("empty command")
	}

	// Execute directly, bypassing cmd /C
	// args[0] is executable (e.g. wt.exe), args[1:] are arguments
	c := exec.Command(args[0], args[1:]...)
	return c.Start()
}

// parseArgs splits a string into arguments, respecting quotes
func parseArgs(cmd string) []string {
	var args []string
	var current []rune
	inQuote := false
	quoteChar := rune(0)

	for _, r := range cmd {
		if inQuote {
			if r == quoteChar {
				inQuote = false
			} else {
				current = append(current, r)
			}
		} else {
			switch r {
			case '"', '\'':
				inQuote = true
				quoteChar = r
			case ' ', '\t':
				if len(current) > 0 {
					args = append(args, string(current))
					current = nil
				}
			default:
				current = append(current, r)
			}
		}
	}
	if len(current) > 0 {
		args = append(args, string(current))
	}
	return args
}

// LaunchPrisma opens Prisma Studio for the project
func (l *Launcher) LaunchPrisma(p *domain.Project) error {
	path := p.Path
	if p.PrismaPath != "" {
		path = p.PrismaPath
	}
	cmdStr := fmt.Sprintf(`wt -w 0 nt --title "Prisma Studio" -d "%s" cmd /k "npx prisma studio"`, path)

	args := parseArgs(cmdStr)
	if len(args) == 0 {
		return fmt.Errorf("failed to create prisma command")
	}

	c := exec.Command(args[0], args[1:]...)
	return c.Start()
}

// LaunchDrizzle opens Drizzle Studio
func (l *Launcher) LaunchDrizzle(p *domain.Project) error {
	path := p.Path
	if p.DrizzlePath != "" {
		path = p.DrizzlePath
	}
	cmdStr := fmt.Sprintf(`wt -w 0 nt --title "Drizzle Studio" -d "%s" cmd /k "npx drizzle-kit studio"`, path)
	return l.runCmd(cmdStr)
}

// LaunchHasura opens Hasura Console
func (l *Launcher) LaunchHasura(p *domain.Project) error {
	path := p.Path
	if p.HasuraPath != "" {
		path = p.HasuraPath
	}
	cmdStr := fmt.Sprintf(`wt -w 0 nt --title "Hasura Console" -d "%s" cmd /k "hasura console"`, path)
	return l.runCmd(cmdStr)
}

// LaunchSupabase opens Supabase Dashboard
func (l *Launcher) LaunchSupabase(p *domain.Project) error {
	path := p.Path
	if p.SupabasePath != "" {
		path = p.SupabasePath
	}
	cmdStr := fmt.Sprintf(`wt -w 0 nt --title "Supabase Status" -d "%s" cmd /k "npx supabase status"`, path)
	return l.runCmd(cmdStr)
}

// LaunchStorybook opens Storybook
func (l *Launcher) LaunchStorybook(p *domain.Project) error {
	path := p.Path
	if p.StorybookPath != "" {
		path = p.StorybookPath
	}
	cmdStr := fmt.Sprintf(`wt -w 0 nt --title "Storybook" -d "%s" cmd /k "npm run storybook"`, path)
	return l.runCmd(cmdStr)
}

// runCmd helper
func (l *Launcher) runCmd(cmdStr string) error {
	args := parseArgs(cmdStr)
	if len(args) == 0 {
		return fmt.Errorf("failed to create command")
	}
	c := exec.Command(args[0], args[1:]...)
	return c.Start()
}
