package inference

import (
	"fmt"
	"sort"
)

// Registry manages language detectors, inferrers, and framework detectors.
// It provides the main entry point for inferring project configuration.
type Registry struct {
	detectors  []LanguageDetector
	inferrers  map[string]LanguageInferrer
	frameworks map[string][]FrameworkDetector
}

// NewRegistry creates an empty registry.
// Use RegisterLanguage and RegisterFramework to add support.
func NewRegistry() *Registry {
	return &Registry{
		detectors:  make([]LanguageDetector, 0),
		inferrers:  make(map[string]LanguageInferrer),
		frameworks: make(map[string][]FrameworkDetector),
	}
}

// RegisterLanguage adds support for a new language.
// The detector identifies the language, and the inferrer extracts project info.
func (r *Registry) RegisterLanguage(detector LanguageDetector, inferrer LanguageInferrer) {
	r.detectors = append(r.detectors, detector)
	r.inferrers[detector.Name()] = inferrer

	// Sort detectors by priority (highest first)
	sort.Slice(r.detectors, func(i, j int) bool {
		return r.detectors[i].Priority() > r.detectors[j].Priority()
	})
}

// RegisterFramework adds support for a new framework.
// Frameworks are detected within their language ecosystem.
func (r *Registry) RegisterFramework(detector FrameworkDetector) {
	lang := detector.Language()
	r.frameworks[lang] = append(r.frameworks[lang], detector)
}

// DetectLanguage identifies the programming language at the given path.
// Returns empty string if no language is detected.
func (r *Registry) DetectLanguage(projectPath string) string {
	for _, d := range r.detectors {
		if d.Detect(projectPath) {
			return d.Name()
		}
	}
	return ""
}

// Infer performs full inference for a project, combining language detection,
// project parsing, and framework detection.
//
// Parameters:
//   - projectPath: Path to the project directory
//   - explicitLang: Optional explicit language override (empty to auto-detect)
//   - explicitFramework: Optional explicit framework override (empty to auto-detect)
//
// Returns ProjectInfo with all detected/inferred values, or an error.
func (r *Registry) Infer(projectPath string, explicitLang string, explicitFramework string) (*ProjectInfo, error) {
	// 1. Detect or use explicit language
	language := explicitLang
	if language == "" {
		language = r.DetectLanguage(projectPath)
	}
	if language == "" {
		return nil, fmt.Errorf("could not detect language for project at %s", projectPath)
	}

	// 2. Run language-specific inference
	inferrer, ok := r.inferrers[language]
	if !ok {
		return nil, fmt.Errorf("no inferrer registered for language: %s", language)
	}

	info, err := inferrer.Infer(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to infer project info: %w", err)
	}
	info.Language = language

	// 3. Detect or use explicit framework
	framework := explicitFramework
	if framework == "" {
		for _, fd := range r.frameworks[language] {
			if fd.Detect(projectPath, info) {
				framework = fd.Name()
				break
			}
		}
	}
	info.Framework = framework

	// 4. Apply framework defaults for any missing values
	if info.Framework != "" {
		r.applyFrameworkDefaults(info)
	}

	return info, nil
}

// applyFrameworkDefaults fills in missing values from framework defaults.
func (r *Registry) applyFrameworkDefaults(info *ProjectInfo) {
	for _, fd := range r.frameworks[info.Language] {
		if fd.Name() == info.Framework {
			defaults := fd.Defaults()
			if defaults == nil {
				return
			}

			if info.DevCommand == "" && defaults.Dev != "" {
				info.DevCommand = defaults.Dev
			}
			if info.BuildCommand == "" && defaults.Build != "" {
				info.BuildCommand = defaults.Build
			}
			if info.StartCommand == "" && defaults.Start != "" {
				info.StartCommand = defaults.Start
			}
			if info.InstallCommand == "" && defaults.Install != "" {
				info.InstallCommand = defaults.Install
			}
			if info.Port == 0 && defaults.Port != 0 {
				info.Port = defaults.Port
			}
			if info.Handler == "" && defaults.Handler != "" {
				info.Handler = defaults.Handler
			}
			return
		}
	}
}

// GetFrameworkDefaults returns the defaults for a specific framework.
func (r *Registry) GetFrameworkDefaults(language, framework string) *FrameworkDefaults {
	for _, fd := range r.frameworks[language] {
		if fd.Name() == framework {
			return fd.Defaults()
		}
	}
	return nil
}

// SupportedLanguages returns a list of all registered language names.
func (r *Registry) SupportedLanguages() []string {
	languages := make([]string, 0, len(r.inferrers))
	for lang := range r.inferrers {
		languages = append(languages, lang)
	}
	sort.Strings(languages)
	return languages
}

// SupportedFrameworks returns a list of all registered framework names for a language.
func (r *Registry) SupportedFrameworks(language string) []string {
	frameworks := make([]string, 0)
	for _, fd := range r.frameworks[language] {
		frameworks = append(frameworks, fd.Name())
	}
	sort.Strings(frameworks)
	return frameworks
}
