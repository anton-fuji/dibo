package cmd

import (
	"runtime/debug"
	"testing"
)

func TestResolvedVersionUsesReleaseVersion(t *testing.T) {
	originalVersion := version
	originalReadBuildInfo := readBuildInfo
	t.Cleanup(func() {
		version = originalVersion
		readBuildInfo = originalReadBuildInfo
	})

	version = "1.2.3"
	if got := resolvedVersion(); got != "1.2.3" {
		t.Errorf("resolvedVersion() = %q, want %q", got, "1.2.3")
	}
}

func TestResolvedVersionUsesGoInstallModuleVersion(t *testing.T) {
	originalVersion := version
	originalReadBuildInfo := readBuildInfo
	t.Cleanup(func() {
		version = originalVersion
		readBuildInfo = originalReadBuildInfo
	})

	version = "dev"
	readBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Main: debug.Module{
				Path:    "github.com/anton-fuji/dibo",
				Version: "v0.4.0",
			},
		}, true
	}

	if got := resolvedVersion(); got != "v0.4.0" {
		t.Errorf("resolvedVersion() = %q, want %q", got, "v0.4.0")
	}
}
