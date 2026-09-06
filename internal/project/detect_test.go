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
