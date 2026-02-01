package inference

import (
	"os"
	"path/filepath"
)

// FileExists checks if a file exists at the given path.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists at the given path.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// FileExistsInProject checks if a file exists in the project directory.
func FileExistsInProject(projectPath, filename string) bool {
	return FileExists(filepath.Join(projectPath, filename))
}

// AnyFileExists checks if any of the given files exist in the project.
func AnyFileExists(projectPath string, filenames ...string) bool {
	for _, f := range filenames {
		if FileExistsInProject(projectPath, f) {
			return true
		}
	}
	return false
}

// FirstNonEmpty returns the first non-empty string from the arguments.
func FirstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// FirstNonZero returns the first non-zero int from the arguments.
func FirstNonZero(values ...int) int {
	for _, v := range values {
		if v != 0 {
			return v
		}
	}
	return 0
}

// HasDependency checks if a dependency exists in the dependencies map.
// Supports checking for exact match or prefix match (for scoped packages).
func HasDependency(deps map[string]string, name string) bool {
	_, ok := deps[name]
	return ok
}

// HasAnyDependency checks if any of the given dependencies exist.
func HasAnyDependency(deps map[string]string, names ...string) bool {
	for _, name := range names {
		if HasDependency(deps, name) {
			return true
		}
	}
	return false
}

// MergeDeps merges multiple dependency maps into one.
func MergeDeps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
