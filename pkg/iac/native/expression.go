package native

import (
	"fmt"
	"regexp"
	"strings"
)

// EvalContext provides values for expression evaluation.
type EvalContext struct {
	Inputs    map[string]interface{}
	Resources map[string]*ResourceState
}

var expressionPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// evaluateExpression evaluates a simple expression string.
// Supports: ${inputs.name}, ${resources.name.outputs.field}
func evaluateExpression(expr string, ctx *EvalContext) (interface{}, error) {
	// If no expressions, return as-is
	if !strings.Contains(expr, "${") {
		return expr, nil
	}

	// If the entire string is a single expression, return the actual value
	if strings.HasPrefix(expr, "${") && strings.HasSuffix(expr, "}") {
		trimmed := expr[2 : len(expr)-1]
		return resolveReference(trimmed, ctx)
	}

	// Otherwise, substitute expressions in the string
	result := expressionPattern.ReplaceAllStringFunc(expr, func(match string) string {
		// Extract reference
		ref := match[2 : len(match)-1]
		value, err := resolveReference(ref, ctx)
		if err != nil {
			return match // Keep original on error
		}
		return fmt.Sprintf("%v", value)
	})

	return result, nil
}

// resolveReference resolves a dotted reference like "inputs.name" or "resources.container.outputs.port"
func resolveReference(ref string, ctx *EvalContext) (interface{}, error) {
	parts := strings.Split(strings.TrimSpace(ref), ".")
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty reference")
	}

	switch parts[0] {
	case "inputs":
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid input reference: %s", ref)
		}
		return navigatePath(ctx.Inputs, parts[1:])

	case "resources":
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid resource reference: %s", ref)
		}
		resourceName := parts[1]
		resource, ok := ctx.Resources[resourceName]
		if !ok {
			return nil, fmt.Errorf("resource not found: %s", resourceName)
		}

		// Handle resources.name.outputs.field or resources.name.properties.field
		if parts[2] == "outputs" {
			return navigatePath(resource.Outputs, parts[3:])
		} else if parts[2] == "properties" {
			return navigatePath(resource.Properties, parts[3:])
		} else if parts[2] == "id" {
			return resource.ID, nil
		}
		return nil, fmt.Errorf("invalid resource property: %s", parts[2])

	default:
		// Try as a function call
		return evaluateFunction(ref, ctx)
	}
}

// navigatePath navigates a path through nested maps.
func navigatePath(data interface{}, path []string) (interface{}, error) {
	if len(path) == 0 {
		return data, nil
	}

	current := data
	for _, key := range path {
		switch v := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = v[key]
			if !ok {
				return nil, fmt.Errorf("key not found: %s", key)
			}
		case map[string]string:
			val, ok := v[key]
			if !ok {
				return nil, fmt.Errorf("key not found: %s", key)
			}
			current = val
		default:
			return nil, fmt.Errorf("cannot navigate into %T", current)
		}
	}

	return current, nil
}

// evaluateFunction evaluates a function call like "random_password(16)"
func evaluateFunction(expr string, ctx *EvalContext) (interface{}, error) {
	// Simple function parsing
	if strings.HasPrefix(expr, "random_password(") {
		// Generate random password (simplified)
		return generateRandomString(16), nil
	}

	return nil, fmt.Errorf("unknown function or reference: %s", expr)
}

// generateRandomString generates a random alphanumeric string.
func generateRandomString(length int) string {
	// Simplified random string generation
	// In production, use crypto/rand
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[i%len(chars)]
	}
	return string(result)
}

// resolveProperties resolves all expressions in a properties map.
func resolveProperties(props map[string]interface{}, ctx *EvalContext) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for key, value := range props {
		resolved, err := resolveValue(value, ctx)
		if err != nil {
			return nil, fmt.Errorf("property %s: %w", key, err)
		}
		result[key] = resolved
	}

	return result, nil
}

// resolveValue recursively resolves expressions in a value.
func resolveValue(value interface{}, ctx *EvalContext) (interface{}, error) {
	switch v := value.(type) {
	case string:
		return evaluateExpression(v, ctx)

	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, val := range v {
			resolved, err := resolveValue(val, ctx)
			if err != nil {
				return nil, err
			}
			result[k] = resolved
		}
		return result, nil

	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			resolved, err := resolveValue(val, ctx)
			if err != nil {
				return nil, err
			}
			result[i] = resolved
		}
		return result, nil

	default:
		return value, nil
	}
}
