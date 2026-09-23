package backend_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const root = "github.com/franciscoHonorat/Sys-Called/backend"

// TestModuleBoundaries keeps the monolith modular. The Go compiler already
// stops a module from importing another module's internal/ packages; this test
// adds the rules the compiler cannot express:
//
//   - a module talks to another one only through its contracts package, never
//     through its facade (the package that builds it);
//   - a contracts package depends on nothing but the standard library, so it
//     never drags a module's internals or a framework along;
//   - the platform packages are shared infrastructure and know no module.
//
// The layers inside each module are tested by modules/<name>/internal/architecture_test.go.
func TestModuleBoundaries(t *testing.T) {
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		from := filepath.ToSlash(filepath.Dir(path))
		for _, spec := range file.Imports {
			imp, _ := strconv.Unquote(spec.Path.Value)
			if reason := boundaryViolation(from, imp); reason != "" {
				t.Errorf("%s imports %s: %s", path, imp, reason)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBoundaryViolations(t *testing.T) {
	cases := []struct {
		from, imp string
		allowed   bool
	}{
		{"modules/tickets/internal/adapters/out/employees", root + "/modules/employees/contracts", true},
		{"modules/tickets/internal/adapters/out/employees", root + "/modules/employees", false},
		{"modules/tickets", root + "/modules/employees/contracts", true},
		{"modules/tickets/internal/adapters/in/events", root + "/internal/platform/eventbus", true},
		{"modules/employees/contracts", "context", true},
		{"modules/employees/contracts", "github.com/google/uuid", false},
		{"modules/employees/contracts", root + "/internal/platform/eventbus", false},
		{"internal/platform/httpserver", root + "/modules/tickets", false},
		{"internal/platform/httpserver", "github.com/gin-gonic/gin", true},
		{"internal/app", root + "/modules/tickets", true},
		{"cmd/sys-called", root + "/internal/app", true},
	}

	for _, c := range cases {
		got := boundaryViolation(c.from, c.imp) == ""
		if got != c.allowed {
			t.Errorf("%s -> %s: allowed=%v, want %v", c.from, c.imp, got, c.allowed)
		}
	}
}

func boundaryViolation(from, imp string) string {
	fromModule := moduleOf(from)

	if strings.HasPrefix(from, "modules/") && strings.HasSuffix(from, "/contracts") {
		if !isStandardLibrary(imp) {
			return "a contracts package may only use the standard library"
		}
		return ""
	}

	if strings.HasPrefix(from, "internal/platform") && strings.HasPrefix(imp, root+"/modules/") {
		return "the platform must not know the modules"
	}

	target, ok := strings.CutPrefix(imp, root+"/")
	if !ok || fromModule == "" {
		return ""
	}
	if toModule := moduleOf(target); toModule != "" && toModule != fromModule {
		if target != "modules/"+toModule+"/contracts" {
			return "modules must talk to each other only through modules/" + toModule + "/contracts"
		}
	}
	return ""
}

// moduleOf returns the module a path belongs to ("tickets" for
// modules/tickets/internal/...), or "" outside modules/.
func moduleOf(path string) string {
	rest, ok := strings.CutPrefix(path, "modules/")
	if !ok {
		return ""
	}
	name, _, _ := strings.Cut(rest, "/")
	return name
}

func isStandardLibrary(imp string) bool {
	first, _, _ := strings.Cut(imp, "/")
	return !strings.Contains(first, ".")
}
