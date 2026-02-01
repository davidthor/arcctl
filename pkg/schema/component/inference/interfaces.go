// Package inference provides automatic detection of language, framework, and
// commands from project files for source-based functions.
package inference

// LanguageDetector detects the programming language of a project.
type LanguageDetector interface {
	// Name returns the language name (e.g., "javascript", "typescript", "python", "go")
	Name() string

	// Detect returns true if this language is detected at the given project path
	Detect(projectPath string) bool

	// Priority returns the detection priority (higher = checked first)
	// This is important when multiple languages might be present (e.g., TypeScript over JavaScript)
	Priority() int
}

// LanguageInferrer extracts project configuration for a specific language.
type LanguageInferrer interface {
	// Language returns which language this inferrer handles
	Language() string

	// Infer extracts project info from the given path
	Infer(projectPath string) (*ProjectInfo, error)
}

// FrameworkDetector detects frameworks within a language ecosystem.
type FrameworkDetector interface {
	// Language returns which language this detector applies to
	Language() string

	// Name returns the framework name (e.g., "nextjs", "fastapi", "gin")
	Name() string

	// Detect returns true if this framework is detected
	Detect(projectPath string, info *ProjectInfo) bool

	// Defaults returns default commands and configuration for this framework
	Defaults() *FrameworkDefaults
}

// ProjectInfo holds inferred project configuration.
// All fields are optional and may be empty if not detected.
type ProjectInfo struct {
	// Language and runtime
	Language string // e.g., "javascript", "typescript", "python", "go"
	Runtime  string // e.g., "nodejs20.x", "python3.11", "go1.21"

	// Framework
	Framework string // e.g., "nextjs", "fastapi", "gin"

	// Package manager
	PackageManager string // e.g., "npm", "pnpm", "yarn", "poetry", "pip"

	// Commands
	InstallCommand string // e.g., "npm install"
	DevCommand     string // e.g., "npm run dev"
	BuildCommand   string // e.g., "npm run build"
	StartCommand   string // e.g., "npm run start"

	// Entry points
	Entry   string // Main entry file
	Handler string // Lambda-style handler (e.g., "index.handler")

	// Network
	Port int // Default port the app listens on

	// Raw data (for framework detection)
	Dependencies    map[string]string // Package dependencies
	DevDependencies map[string]string // Dev dependencies
	Scripts         map[string]string // Package.json scripts or Makefile targets
}

// FrameworkDefaults provides fallback values for a framework.
type FrameworkDefaults struct {
	Framework string
	Language  string
	Dev       string
	Build     string
	Start     string
	Install   string
	Port      int
	Handler   string // Default handler pattern
}
