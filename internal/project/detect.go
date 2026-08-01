// Package project detects common project types from files in a directory.
package project

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Type describes a dibo template and the files that identify it.
type Type struct {
	Template string
	Signals  []string
}

var types = []Type{
	{Template: "Go", Signals: []string{"go.mod"}},
	{Template: "Node", Signals: []string{"package.json"}},
	{Template: "Python", Signals: []string{"pyproject.toml", "requirements.txt", "Pipfile", "setup.py"}},
	{Template: "Ruby", Signals: []string{"Gemfile", ".ruby-version"}},
	{Template: "Rust", Signals: []string{"Cargo.toml"}},
	{Template: "Java", Signals: []string{"pom.xml", "build.gradle", "build.gradle.kts", "settings.gradle", "settings.gradle.kts"}},
	{Template: "PHP", Signals: []string{"composer.json"}},
}

// Detect returns templates whose identifying files occur at the project root.
// .NET is detected from root-level solution and project files as their names vary.
func Detect(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read project directory %s: %w", dir, err)
	}

	present := make(map[string]bool, len(entries))
	dotNet := false
	for _, entry := range entries {
		name := entry.Name()
		present[name] = true
		if filepath.Ext(name) == ".sln" || filepath.Ext(name) == ".csproj" || filepath.Ext(name) == ".fsproj" {
			dotNet = true
		}
	}

	result := make([]string, 0, len(types)+1)
	for _, typ := range types {
		for _, signal := range typ.Signals {
			if present[signal] {
				result = append(result, typ.Template)
				break
			}
		}
	}
	if dotNet {
		result = append(result, "dotNet")
	}
	sort.Strings(result)
	return result, nil
}
