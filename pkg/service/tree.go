package service

import (
	"os"
	"path/filepath"
	"strings"

	"devterminal/pkg/domain"
)

// TreeGenerator generates an ASCII tree of the project
type TreeGenerator struct {
	Config *domain.Config
}

func NewTreeGenerator(cfg *domain.Config) *TreeGenerator {
	return &TreeGenerator{Config: cfg}
}

// GenerateTree returns a string representation of the file structure
func (t *TreeGenerator) GenerateTree(rootPath string) (string, error) {
	var sb strings.Builder
	sb.WriteString(rootPath + "\n")
	err := t.walk(rootPath, "", &sb)
	return sb.String(), err
}

func (t *TreeGenerator) walk(path string, prefix string, sb *strings.Builder) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Filter ignored files/dirs
	var filtered []os.DirEntry
	for _, entry := range entries {
		if !t.isIgnored(entry.Name()) {
			filtered = append(filtered, entry)
		}
	}

	for i, entry := range filtered {
		isLast := i == len(filtered)-1

		marker := "├── "
		if isLast {
			marker = "└── "
		}

		sb.WriteString(prefix + marker + entry.Name() + "\n")

		if entry.IsDir() {
			newPrefix := prefix + "│   "
			if isLast {
				newPrefix = prefix + "    "
			}
			t.walk(filepath.Join(path, entry.Name()), newPrefix, sb)
		}
	}
	return nil
}

func (t *TreeGenerator) isIgnored(name string) bool {
	for _, ignored := range t.Config.IgnoredFiles {
		// Simple match or glob could be implemented. For now, exact match or contains.
		// "node_modules" exact match
		// ".git" exact match
		if name == ignored {
			return true
		}
	}
	return false
}
