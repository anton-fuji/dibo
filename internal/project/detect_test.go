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
