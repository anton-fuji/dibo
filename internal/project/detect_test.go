package project

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDetect(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"go.mod", "package.json", "service.csproj"} {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := Detect(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Go", "Node", "dotNet"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Detect() = %v, want %v", got, want)
	}
}

func TestDetectWithEvidence(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"go.mod", "pyproject.toml", "requirements.txt", "service.csproj"} {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "nested", "Cargo.toml"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := DetectWithEvidence(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := []Detection{
		{Template: "Go", Signals: []string{"go.mod"}},
		{Template: "Python", Signals: []string{"pyproject.toml", "requirements.txt"}},
		{Template: "dotNet", Signals: []string{"service.csproj"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DetectWithEvidence() = %v, want %v", got, want)
	}
}

func TestDetectWithEvidenceForEachType(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		template string
	}{
		{name: "go", filename: "go.mod", template: "Go"},
		{name: "node", filename: "package.json", template: "Node"},
		{name: "python", filename: "pyproject.toml", template: "Python"},
		{name: "ruby", filename: "Gemfile", template: "Ruby"},
		{name: "rust", filename: "Cargo.toml", template: "Rust"},
		{name: "java", filename: "pom.xml", template: "Java"},
		{name: "php", filename: "composer.json", template: "PHP"},
		{name: "dotnet", filename: "service.csproj", template: "dotNet"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, tt.filename), nil, 0o644); err != nil {
				t.Fatal(err)
			}

			got, err := DetectWithEvidence(dir)
			if err != nil {
				t.Fatal(err)
			}
			want := []Detection{{Template: tt.template, Signals: []string{tt.filename}}}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("DetectWithEvidence() = %v, want %v", got, want)
			}
		})
	}
}

func TestDetectWithEvidenceNoMatch(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "README.md"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "nested", "go.mod"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := DetectWithEvidence(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("DetectWithEvidence() = %v, want no detections", got)
	}
}

func TestDetectWithEvidenceRecursive(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"go.mod":                            "",
		"apps/api/requirements.txt":         "",
		"packages/web/package.json":         "",
		"node_modules/ignored/package.json": "",
		"target/ignored/Cargo.toml":         "",
	}
	for name := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	got, err := DetectWithEvidenceRecursive(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := []Detection{
		{Template: "Go", Signals: []string{"go.mod"}, Directory: "."},
		{Template: "Python", Signals: []string{"requirements.txt"}, Directory: "apps/api"},
		{Template: "Node", Signals: []string{"package.json"}, Directory: "packages/web"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DetectWithEvidenceRecursive() = %v, want %v", got, want)
	}
}
