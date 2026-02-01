package inference

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GoDetector detects Go projects.
type GoDetector struct{}

func (d *GoDetector) Name() string  { return "go" }
func (d *GoDetector) Priority() int { return 80 }

func (d *GoDetector) Detect(projectPath string) bool {
	return FileExistsInProject(projectPath, "go.mod")
}

// GoInferrer infers configuration from Go projects.
type GoInferrer struct{}

func (i *GoInferrer) Language() string { return "go" }

func (i *GoInferrer) Infer(projectPath string) (*ProjectInfo, error) {
	info := &ProjectInfo{
		Language:       "go",
		PackageManager: "go",
		InstallCommand: "go mod download",
		Scripts:        make(map[string]string),
	}

	// Parse go.mod for version
	goModPath := filepath.Join(projectPath, "go.mod")
	if version := parseGoModVersion(goModPath); version != "" {
		info.Runtime = "go" + version
	}

	// Check for Makefile targets
	makefilePath := filepath.Join(projectPath, "Makefile")
	if FileExists(makefilePath) {
		targets := parseMakefileTargets(makefilePath)
		info.Scripts = targets

		// Use Makefile targets for commands if available
		if _, ok := targets["dev"]; ok {
			info.DevCommand = "make dev"
		}
		if _, ok := targets["build"]; ok {
			info.BuildCommand = "make build"
		}
		if _, ok := targets["run"]; ok {
			info.StartCommand = "make run"
		} else if _, ok := targets["start"]; ok {
			info.StartCommand = "make start"
		}
	}

	// Default commands if not found in Makefile
	if info.DevCommand == "" {
		info.DevCommand = "go run ."
	}
	if info.BuildCommand == "" {
		info.BuildCommand = "go build -o bin/app ."
	}
	if info.StartCommand == "" {
		info.StartCommand = "./bin/app"
	}

	// Default port for Go web apps
	info.Port = 8080

	return info, nil
}

// parseGoModVersion extracts the Go version from go.mod.
func parseGoModVersion(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`^go\s+(\d+\.\d+)`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if matches := re.FindStringSubmatch(line); len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

// parseMakefileTargets extracts target names from a Makefile.
func parseMakefileTargets(path string) map[string]string {
	targets := make(map[string]string)

	file, err := os.Open(path)
	if err != nil {
		return targets
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_-]*)\s*:`)

	for scanner.Scan() {
		line := scanner.Text()
		if matches := re.FindStringSubmatch(line); len(matches) > 1 {
			target := matches[1]
			// Skip common non-command targets
			if target != ".PHONY" && target != ".DEFAULT" {
				targets[target] = target
			}
		}
	}

	return targets
}

// Go Framework Detectors

// GinDetector detects Gin framework projects.
type GinDetector struct{}

func (d *GinDetector) Language() string { return "go" }
func (d *GinDetector) Name() string     { return "gin" }

func (d *GinDetector) Detect(projectPath string, info *ProjectInfo) bool {
	// Check go.mod for gin import
	goModPath := filepath.Join(projectPath, "go.mod")
	return fileContains(goModPath, "github.com/gin-gonic/gin")
}

func (d *GinDetector) Defaults() *FrameworkDefaults {
	return &FrameworkDefaults{
		Framework: "gin",
		Language:  "go",
		Dev:       "go run .",
		Build:     "go build -o bin/app .",
		Start:     "./bin/app",
		Port:      8080,
	}
}

// EchoDetector detects Echo framework projects.
type EchoDetector struct{}

func (d *EchoDetector) Language() string { return "go" }
func (d *EchoDetector) Name() string     { return "echo" }

func (d *EchoDetector) Detect(projectPath string, info *ProjectInfo) bool {
	goModPath := filepath.Join(projectPath, "go.mod")
	return fileContains(goModPath, "github.com/labstack/echo")
}

func (d *EchoDetector) Defaults() *FrameworkDefaults {
	return &FrameworkDefaults{
		Framework: "echo",
		Language:  "go",
		Dev:       "go run .",
		Build:     "go build -o bin/app .",
		Start:     "./bin/app",
		Port:      8080,
	}
}

// FiberDetector detects Fiber framework projects.
type FiberDetector struct{}

func (d *FiberDetector) Language() string { return "go" }
func (d *FiberDetector) Name() string     { return "fiber" }

func (d *FiberDetector) Detect(projectPath string, info *ProjectInfo) bool {
	goModPath := filepath.Join(projectPath, "go.mod")
	return fileContains(goModPath, "github.com/gofiber/fiber")
}

func (d *FiberDetector) Defaults() *FrameworkDefaults {
	return &FrameworkDefaults{
		Framework: "fiber",
		Language:  "go",
		Dev:       "go run .",
		Build:     "go build -o bin/app .",
		Start:     "./bin/app",
		Port:      3000,
	}
}

// ChiDetector detects Chi router projects.
type ChiDetector struct{}

func (d *ChiDetector) Language() string { return "go" }
func (d *ChiDetector) Name() string     { return "chi" }

func (d *ChiDetector) Detect(projectPath string, info *ProjectInfo) bool {
	goModPath := filepath.Join(projectPath, "go.mod")
	return fileContains(goModPath, "github.com/go-chi/chi")
}

func (d *ChiDetector) Defaults() *FrameworkDefaults {
	return &FrameworkDefaults{
		Framework: "chi",
		Language:  "go",
		Dev:       "go run .",
		Build:     "go build -o bin/app .",
		Start:     "./bin/app",
		Port:      8080,
	}
}

// fileContains checks if a file contains a specific string.
func fileContains(path, substr string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), substr)
}

// RegisterGo registers Go language support with the registry.
func RegisterGo(r *Registry) {
	r.RegisterLanguage(&GoDetector{}, &GoInferrer{})

	// Framework detectors
	r.RegisterFramework(&GinDetector{})
	r.RegisterFramework(&EchoDetector{})
	r.RegisterFramework(&FiberDetector{})
	r.RegisterFramework(&ChiDetector{})
}
