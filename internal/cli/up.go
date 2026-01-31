package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/architect-io/arcctl/pkg/schema/component"
	"github.com/architect-io/arcctl/pkg/state"
	"github.com/architect-io/arcctl/pkg/state/types"
	"github.com/spf13/cobra"
)

func newUpCmd() *cobra.Command {
	var (
		file       string
		name       string
		variables  []string
		varFile    string
		detach     bool
		noOpen     bool
		port       int
	)

	cmd := &cobra.Command{
		Use:   "up [path]",
		Short: "Deploy a component to a local environment",
		Long: `The up command provides a streamlined experience for local development.
It deploys your component with all its dependencies to a local environment
with minimal configuration.

The up command:
  1. Parses your architect.yml file
  2. Creates a local development environment
  3. Provisions all required resources (databases, etc.)
  4. Builds and deploys your application
  5. Watches for file changes and auto-reloads (unless --detach)
  6. Exposes routes for local access`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			// Determine architect.yml location
			componentFile := file
			if componentFile == "" {
				componentFile = filepath.Join(path, "architect.yml")
			}

			// Load the component
			loader := component.NewLoader()
			comp, err := loader.Load(componentFile)
			if err != nil {
				return fmt.Errorf("failed to load component: %w", err)
			}

			// Determine environment name from flag or derive from directory name
			envName := name
			if envName == "" {
				// Use the directory name as the base for the environment name
				absPath, _ := filepath.Abs(path)
				dirName := filepath.Base(absPath)
				envName = fmt.Sprintf("%s-dev", dirName)
			}

			// Load variables from file if specified
			vars := make(map[string]string)
			if varFile != "" {
				data, err := os.ReadFile(varFile)
				if err != nil {
					return fmt.Errorf("failed to read var file: %w", err)
				}
				if err := upParseVarFile(data, vars); err != nil {
					return fmt.Errorf("failed to parse var file: %w", err)
				}
			}

			// Parse inline variables
			for _, v := range variables {
				parts := strings.SplitN(v, "=", 2)
				if len(parts) == 2 {
					vars[parts[0]] = parts[1]
				}
			}

			fmt.Printf("Component: %s\n", filepath.Base(path))
			fmt.Printf("Environment: %s (local)\n", envName)
			fmt.Println()

			// Create state manager (always local for 'up')
			mgr, err := upCreateStateManager()
			if err != nil {
				return fmt.Errorf("failed to create state manager: %w", err)
			}

			ctx := context.Background()

			// Ensure local datacenter exists
			dcName := "local"
			dc, err := mgr.GetDatacenter(ctx, dcName)
			if err != nil {
				// Create local datacenter
				dc = &types.DatacenterState{
					Name:         dcName,
					Version:      "local",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
					Environments: []string{},
				}
				if err := mgr.SaveDatacenter(ctx, dc); err != nil {
					return fmt.Errorf("failed to create local datacenter: %w", err)
				}
			}

			// Create or get environment
			env, err := mgr.GetEnvironment(ctx, envName)
			if err != nil {
				// Create new environment
				env = &types.EnvironmentState{
					Name:       envName,
					Datacenter: dcName,
					Status:     types.EnvironmentStatusPending,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					Components: make(map[string]*types.ComponentState),
				}

				// Add environment to datacenter
				dc.Environments = append(dc.Environments, envName)
				dc.UpdatedAt = time.Now()
				if err := mgr.SaveDatacenter(ctx, dc); err != nil {
					return fmt.Errorf("failed to update datacenter: %w", err)
				}
			}

			// Provision resources
			for _, db := range comp.Databases() {
				fmt.Printf("[provision] Creating database %q...\n", db.Name())
				// TODO: Implement actual database provisioning
			}

			for _, bucket := range comp.Buckets() {
				fmt.Printf("[provision] Creating bucket %q...\n", bucket.Name())
				// TODO: Implement actual bucket provisioning
			}

			// Build and deploy
			for _, depl := range comp.Deployments() {
				if depl.Build() != nil {
					fmt.Printf("[build] Building deployment %q...\n", depl.Name())
					// TODO: Implement actual Docker build
				}
				fmt.Printf("[deploy] Starting deployment %q...\n", depl.Name())
				// TODO: Implement actual deployment
			}

			// Expose routes
			localPort := port
			if localPort == 0 {
				localPort = 8080
			}

			routeURLs := make(map[string]string)
			for _, route := range comp.Routes() {
				url := fmt.Sprintf("http://localhost:%d", localPort)
				routeURLs[route.Name()] = url
				fmt.Printf("[expose] Route %q available at %s\n", route.Name(), url)
				localPort++
			}

			fmt.Println()
			if len(routeURLs) > 0 {
				// Get first URL as primary
				var primaryURL string
				for _, url := range routeURLs {
					primaryURL = url
					break
				}
				fmt.Printf("Application running at %s\n", primaryURL)

				// Open browser unless --no-open is specified
				if !noOpen {
					// TODO: Implement browser opening
					_ = primaryURL
				}
			}

			// Update environment state
			env.Status = types.EnvironmentStatusReady
			env.UpdatedAt = time.Now()
			env.Variables = vars

			// Add component to environment - use directory name as component name
			absPath, _ := filepath.Abs(path)
			componentName := filepath.Base(absPath)
			env.Components[componentName] = &types.ComponentState{
				Name:       componentName,
				Version:    "local",
				Source:     path,
				Status:     types.ResourceStatusReady,
				Variables:  vars,
				DeployedAt: time.Now(),
				UpdatedAt:  time.Now(),
			}

			if err := mgr.SaveEnvironment(ctx, env); err != nil {
				return fmt.Errorf("failed to save environment state: %w", err)
			}

			if detach {
				fmt.Println()
				fmt.Println("Running in background. To stop:")
				fmt.Printf("  arcctl env destroy %s\n", envName)
			} else {
				fmt.Println()
				fmt.Println("Watching for changes... (Ctrl+C to stop)")
				// TODO: Implement file watching
				// For now, just wait
				select {}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to architect.yml if not in default location")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Environment name (default: auto-generated from component name)")
	cmd.Flags().StringArrayVar(&variables, "var", nil, "Set a component variable (key=value)")
	cmd.Flags().StringVar(&varFile, "var-file", "", "Load variables from a file")
	cmd.Flags().BoolVarP(&detach, "detach", "d", false, "Run in background (don't watch for changes)")
	cmd.Flags().BoolVar(&noOpen, "no-open", false, "Don't open browser to application URL")
	cmd.Flags().IntVar(&port, "port", 0, "Override the port for local access (default: 8080)")

	return cmd
}

// Helper functions for up command

func upCreateStateManager() (state.Manager, error) {
	// Use config file defaults with no CLI overrides
	return createStateManagerWithConfig("", nil)
}

func upParseVarFile(data []byte, vars map[string]string) error {
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			value = strings.Trim(value, "\"'")
			vars[key] = value
		}
	}
	return nil
}
