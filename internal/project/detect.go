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

// Detection describes a detected project type and the files that identified it.
type Detection struct {
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
	detected, err := DetectWithEvidence(dir)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(detected))
	for _, detection := range detected {
		result = append(result, detection.Template)
	}
	return result, nil
}

// DetectWithEvidence returns detected project types and the root-level files
// that identified each type.
func DetectWithEvidence(dir string) ([]Detection, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read project directory %s: %w", dir, err)
	}

	present := make(map[string]bool, len(entries))
	dotNetSignals := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		present[name] = true
		if filepath.Ext(name) == ".sln" || filepath.Ext(name) == ".csproj" || filepath.Ext(name) == ".fsproj" {
			dotNetSignals = append(dotNetSignals, name)
		}
	}

	result := make([]Detection, 0, len(types)+1)
	for _, typ := range types {
		signals := make([]string, 0, len(typ.Signals))
		for _, signal := range typ.Signals {
			if present[signal] {
				signals = append(signals, signal)
			}
		}
		if len(signals) > 0 {
			result = append(result, Detection{Template: typ.Template, Signals: signals})
		}
	}
	if len(dotNetSignals) > 0 {
		sort.Strings(dotNetSignals)
		result = append(result, Detection{Template: "dotNet", Signals: dotNetSignals})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Template < result[j].Template
	})
	return result, nil
}
