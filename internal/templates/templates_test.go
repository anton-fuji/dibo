package templates

import (
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	names, err := List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(names) == 0 {
		t.Fatal("List() returned no templates")
	}
	// sorted
	for i := 1; i < len(names); i++ {
		if names[i-1] > names[i] {
			t.Errorf("List() not sorted: %v", names)
			break
		}
	}
	// expected entries present
	want := map[string]bool{"Common": false, "Go": false, "Node": false, "Secrets": false}
	for _, n := range names {
		if _, ok := want[n]; ok {
			want[n] = true
		}
	}
	for k, found := range want {
		if !found {
			t.Errorf("List() missing %q", k)
		}
	}
}

func TestRead(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		wantCan string
		wantErr bool
	}{
		{"exact", "Go", "Go", false},
		{"lowercase", "go", "Go", false},
		{"uppercase", "NODE", "Node", false},
		{"unknown", "nope", "", true},
		{"empty", "", "", true},
		{"traversal slash", "../templates", "", true},
		{"traversal dots", "..", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, canon, err := Read(tc.arg)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("Read(%q) expected error, got nil", tc.arg)
				}
				return
			}
			if err != nil {
				t.Fatalf("Read(%q) unexpected error: %v", tc.arg, err)
			}
			if canon != tc.wantCan {
				t.Errorf("Read(%q) canonical = %q, want %q", tc.arg, canon, tc.wantCan)
			}
		})
	}
}

func TestCombine(t *testing.T) {
	// missing names are reported, not fatal, as long as one resolves
	out, missing, err := Combine([]string{"Common", "nope"})
	if err != nil {
		t.Fatalf("Combine error: %v", err)
	}
	if len(missing) != 1 || missing[0] != "nope" {
		t.Errorf("missing = %v, want [nope]", missing)
	}
	if !strings.Contains(out, "### Common ###") {
		t.Errorf("output missing canonical header:\n%s", out)
	}

	// all-missing is an error
	if _, _, err := Combine([]string{"nope", "alsono"}); err == nil {
		t.Error("Combine(all-missing) expected error, got nil")
	}

	// de-duplication across templates (Secrets and Node both carry .env)
	out, _, err = Combine([]string{"Secrets", "Node"})
	if err != nil {
		t.Fatal(err)
	}
	if c := strings.Count("\n"+out, "\n.env\n"); c != 1 {
		t.Errorf(".env appears %d times after dedup, want 1", c)
	}

	// comments are preserved (not de-duplicated)
	if !strings.Contains(out, "# --- Dependencies ---") {
		t.Errorf("expected comment headers to be preserved")
	}
}
