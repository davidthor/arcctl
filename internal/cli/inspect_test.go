package cli

import (
	"testing"

	"github.com/architect-io/arcctl/pkg/graph"
	"github.com/architect-io/arcctl/pkg/resolver"
	"github.com/stretchr/testify/assert"
)

func TestNewInspectCmd(t *testing.T) {
	cmd := newInspectCmd()

	assert.Equal(t, "inspect", cmd.Use)
	assert.Contains(t, cmd.Short, "Inspect")

	// Should have component subcommand
	subcommands := cmd.Commands()
	assert.Len(t, subcommands, 1)
	assert.Equal(t, "component [path|image]", subcommands[0].Use)
}

func TestInspectComponentCmd_Flags(t *testing.T) {
	cmd := newInspectComponentCmd()

	// Check expand flag
	expandFlag := cmd.Flags().Lookup("expand")
	assert.NotNil(t, expandFlag)
	assert.Equal(t, "false", expandFlag.DefValue)

	// Check file flag
	fileFlag := cmd.Flags().Lookup("file")
	assert.NotNil(t, fileFlag)
	assert.Equal(t, "f", fileFlag.Shorthand)
}

func TestExtractComponentName(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		resolved resolver.ResolvedComponent
		want     string
	}{
		{
			name: "OCI reference with tag",
			ref:  "ghcr.io/myorg/myapp:v1.0.0",
			resolved: resolver.ResolvedComponent{
				Type: resolver.ReferenceTypeOCI,
			},
			want: "myapp",
		},
		{
			name: "OCI reference without tag",
			ref:  "ghcr.io/myorg/myapp",
			resolved: resolver.ResolvedComponent{
				Type: resolver.ReferenceTypeOCI,
			},
			want: "myapp",
		},
		{
			name: "local directory path",
			ref:  "./my-component",
			resolved: resolver.ResolvedComponent{
				Type: resolver.ReferenceTypeLocal,
				Path: "/Users/test/my-component/architect.yml",
			},
			want: "my-component",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractComponentName(tt.ref, tt.resolved)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatNodeID(t *testing.T) {
	tests := []struct {
		name      string
		nodeType  graph.NodeType
		component string
		nodeName  string
		want      string
	}{
		{
			name:      "database node",
			nodeType:  graph.NodeTypeDatabase,
			component: "myapp",
			nodeName:  "main",
			want:      "[DB] myapp/main",
		},
		{
			name:      "deployment node",
			nodeType:  graph.NodeTypeDeployment,
			component: "myapp",
			nodeName:  "api",
			want:      "[DP] myapp/api",
		},
		{
			name:      "function node",
			nodeType:  graph.NodeTypeFunction,
			component: "myapp",
			nodeName:  "web",
			want:      "[FN] myapp/web",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := graph.NewNode(tt.nodeType, tt.component, tt.nodeName)
			got := formatNodeID(node)
			assert.Equal(t, tt.want, got)
		})
	}
}
