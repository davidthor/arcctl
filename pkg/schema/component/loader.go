package component

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/architect-io/arcctl/pkg/errors"
	"github.com/architect-io/arcctl/pkg/schema/component/internal"
	"github.com/architect-io/arcctl/pkg/schema/component/v1"
	"gopkg.in/yaml.v3"
)

// Parser interface for version-specific parsers.
type Parser interface {
	ParseBytes(data []byte) (*v1.SchemaV1, error)
}

// Transformer interface for version-specific transformers.
type Transformer interface {
	Transform(schema *v1.SchemaV1) (*internal.InternalComponent, error)
}

// versionDetectingLoader implements the Loader interface with automatic version detection.
type versionDetectingLoader struct {
	parsers        map[string]*v1.Parser
	validators     map[string]*v1.Validator
	transformers   map[string]*v1.Transformer
	defaultVersion string
}

// NewLoader creates a new component loader that auto-detects schema version.
func NewLoader() Loader {
	return &versionDetectingLoader{
		parsers: map[string]*v1.Parser{
			"v1": v1.NewParser(),
		},
		validators: map[string]*v1.Validator{
			"v1": v1.NewValidator(),
		},
		transformers: map[string]*v1.Transformer{
			"v1": v1.NewTransformer(),
		},
		defaultVersion: "v1",
	}
}

// Load parses a component from the given path.
func (l *versionDetectingLoader) Load(path string) (Component, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeParse, fmt.Sprintf("failed to read %s", path), err)
	}

	comp, err := l.LoadFromBytes(data, path)
	if err != nil {
		return nil, err
	}

	// Set source path on internal component
	comp.Internal().SourcePath = path

	// Try to load README from the same directory
	dir := filepath.Dir(path)
	readme := l.loadReadme(dir)
	if readme != "" {
		comp.Internal().Readme = readme
	}

	return comp, nil
}

// loadReadme attempts to load a README file from the given directory.
// It checks for README.md, README.MD, readme.md, and README (in that order).
func (l *versionDetectingLoader) loadReadme(dir string) string {
	readmeNames := []string{"README.md", "README.MD", "readme.md", "Readme.md", "README"}
	for _, name := range readmeNames {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err == nil {
			return string(data)
		}
	}
	return ""
}

// LoadFromBytes parses a component from raw bytes.
func (l *versionDetectingLoader) LoadFromBytes(data []byte, sourcePath string) (Component, error) {
	// Detect version
	version, err := l.detectVersion(data)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeParse, "failed to detect schema version", err)
	}

	// Get parser for version
	parser, ok := l.parsers[version]
	if !ok {
		return nil, errors.New(errors.ErrCodeParse, fmt.Sprintf("unsupported schema version: %s", version))
	}

	// Parse schema
	schema, err := parser.ParseBytes(data)
	if err != nil {
		return nil, errors.ParseError(sourcePath, err)
	}

	// Validate schema
	validator, ok := l.validators[version]
	if !ok {
		return nil, errors.New(errors.ErrCodeParse, fmt.Sprintf("no validator for version: %s", version))
	}

	validationErrors := validator.Validate(schema)
	if len(validationErrors) > 0 {
		errMsgs := make([]string, len(validationErrors))
		for i, e := range validationErrors {
			errMsgs[i] = e.Error()
		}
		return nil, errors.ValidationError(
			"component validation failed",
			map[string]interface{}{
				"errors": errMsgs,
			},
		)
	}

	// Transform to internal representation
	transformer, ok := l.transformers[version]
	if !ok {
		return nil, errors.New(errors.ErrCodeParse, fmt.Sprintf("no transformer for version: %s", version))
	}

	internal, err := transformer.Transform(schema)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeParse, "failed to transform schema", err)
	}

	internal.SourcePath = sourcePath

	return newComponentWrapper(internal), nil
}

// Validate validates a component without fully parsing.
func (l *versionDetectingLoader) Validate(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(errors.ErrCodeParse, fmt.Sprintf("failed to read %s", path), err)
	}

	// Detect version
	version, err := l.detectVersion(data)
	if err != nil {
		return errors.Wrap(errors.ErrCodeParse, "failed to detect schema version", err)
	}

	// Get parser for version
	parser, ok := l.parsers[version]
	if !ok {
		return errors.New(errors.ErrCodeParse, fmt.Sprintf("unsupported schema version: %s", version))
	}

	// Parse schema
	schema, err := parser.ParseBytes(data)
	if err != nil {
		return errors.ParseError(path, err)
	}

	// Validate schema
	validator, ok := l.validators[version]
	if !ok {
		return errors.New(errors.ErrCodeParse, fmt.Sprintf("no validator for version: %s", version))
	}

	validationErrors := validator.Validate(schema)
	if len(validationErrors) > 0 {
		errMsgs := make([]string, len(validationErrors))
		for i, e := range validationErrors {
			errMsgs[i] = e.Error()
		}
		return errors.ValidationError(
			"component validation failed",
			map[string]interface{}{
				"errors": errMsgs,
			},
		)
	}

	return nil
}

// detectVersion detects the schema version from the data.
func (l *versionDetectingLoader) detectVersion(data []byte) (string, error) {
	// First, try to parse the version field
	var versionOnly struct {
		Version string `yaml:"version"`
	}
	if err := yaml.Unmarshal(data, &versionOnly); err == nil && versionOnly.Version != "" {
		// Normalize version format (strip leading 'v' if present for comparison)
		version := strings.ToLower(versionOnly.Version)
		if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}
		return version, nil
	}

	// Default to v1
	return l.defaultVersion, nil
}
