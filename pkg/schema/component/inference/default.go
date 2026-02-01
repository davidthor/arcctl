package inference

// DefaultRegistry returns a registry with all built-in language and framework support.
func DefaultRegistry() *Registry {
	r := NewRegistry()

	// Register all built-in languages
	RegisterJavaScript(r)
	RegisterGo(r)
	RegisterPython(r)

	return r
}

// InferProject is a convenience function that uses the default registry.
func InferProject(projectPath string) (*ProjectInfo, error) {
	return DefaultRegistry().Infer(projectPath, "", "")
}

// InferProjectWithOverrides is a convenience function that allows explicit overrides.
func InferProjectWithOverrides(projectPath, language, framework string) (*ProjectInfo, error) {
	return DefaultRegistry().Infer(projectPath, language, framework)
}
