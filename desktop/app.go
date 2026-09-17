package main

import (
	"context"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
	"seal-go/analyzer"
)

type Finding struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Message string `json:"message"`
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SelectFolder() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choose a Go project",
	})
}

func (a *App) SelectGoFile() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choose a Go file",
		Filters: []runtime.FileFilter{
			{DisplayName: "Go Files (*.go)", Pattern: "*.go"},
		},
	})
}

func (a *App) AnalyzeProject(path string) ([]Finding, error) {
	if path == "" {
		return nil, fmt.Errorf("choose a project folder or Go file first")
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open selected path: %w", err)
	}

	dir := path
	pattern := "./..."
	selectedFile := ""
	if !info.IsDir() {
		if !strings.EqualFold(filepath.Ext(path), ".go") {
			return nil, fmt.Errorf("selected file is not a Go file")
		}
		selectedFile, err = filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("cannot resolve selected file: %w", err)
		}
		selectedFile = filepath.Clean(selectedFile)
		dir = filepath.Dir(selectedFile)
		pattern = "file=" + selectedFile
	}

	fset := token.NewFileSet()
	cfg := &packages.Config{
		Context: a.ctx,
		Dir:     dir,
		Fset:    fset,
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes,
	}

	pkgs, err := packages.Load(cfg, pattern)
	loadedFile := selectedFile
	// A file query can return no package for files excluded by Go's
	// directory matching rules. Load the explicitly selected file instead.
	if err == nil && len(pkgs) == 0 && selectedFile != "" {
		// Use an ordinary filename so names beginning with '_' or '.'
		// can still be inspected when the user explicitly selects them.
		source, readErr := os.ReadFile(selectedFile)
		if readErr != nil {
			return nil, fmt.Errorf("cannot read selected Go file: %w", readErr)
		}
		tempDir, tempErr := os.MkdirTemp("", "seal-go-")
		if tempErr != nil {
			return nil, fmt.Errorf("cannot prepare selected Go file: %w", tempErr)
		}
		defer os.RemoveAll(tempDir)
		loadedFile = filepath.Join(tempDir, "selected.go")
		if writeErr := os.WriteFile(loadedFile, source, 0o600); writeErr != nil {
			return nil, fmt.Errorf("cannot prepare selected Go file: %w", writeErr)
		}
		pkgs, err = packages.Load(cfg, loadedFile)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot load Go packages: %w", err)
	}
	if len(pkgs) == 0 {
		if selectedFile != "" {
			return nil, fmt.Errorf("cannot load selected Go file: %s", selectedFile)
		}
		return nil, fmt.Errorf("no Go packages found in selected folder")
	}

	findings := make([]Finding, 0)
	for _, pkg := range pkgs {
		for _, loadErr := range pkg.Errors {
			findings = append(findings, Finding{
				File:    loadErr.Pos,
				Message: loadErr.Msg,
			})
		}

		if len(pkg.Syntax) == 0 || pkg.TypesInfo == nil {
			continue
		}

		pass := &analysis.Pass{
			Analyzer:   analyzer.Analyzer,
			Fset:       fset,
			Files:      pkg.Syntax,
			Pkg:        pkg.Types,
			TypesInfo:  pkg.TypesInfo,
			TypesSizes: pkg.TypesSizes,
			ResultOf:   map[*analysis.Analyzer]any{},
			ReadFile:   os.ReadFile,
		}
		pass.Report = func(diagnostic analysis.Diagnostic) {
			position := fset.Position(diagnostic.Pos)
			if loadedFile != selectedFile && sameFile(position.Filename, loadedFile) {
				position.Filename = selectedFile
			}
			if selectedFile != "" && !sameFile(position.Filename, selectedFile) {
				return
			}
			findings = append(findings, Finding{
				File:    position.Filename,
				Line:    position.Line,
				Message: diagnostic.Message,
			})
		}

		if _, err := analyzer.Analyzer.Run(pass); err != nil {
			return nil, fmt.Errorf("analysis failed for package %s: %w", pkg.PkgPath, err)
		}
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].File != findings[j].File {
			return findings[i].File < findings[j].File
		}
		if findings[i].Line != findings[j].Line {
			return findings[i].Line < findings[j].Line
		}
		return findings[i].Message < findings[j].Message
	})

	return findings, nil
}

func sameFile(left string, right string) bool {
	absolute, err := filepath.Abs(left)
	if err != nil {
		return filepath.Clean(left) == right
	}
	return filepath.Clean(absolute) == right
}
