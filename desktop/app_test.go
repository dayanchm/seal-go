package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeProjectFindsUnclosedHTTPBody(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module sample\n\ngo 1.22\n")
	writeFile(t, dir, "main.go", `package main

import "net/http"

func main() {
	resp, _ := http.Get("https://example.com")
	_ = resp
}
`)

	findings, err := NewApp().AnalyzeProject(dir)
	if err != nil {
		t.Fatalf("AnalyzeProject returned error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("got %d findings, want 1: %#v", len(findings), findings)
	}
	if findings[0].Line != 6 {
		t.Fatalf("got line %d, want 6", findings[0].Line)
	}
}

func TestAnalyzeProjectAcceptsGoFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module sample\n\ngo 1.22\n")
	writeFile(t, dir, "main.go", `package main

import "net/http"

func main() {
	resp, _ := http.Get("https://example.com")
	_ = resp
}
`)

	findings, err := NewApp().AnalyzeProject(filepath.Join(dir, "main.go"))
	if err != nil {
		t.Fatalf("AnalyzeProject returned error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("got %d findings, want 1: %#v", len(findings), findings)
	}
}

func writeFile(t *testing.T, dir string, name string, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestAnalyzeProjectStandaloneGoFile(t *testing.T) {
	for _, name := range []string{"main.go", "_example.go", ".example.go"} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, dir, name, `package main
import "net/http"
func main() {
 resp, _ := http.Get("https://example.com")
 _ = resp
}
`)
			findings, err := NewApp().AnalyzeProject(filepath.Join(dir, name))
			if err != nil {
				t.Fatal(err)
			}
			if len(findings) != 1 || findings[0].Line != 4 {
				t.Fatalf("expected HTTP cleanup finding at line 4, got %#v", findings)
			}
			if findings[0].File != filepath.Join(dir, name) {
				t.Fatalf("finding refers to wrong file: %s", findings[0].File)
			}
		})
	}
}
