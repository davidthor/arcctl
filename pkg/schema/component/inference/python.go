package inference

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// PythonDetector detects Python projects.
type PythonDetector struct{}

func (d *PythonDetector) Name() string  { return "python" }
func (d *PythonDetector) Priority() int { return 70 }

func (d *PythonDetector) Detect(projectPath string) bool {
	return AnyFileExists(projectPath,
		"pyproject.toml",
		"requirements.txt",
		"setup.py",
		"Pipfile",
	)
}

// PythonInferrer infers configuration from Python projects.
type PythonInferrer struct{}

func (i *PythonInferrer) Language() string { return "python" }

func (i *PythonInferrer) Infer(projectPath string) (*ProjectInfo, error) {
	info := &ProjectInfo{
		Language:     "python",
		Dependencies: make(map[string]string),
		Scripts:      make(map[string]string),
		Port:         8000, // Common default for Python web apps
	}

	// Detect package manager and parse dependencies
	if FileExistsInProject(projectPath, "poetry.lock") {
		info.PackageManager = "poetry"
		info.InstallCommand = "poetry install"
	} else if FileExistsInProject(projectPath, "Pipfile.lock") {
		info.PackageManager = "pipenv"
		info.InstallCommand = "pipenv install"
	} else if FileExistsInProject(projectPath, "requirements.txt") {
		info.PackageManager = "pip"
		info.InstallCommand = "pip install -r requirements.txt"
		// Parse requirements.txt for dependencies
		parsePythonRequirements(filepath.Join(projectPath, "requirements.txt"), info)
	} else if FileExistsInProject(projectPath, "pyproject.toml") {
		info.PackageManager = "pip"
		info.InstallCommand = "pip install -e ."
	}

	// Try to parse pyproject.toml for more info
	pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
	if FileExists(pyprojectPath) {
		parsePyprojectToml(pyprojectPath, info)
	}

	// Check for .python-version file for runtime
	pythonVersionPath := filepath.Join(projectPath, ".python-version")
	if FileExists(pythonVersionPath) {
		if version := readFirstLine(pythonVersionPath); version != "" {
			info.Runtime = "python" + strings.TrimSpace(version)
		}
	}

	// Default handler for Lambda
	info.Handler = "handler.handler"

	return info, nil
}

// parsePythonRequirements parses requirements.txt for dependencies.
func parsePythonRequirements(path string, info *ProjectInfo) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`^([a-zA-Z0-9_-]+)`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "-") {
			continue
		}
		if matches := re.FindStringSubmatch(line); len(matches) > 1 {
			info.Dependencies[strings.ToLower(matches[1])] = ""
		}
	}
}

// parsePyprojectToml parses pyproject.toml for project metadata.
// This is a simplified parser that looks for common patterns.
func parsePyprojectToml(path string, info *ProjectInfo) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inDependencies := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check for Python version requirement
		if strings.HasPrefix(line, "requires-python") {
			if version := extractPythonVersion(line); version != "" {
				info.Runtime = "python" + version
			}
		}

		// Track dependencies section
		if strings.Contains(line, "[project.dependencies]") || strings.Contains(line, "[tool.poetry.dependencies]") {
			inDependencies = true
			continue
		}
		if strings.HasPrefix(line, "[") {
			inDependencies = false
		}

		if inDependencies && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				depName := strings.Trim(strings.TrimSpace(parts[0]), `"'`)
				info.Dependencies[strings.ToLower(depName)] = ""
			}
		}

		// Check for scripts/entry points
		if strings.HasPrefix(line, "dev = ") || strings.HasPrefix(line, `"dev" = `) {
			// Poetry script
			info.DevCommand = extractScriptCommand(line)
		}
	}
}

// extractPythonVersion extracts version from a requires-python line.
func extractPythonVersion(line string) string {
	re := regexp.MustCompile(`(\d+\.\d+)`)
	if matches := re.FindStringSubmatch(line); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractScriptCommand extracts the command from a script definition.
func extractScriptCommand(line string) string {
	re := regexp.MustCompile(`=\s*["'](.+)["']`)
	if matches := re.FindStringSubmatch(line); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// readFirstLine reads the first line of a file.
func readFirstLine(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

// Python Framework Detectors

// FastAPIDetector detects FastAPI projects.
type FastAPIDetector struct{}

func (d *FastAPIDetector) Language() string { return "python" }
func (d *FastAPIDetector) Name() string     { return "fastapi" }

func (d *FastAPIDetector) Detect(projectPath string, info *ProjectInfo) bool {
	return HasDependency(info.Dependencies, "fastapi")
}

func (d *FastAPIDetector) Defaults() *FrameworkDefaults {
	return &FrameworkDefaults{
		Framework: "fastapi",
		Language:  "python",
		Dev:       "uvicorn main:app --reload",
		Start:     "uvicorn main:app --host 0.0.0.0",
		Port:      8000,
		Handler:   "main.handler",
	}
}

// FlaskDetector detects Flask projects.
type FlaskDetector struct{}

func (d *FlaskDetector) Language() string { return "python" }
func (d *FlaskDetector) Name() string     { return "flask" }

func (d *FlaskDetector) Detect(projectPath string, info *ProjectInfo) bool {
	return HasDependency(info.Dependencies, "flask")
}

func (d *FlaskDetector) Defaults() *FrameworkDefaults {
	return &FrameworkDefaults{
		Framework: "flask",
		Language:  "python",
		Dev:       "flask run --reload",
		Start:     "gunicorn app:app",
		Port:      5000,
		Handler:   "app.handler",
	}
}

// DjangoDetector detects Django projects.
type DjangoDetector struct{}

func (d *DjangoDetector) Language() string { return "python" }
func (d *DjangoDetector) Name() string     { return "django" }

func (d *DjangoDetector) Detect(projectPath string, info *ProjectInfo) bool {
	return HasDependency(info.Dependencies, "django")
}

func (d *DjangoDetector) Defaults() *FrameworkDefaults {
	return &FrameworkDefaults{
		Framework: "django",
		Language:  "python",
		Dev:       "python manage.py runserver",
		Build:     "python manage.py collectstatic --noinput",
		Start:     "gunicorn config.wsgi:application",
		Port:      8000,
	}
}

// StarletteDetector detects Starlette projects.
type StarletteDetector struct{}

func (d *StarletteDetector) Language() string { return "python" }
func (d *StarletteDetector) Name() string     { return "starlette" }

func (d *StarletteDetector) Detect(projectPath string, info *ProjectInfo) bool {
	// Only detect if starlette is present but fastapi is not
	// (FastAPI includes Starlette)
	return HasDependency(info.Dependencies, "starlette") && !HasDependency(info.Dependencies, "fastapi")
}

func (d *StarletteDetector) Defaults() *FrameworkDefaults {
	return &FrameworkDefaults{
		Framework: "starlette",
		Language:  "python",
		Dev:       "uvicorn main:app --reload",
		Start:     "uvicorn main:app --host 0.0.0.0",
		Port:      8000,
	}
}

// RegisterPython registers Python language support with the registry.
func RegisterPython(r *Registry) {
	r.RegisterLanguage(&PythonDetector{}, &PythonInferrer{})

	// Framework detectors (order matters - more specific first)
	r.RegisterFramework(&FastAPIDetector{})
	r.RegisterFramework(&FlaskDetector{})
	r.RegisterFramework(&DjangoDetector{})
	r.RegisterFramework(&StarletteDetector{})
}
