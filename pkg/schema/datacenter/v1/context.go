// Package v1 implements the v1 datacenter schema.
package v1

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// EvalContext holds the evaluation context for HCL expressions.
type EvalContext struct {
	// Variables holds datacenter-level variables
	Variables map[string]cty.Value

	// Environment holds environment-level context
	Environment *EnvironmentContext

	// Node holds current node context (during hook evaluation)
	Node *NodeContext

	// Modules holds outputs from evaluated modules
	Modules map[string]cty.Value
}

// EnvironmentContext holds environment-level values.
type EnvironmentContext struct {
	Name       string
	Datacenter string
	Account    string
	Region     string
}

// NodeContext holds the current resource node being processed.
type NodeContext struct {
	// Type is the resource type (database, deployment, service, etc.)
	Type string
	// Name is the resource name
	Name string
	// Component is the component this resource belongs to
	Component string
	// Inputs are the resource inputs from the component spec
	Inputs map[string]cty.Value
}

// NewEvalContext creates a new evaluation context.
func NewEvalContext() *EvalContext {
	return &EvalContext{
		Variables: make(map[string]cty.Value),
		Modules:   make(map[string]cty.Value),
	}
}

// WithVariable adds a variable to the context.
func (ctx *EvalContext) WithVariable(name string, value interface{}) *EvalContext {
	ctx.Variables[name] = toCtyValue(value)
	return ctx
}

// WithEnvironment sets the environment context.
func (ctx *EvalContext) WithEnvironment(env *EnvironmentContext) *EvalContext {
	ctx.Environment = env
	return ctx
}

// WithNode sets the current node context.
func (ctx *EvalContext) WithNode(node *NodeContext) *EvalContext {
	ctx.Node = node
	return ctx
}

// WithModule adds module outputs to the context.
func (ctx *EvalContext) WithModule(name string, outputs map[string]interface{}) *EvalContext {
	outputValues := make(map[string]cty.Value)
	for k, v := range outputs {
		outputValues[k] = toCtyValue(v)
	}
	ctx.Modules[name] = cty.ObjectVal(outputValues)
	return ctx
}

// ToHCLContext converts the EvalContext to an HCL evaluation context.
func (ctx *EvalContext) ToHCLContext() *hcl.EvalContext {
	vars := make(map[string]cty.Value)

	// Add variables
	if len(ctx.Variables) > 0 {
		vars["variable"] = cty.ObjectVal(ctx.Variables)
		vars["var"] = cty.ObjectVal(ctx.Variables) // Alias
	}

	// Add environment
	if ctx.Environment != nil {
		envValues := map[string]cty.Value{
			"name": cty.StringVal(ctx.Environment.Name),
		}
		if ctx.Environment.Datacenter != "" {
			envValues["datacenter"] = cty.StringVal(ctx.Environment.Datacenter)
		}
		if ctx.Environment.Account != "" {
			envValues["account"] = cty.StringVal(ctx.Environment.Account)
		}
		if ctx.Environment.Region != "" {
			envValues["region"] = cty.StringVal(ctx.Environment.Region)
		}
		vars["environment"] = cty.ObjectVal(envValues)
	}

	// Add node context
	if ctx.Node != nil {
		nodeValues := map[string]cty.Value{
			"type":      cty.StringVal(ctx.Node.Type),
			"name":      cty.StringVal(ctx.Node.Name),
			"component": cty.StringVal(ctx.Node.Component),
		}
		if len(ctx.Node.Inputs) > 0 {
			nodeValues["inputs"] = cty.ObjectVal(ctx.Node.Inputs)
		} else {
			nodeValues["inputs"] = cty.EmptyObjectVal
		}
		vars["node"] = cty.ObjectVal(nodeValues)
	}

	// Add modules
	if len(ctx.Modules) > 0 {
		vars["module"] = cty.ObjectVal(ctx.Modules)
	}

	return &hcl.EvalContext{
		Variables: vars,
		Functions: standardFunctions(),
	}
}

// Clone creates a copy of the evaluation context.
func (ctx *EvalContext) Clone() *EvalContext {
	newCtx := &EvalContext{
		Variables:   make(map[string]cty.Value),
		Modules:     make(map[string]cty.Value),
		Environment: ctx.Environment,
		Node:        ctx.Node,
	}

	for k, v := range ctx.Variables {
		newCtx.Variables[k] = v
	}

	for k, v := range ctx.Modules {
		newCtx.Modules[k] = v
	}

	return newCtx
}

// toCtyValue converts a Go value to a cty.Value.
func toCtyValue(v interface{}) cty.Value {
	if v == nil {
		return cty.NullVal(cty.DynamicPseudoType)
	}

	switch val := v.(type) {
	case string:
		return cty.StringVal(val)
	case int:
		return cty.NumberIntVal(int64(val))
	case int64:
		return cty.NumberIntVal(val)
	case float64:
		return cty.NumberFloatVal(val)
	case bool:
		return cty.BoolVal(val)
	case []string:
		if len(val) == 0 {
			return cty.ListValEmpty(cty.String)
		}
		vals := make([]cty.Value, len(val))
		for i, s := range val {
			vals[i] = cty.StringVal(s)
		}
		return cty.ListVal(vals)
	case []interface{}:
		if len(val) == 0 {
			return cty.ListValEmpty(cty.DynamicPseudoType)
		}
		vals := make([]cty.Value, len(val))
		for i, item := range val {
			vals[i] = toCtyValue(item)
		}
		return cty.TupleVal(vals)
	case map[string]interface{}:
		if len(val) == 0 {
			return cty.EmptyObjectVal
		}
		vals := make(map[string]cty.Value)
		for k, item := range val {
			vals[k] = toCtyValue(item)
		}
		return cty.ObjectVal(vals)
	case map[string]string:
		if len(val) == 0 {
			return cty.EmptyObjectVal
		}
		vals := make(map[string]cty.Value)
		for k, item := range val {
			vals[k] = cty.StringVal(item)
		}
		return cty.ObjectVal(vals)
	case cty.Value:
		return val
	default:
		// Fallback to string representation
		return cty.StringVal(fmt.Sprintf("%v", v))
	}
}

// fromCtyValue converts a cty.Value back to a Go value.
func fromCtyValue(v cty.Value) interface{} {
	if v.IsNull() {
		return nil
	}

	switch v.Type() {
	case cty.String:
		return v.AsString()
	case cty.Number:
		bf := v.AsBigFloat()
		if bf.IsInt() {
			i, _ := bf.Int64()
			return i
		}
		f, _ := bf.Float64()
		return f
	case cty.Bool:
		return v.True()
	default:
		if v.Type().IsListType() || v.Type().IsTupleType() {
			var result []interface{}
			for it := v.ElementIterator(); it.Next(); {
				_, elem := it.Element()
				result = append(result, fromCtyValue(elem))
			}
			return result
		}
		if v.Type().IsMapType() || v.Type().IsObjectType() {
			result := make(map[string]interface{})
			for it := v.ElementIterator(); it.Next(); {
				key, elem := it.Element()
				result[key.AsString()] = fromCtyValue(elem)
			}
			return result
		}
		return v.GoString()
	}
}
