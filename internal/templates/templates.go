// Package templates provides access to the embedded .dockerignore boilerplates.
package templates

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

const ext = ".dockerignore"

// FS holds the embedded .dockerignore templates.
//
//go:embed *.dockerignore
var FS embed.FS

// validate rejects names that could escape the embedded FS root.
func validate(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("template name is empty")
	}
	if strings.ContainsAny(name, `/\`) || strings.Contains(name, "..") {
		return fmt.Errorf("invalid template name: %q", name)
	}
	return nil
}

// Read returns the content of a template by name (case-insensitive) along with
// its canonical name. It never reads outside the embedded template set.
func Read(name string) (content []byte, canonical string, err error) {
	if err = validate(name); err != nil {
		return nil, "", err
	}
	entries, err := fs.ReadDir(FS, ".")
	if err != nil {
		return nil, "", fmt.Errorf("read templates: %w", err)
	}
	for _, e := range entries {
		base := strings.TrimSuffix(e.Name(), ext)
		if strings.EqualFold(base, name) {
			b, rerr := FS.ReadFile(e.Name())
			if rerr != nil {
				return nil, "", fmt.Errorf("read %q: %w", base, rerr)
			}
			return b, base, nil
		}
	}
	return nil, "", fmt.Errorf("template %q not found", name)
}

// List returns all available template names, sorted.
func List() ([]string, error) {
	entries, err := fs.ReadDir(FS, ".")
	if err != nil {
		return nil, fmt.Errorf("read templates: %w", err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ext) {
			names = append(names, strings.TrimSuffix(e.Name(), ext))
		}
	}
	sort.Strings(names)
	return names, nil
}

func Combine(names []string) (result string, missing []string, err error) {
	var b strings.Builder
	seen := make(map[string]struct{})
	found := 0

	for _, name := range names {
		content, canonical, rerr := Read(name)
		if rerr != nil {
			missing = append(missing, name)
			continue
		}
		fmt.Fprintf(&b, "### %s ###\n", canonical)
		for _, line := range strings.Split(string(content), "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				b.WriteString(line + "\n")
				continue
			}
			if _, dup := seen[trimmed]; dup {
				continue
			}
			seen[trimmed] = struct{}{}
			b.WriteString(line + "\n")
		}
		b.WriteString("\n")
		found++
	}

	if found == 0 {
		return "", missing, fmt.Errorf("no valid templates found")
	}
	return b.String(), missing, nil
}
