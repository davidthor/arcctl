package resolver

import (
	"testing"
)

func TestNewDependencyResolver(t *testing.T) {
	mockResolver := NewResolver(Options{
		AllowLocal: true,
	})

	depResolver := NewDependencyResolver(mockResolver)
	if depResolver == nil {
		t.Fatal("NewDependencyResolver returned nil")
	}

	if depResolver.resolver == nil {
		t.Error("resolver is nil")
	}
	if depResolver.loader == nil {
		t.Error("loader is nil")
	}
	if depResolver.resolved == nil {
		t.Error("resolved map is nil")
	}
	if depResolver.visiting == nil {
		t.Error("visiting map is nil")
	}
}

func TestDependencyGraph_GetDeploymentOrder(t *testing.T) {
	graph := &DependencyGraph{
		Order: []string{"dep1", "dep2", "root"},
	}

	order := graph.GetDeploymentOrder()
	if len(order) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(order))
	}

	expected := []string{"dep1", "dep2", "root"}
	for i, name := range order {
		if name != expected[i] {
			t.Errorf("Order[%d]: got %q, want %q", i, name, expected[i])
		}
	}
}

func TestDependencyGraph_GetDestroyOrder(t *testing.T) {
	graph := &DependencyGraph{
		Order: []string{"dep1", "dep2", "root"},
	}

	order := graph.GetDestroyOrder()
	if len(order) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(order))
	}

	// Destroy order should be reverse
	expected := []string{"root", "dep2", "dep1"}
	for i, name := range order {
		if name != expected[i] {
			t.Errorf("DestroyOrder[%d]: got %q, want %q", i, name, expected[i])
		}
	}
}

func TestDependencyGraph_GetDependency(t *testing.T) {
	graph := &DependencyGraph{
		All: map[string]ResolvedDependency{
			"api": {
				Name:  "api",
				Depth: 0,
			},
			"database": {
				Name:  "database",
				Depth: 1,
			},
		},
	}

	t.Run("existing dependency", func(t *testing.T) {
		dep, ok := graph.GetDependency("api")
		if !ok {
			t.Error("Expected to find 'api' dependency")
		}
		if dep.Name != "api" {
			t.Errorf("Name: got %q, want %q", dep.Name, "api")
		}
	})

	t.Run("non-existing dependency", func(t *testing.T) {
		_, ok := graph.GetDependency("nonexistent")
		if ok {
			t.Error("Should not find 'nonexistent' dependency")
		}
	})
}

func TestDependencyGraph_FlattenDependencies(t *testing.T) {
	graph := &DependencyGraph{
		Order: []string{"database", "api", "root"},
		All: map[string]ResolvedDependency{
			"database": {Name: "database", Depth: 2},
			"api":      {Name: "api", Depth: 1},
			"root":     {Name: "root", Depth: 0},
		},
	}

	deps := graph.FlattenDependencies()
	if len(deps) != 3 {
		t.Fatalf("Expected 3 dependencies, got %d", len(deps))
	}

	// Should be in topological order
	expectedNames := []string{"database", "api", "root"}
	for i, dep := range deps {
		if dep.Name != expectedNames[i] {
			t.Errorf("Dependency[%d]: got %q, want %q", i, dep.Name, expectedNames[i])
		}
	}
}

func TestDependencyGraph_HasCircularDependencies(t *testing.T) {
	t.Run("no circular dependencies", func(t *testing.T) {
		graph := &DependencyGraph{
			All: map[string]ResolvedDependency{
				"root": {
					Name: "root",
					Dependencies: []ResolvedDependency{
						{Name: "api"},
					},
				},
				"api": {
					Name: "api",
					Dependencies: []ResolvedDependency{
						{Name: "database"},
					},
				},
				"database": {
					Name:         "database",
					Dependencies: []ResolvedDependency{},
				},
			},
		}

		if graph.HasCircularDependencies() {
			t.Error("Should not detect circular dependencies")
		}
	})

	t.Run("with circular dependencies", func(t *testing.T) {
		graph := &DependencyGraph{
			All: map[string]ResolvedDependency{
				"a": {
					Name: "a",
					Dependencies: []ResolvedDependency{
						{Name: "b"},
					},
				},
				"b": {
					Name: "b",
					Dependencies: []ResolvedDependency{
						{Name: "c"},
					},
				},
				"c": {
					Name: "c",
					Dependencies: []ResolvedDependency{
						{Name: "a"}, // Circular back to 'a'
					},
				},
			},
		}

		if !graph.HasCircularDependencies() {
			t.Error("Should detect circular dependencies")
		}
	})

	t.Run("self-referencing dependency", func(t *testing.T) {
		graph := &DependencyGraph{
			All: map[string]ResolvedDependency{
				"self": {
					Name: "self",
					Dependencies: []ResolvedDependency{
						{Name: "self"}, // Self-reference
					},
				},
			},
		}

		if !graph.HasCircularDependencies() {
			t.Error("Should detect self-referencing circular dependency")
		}
	})

	t.Run("empty graph", func(t *testing.T) {
		graph := &DependencyGraph{
			All: map[string]ResolvedDependency{},
		}

		if graph.HasCircularDependencies() {
			t.Error("Empty graph should not have circular dependencies")
		}
	})
}

func TestResolvedDependency(t *testing.T) {
	dep := ResolvedDependency{
		Name: "my-component",
		Component: ResolvedComponent{
			Reference: "./my-component",
			Type:      ReferenceTypeLocal,
			Path:      "/path/to/component.yml",
		},
		Variables: map[string]string{
			"env": "production",
		},
		Depth: 1,
		Dependencies: []ResolvedDependency{
			{Name: "child-dep", Depth: 2},
		},
	}

	if dep.Name != "my-component" {
		t.Errorf("Name: got %q", dep.Name)
	}
	if dep.Component.Type != ReferenceTypeLocal {
		t.Errorf("Component.Type: got %q", dep.Component.Type)
	}
	if dep.Variables["env"] != "production" {
		t.Error("Variables not preserved")
	}
	if dep.Depth != 1 {
		t.Errorf("Depth: got %d", dep.Depth)
	}
	if len(dep.Dependencies) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(dep.Dependencies))
	}
}

func TestDependencyResolver_TopologicalSort(t *testing.T) {
	// Create a mock resolver for testing
	r := &DependencyResolver{
		resolved: map[string]ResolvedDependency{
			"root": {
				Name: "root",
				Dependencies: []ResolvedDependency{
					{Name: "api"},
					{Name: "web"},
				},
			},
			"api": {
				Name: "api",
				Dependencies: []ResolvedDependency{
					{Name: "database"},
				},
			},
			"web": {
				Name: "web",
				Dependencies: []ResolvedDependency{
					{Name: "api"},
				},
			},
			"database": {
				Name:         "database",
				Dependencies: []ResolvedDependency{},
			},
		},
		visiting: make(map[string]bool),
	}

	order := r.topologicalSort()

	// Verify that dependencies come before dependents
	positions := make(map[string]int)
	for i, name := range order {
		positions[name] = i
	}

	// database should come before api
	if positions["database"] > positions["api"] {
		t.Error("database should come before api in topological order")
	}

	// api should come before web
	if positions["api"] > positions["web"] {
		t.Error("api should come before web in topological order")
	}

	// api should come before root
	if positions["api"] > positions["root"] {
		t.Error("api should come before root in topological order")
	}

	// web should come before root
	if positions["web"] > positions["root"] {
		t.Error("web should come before root in topological order")
	}
}
