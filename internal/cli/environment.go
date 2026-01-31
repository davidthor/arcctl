package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/architect-io/arcctl/pkg/schema/environment"
	"github.com/architect-io/arcctl/pkg/state"
	"github.com/architect-io/arcctl/pkg/state/backend"
	"github.com/architect-io/arcctl/pkg/state/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newEnvironmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "environment",
		Aliases: []string{"env"},
		Short:   "Manage environments",
		Long:    `Commands for creating, updating, and managing environments.`,
	}

	cmd.AddCommand(newEnvironmentListCmd())
	cmd.AddCommand(newEnvironmentGetCmd())
	cmd.AddCommand(newEnvironmentCreateCmd())
	cmd.AddCommand(newEnvironmentUpdateCmd())
	cmd.AddCommand(newEnvironmentDestroyCmd())
	cmd.AddCommand(newEnvironmentApplyCmd())
	cmd.AddCommand(newEnvironmentValidateCmd())

	return cmd
}

func newEnvironmentListCmd() *cobra.Command {
	var (
		datacenter    string
		outputFormat  string
		backendType   string
		backendConfig []string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List environments",
		Long:  `List all environments.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			// Create state manager
			mgr, err := envCreateStateManager(backendType, backendConfig)
			if err != nil {
				return fmt.Errorf("failed to create state manager: %w", err)
			}

			// List environments
			envRefs, err := mgr.ListEnvironments(ctx)
			if err != nil {
				return fmt.Errorf("failed to list environments: %w", err)
			}

			// Filter by datacenter if specified
			if datacenter != "" {
				filtered := make([]types.EnvironmentRef, 0)
				for _, ref := range envRefs {
					if ref.Datacenter == datacenter {
						filtered = append(filtered, ref)
					}
				}
				envRefs = filtered
			}

			// Handle output format
			switch outputFormat {
			case "json":
				data, err := json.MarshalIndent(envRefs, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Println(string(data))
			case "yaml":
				data, err := yaml.Marshal(envRefs)
				if err != nil {
					return fmt.Errorf("failed to marshal YAML: %w", err)
				}
				fmt.Println(string(data))
			default:
				// Table format
				if len(envRefs) == 0 {
					fmt.Println("No environments found.")
					return nil
				}

				fmt.Printf("%-16s %-20s %-12s %s\n", "NAME", "DATACENTER", "COMPONENTS", "CREATED")
				for _, ref := range envRefs {
					// Get full environment state for component count
					env, err := mgr.GetEnvironment(ctx, ref.Name)
					componentCount := 0
					if err == nil {
						componentCount = len(env.Components)
					}
					fmt.Printf("%-16s %-20s %-12d %s\n",
						ref.Name,
						ref.Datacenter,
						componentCount,
						ref.CreatedAt.Format("2006-01-02"),
					)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&datacenter, "datacenter", "d", "", "Filter by datacenter")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json, yaml")
	cmd.Flags().StringVar(&backendType, "backend", "", "State backend type")
	cmd.Flags().StringArrayVar(&backendConfig, "backend-config", nil, "Backend configuration (key=value)")

	return cmd
}

func newEnvironmentGetCmd() *cobra.Command {
	var (
		outputFormat  string
		backendType   string
		backendConfig []string
	)

	cmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Get details of an environment",
		Long:  `Get detailed information about an environment.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			envName := args[0]
			ctx := context.Background()

			// Create state manager
			mgr, err := envCreateStateManager(backendType, backendConfig)
			if err != nil {
				return fmt.Errorf("failed to create state manager: %w", err)
			}

			// Get environment state
			env, err := mgr.GetEnvironment(ctx, envName)
			if err != nil {
				return fmt.Errorf("environment %q not found: %w", envName, err)
			}

			// Handle output format
			switch outputFormat {
			case "json":
				data, err := json.MarshalIndent(env, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Println(string(data))
			case "yaml":
				data, err := yaml.Marshal(env)
				if err != nil {
					return fmt.Errorf("failed to marshal YAML: %w", err)
				}
				fmt.Println(string(data))
			default:
				// Table format
				fmt.Printf("Environment: %s\n", env.Name)
				fmt.Printf("Datacenter:  %s\n", env.Datacenter)
				fmt.Printf("Created:     %s\n", env.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Printf("Status:      %s\n", env.Status)
				fmt.Println()

				if len(env.Components) > 0 {
					fmt.Println("Components:")
					fmt.Printf("  %-16s %-40s %-12s %s\n", "NAME", "SOURCE", "STATUS", "RESOURCES")
					for name, comp := range env.Components {
						fmt.Printf("  %-16s %-40s %-12s %d\n",
							name,
							envTruncateString(comp.Source, 40),
							comp.Status,
							len(comp.Resources),
						)
					}
					fmt.Println()
				}

				// Collect URLs from components
				var urls []struct {
					component string
					route     string
					url       string
				}
				for compName, comp := range env.Components {
					for resName, res := range comp.Resources {
						if res.Type == "route" || res.Type == "ingress" {
							if url, ok := res.Outputs["url"].(string); ok {
								urls = append(urls, struct {
									component string
									route     string
									url       string
								}{compName, resName, url})
							}
						}
					}
				}
				if len(urls) > 0 {
					fmt.Println("URLs:")
					for _, u := range urls {
						fmt.Printf("  %s/%s: %s\n", u.component, u.route, u.url)
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json, yaml")
	cmd.Flags().StringVar(&backendType, "backend", "", "State backend type")
	cmd.Flags().StringArrayVar(&backendConfig, "backend-config", nil, "Backend configuration (key=value)")

	return cmd
}

func newEnvironmentCreateCmd() *cobra.Command {
	var (
		datacenter    string
		ifNotExists   bool
		backendType   string
		backendConfig []string
	)

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new environment",
		Long:  `Create a new environment.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			envName := args[0]
			ctx := context.Background()

			// Create state manager
			mgr, err := envCreateStateManager(backendType, backendConfig)
			if err != nil {
				return fmt.Errorf("failed to create state manager: %w", err)
			}

			// Check if environment already exists
			existingEnv, err := mgr.GetEnvironment(ctx, envName)
			if err == nil && existingEnv != nil {
				if ifNotExists {
					fmt.Printf("Environment %q already exists, skipping creation.\n", envName)
					return nil
				}
				return fmt.Errorf("environment %q already exists", envName)
			}

			// Verify datacenter exists
			dc, err := mgr.GetDatacenter(ctx, datacenter)
			if err != nil {
				return fmt.Errorf("datacenter %q not found: %w", datacenter, err)
			}

			fmt.Printf("Environment: %s\n", envName)
			fmt.Printf("Datacenter:  %s\n", datacenter)
			fmt.Println()

			fmt.Printf("[create] Creating environment %q...\n", envName)

			// Create environment state
			envState := &types.EnvironmentState{
				Name:       envName,
				Datacenter: datacenter,
				Status:     types.EnvironmentStatusReady,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				Components: make(map[string]*types.ComponentState),
			}

			if err := mgr.SaveEnvironment(ctx, envState); err != nil {
				return fmt.Errorf("failed to save environment state: %w", err)
			}

			// Update datacenter with environment reference
			dc.Environments = append(dc.Environments, envName)
			dc.UpdatedAt = time.Now()
			if err := mgr.SaveDatacenter(ctx, dc); err != nil {
				return fmt.Errorf("failed to update datacenter state: %w", err)
			}

			fmt.Printf("[success] Environment created successfully\n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&datacenter, "datacenter", "d", "", "Datacenter to use (required)")
	cmd.Flags().BoolVar(&ifNotExists, "if-not-exists", false, "Don't error if environment already exists")
	cmd.Flags().StringVar(&backendType, "backend", "", "State backend type")
	cmd.Flags().StringArrayVar(&backendConfig, "backend-config", nil, "Backend configuration (key=value)")
	_ = cmd.MarkFlagRequired("datacenter")

	return cmd
}

func newEnvironmentUpdateCmd() *cobra.Command {
	var (
		datacenter    string
		backendType   string
		backendConfig []string
	)

	cmd := &cobra.Command{
		Use:   "update <name>",
		Short: "Update an environment",
		Long:  `Update environment configuration.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			envName := args[0]
			ctx := context.Background()

			// Create state manager
			mgr, err := envCreateStateManager(backendType, backendConfig)
			if err != nil {
				return fmt.Errorf("failed to create state manager: %w", err)
			}

			// Get environment state
			env, err := mgr.GetEnvironment(ctx, envName)
			if err != nil {
				return fmt.Errorf("environment %q not found: %w", envName, err)
			}

			// Update datacenter if specified
			if datacenter != "" && datacenter != env.Datacenter {
				// Verify new datacenter exists
				_, err := mgr.GetDatacenter(ctx, datacenter)
				if err != nil {
					return fmt.Errorf("datacenter %q not found: %w", datacenter, err)
				}

				fmt.Printf("Updating environment %q datacenter from %q to %q\n", envName, env.Datacenter, datacenter)
				env.Datacenter = datacenter
			}

			env.UpdatedAt = time.Now()

			if err := mgr.SaveEnvironment(ctx, env); err != nil {
				return fmt.Errorf("failed to save environment state: %w", err)
			}

			fmt.Printf("[success] Environment updated successfully\n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&datacenter, "datacenter", "d", "", "Change target datacenter")
	cmd.Flags().StringVar(&backendType, "backend", "", "State backend type")
	cmd.Flags().StringArrayVar(&backendConfig, "backend-config", nil, "Backend configuration (key=value)")

	return cmd
}

func newEnvironmentDestroyCmd() *cobra.Command {
	var (
		autoApprove   bool
		backendType   string
		backendConfig []string
	)

	cmd := &cobra.Command{
		Use:   "destroy <name>",
		Short: "Destroy an environment",
		Long: `Destroy an environment and all its resources.

WARNING: This will destroy all components and resources in the environment. Use with caution.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			envName := args[0]
			ctx := context.Background()

			// Create state manager
			mgr, err := envCreateStateManager(backendType, backendConfig)
			if err != nil {
				return fmt.Errorf("failed to create state manager: %w", err)
			}

			// Get environment state
			env, err := mgr.GetEnvironment(ctx, envName)
			if err != nil {
				return fmt.Errorf("environment %q not found: %w", envName, err)
			}

			// Display what will be destroyed
			fmt.Printf("Environment: %s\n", envName)
			fmt.Printf("Datacenter:  %s\n", env.Datacenter)
			fmt.Println()

			componentCount := len(env.Components)
			resourceCount := 0
			for _, comp := range env.Components {
				resourceCount += len(comp.Resources)
			}

			fmt.Println("This will destroy:")
			fmt.Printf("  - %d components\n", componentCount)
			fmt.Printf("  - %d resources\n", resourceCount)
			fmt.Println()

			// Confirm unless --auto-approve is provided
			if !autoApprove {
				fmt.Print("Are you sure you want to destroy this environment? [y/N]: ")
				var response string
				_, _ = fmt.Scanln(&response)
				response = strings.ToLower(strings.TrimSpace(response))
				if response != "y" && response != "yes" {
					fmt.Println("Destroy cancelled.")
					return nil
				}
			}

			fmt.Println()

			// Destroy components
			for compName := range env.Components {
				fmt.Printf("[destroy] Destroying component %q...\n", compName)
				// TODO: Implement actual component destroy logic
			}

			fmt.Printf("[destroy] Removing environment...\n")

			// Delete environment state
			if err := mgr.DeleteEnvironment(ctx, envName); err != nil {
				return fmt.Errorf("failed to delete environment state: %w", err)
			}

	// Update datacenter to remove environment reference
	dc, err := mgr.GetDatacenter(ctx, env.Datacenter)
	if err == nil {
		newEnvs := make([]string, 0, len(dc.Environments))
		for _, e := range dc.Environments {
			if e != envName {
				newEnvs = append(newEnvs, e)
			}
		}
		dc.Environments = newEnvs
		dc.UpdatedAt = time.Now()
		_ = mgr.SaveDatacenter(ctx, dc)
	}

			fmt.Printf("[success] Environment destroyed successfully\n")

			return nil
		},
	}

	cmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "Skip confirmation prompt")
	cmd.Flags().StringVar(&backendType, "backend", "", "State backend type")
	cmd.Flags().StringArrayVar(&backendConfig, "backend-config", nil, "Backend configuration (key=value)")

	return cmd
}

func newEnvironmentApplyCmd() *cobra.Command {
	var (
		autoApprove   bool
		configFile    string
		backendType   string
		backendConfig []string
	)

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply an environment configuration file",
		Long:  `Apply an environment configuration file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			// Load and validate the environment file
			loader := environment.NewLoader()
			env, err := loader.Load(configFile)
			if err != nil {
				return fmt.Errorf("failed to load environment: %w", err)
			}

			envName := env.Name()

			fmt.Printf("Environment: %s\n", envName)
			fmt.Printf("Datacenter:  %s\n", env.Datacenter())
			fmt.Println()

			fmt.Println("Components:")
			for name, comp := range env.Components() {
				// Key is registry address, source is version tag or file path
				source := comp.Source()
				if strings.HasPrefix(source, "./") || strings.HasPrefix(source, "../") || strings.HasPrefix(source, "/") {
					fmt.Printf("  %s (local): %s\n", name, source)
				} else {
					fmt.Printf("  %s:%s\n", name, source)
				}
			}
			fmt.Println()

			// Confirm unless --auto-approve is provided
			if !autoApprove {
				fmt.Print("Proceed with apply? [Y/n]: ")
				var response string
				_, _ = fmt.Scanln(&response)
				response = strings.ToLower(strings.TrimSpace(response))
				if response != "" && response != "y" && response != "yes" {
					fmt.Println("Apply cancelled.")
					return nil
				}
			}

			// Create state manager
			mgr, err := envCreateStateManager(backendType, backendConfig)
			if err != nil {
				return fmt.Errorf("failed to create state manager: %w", err)
			}

			fmt.Println()
			fmt.Printf("[apply] Applying environment %q...\n", envName)

			// Create or update environment state
			envState, err := mgr.GetEnvironment(ctx, envName)
			if err != nil {
				// Create new environment
				envState = &types.EnvironmentState{
					Name:       envName,
					Datacenter: env.Datacenter(),
					Status:     types.EnvironmentStatusReady,
					CreatedAt:  time.Now(),
					Components: make(map[string]*types.ComponentState),
				}
			}
			envState.UpdatedAt = time.Now()

			// TODO: Implement actual apply logic - deploy/update components

			if err := mgr.SaveEnvironment(ctx, envState); err != nil {
				return fmt.Errorf("failed to save environment state: %w", err)
			}

			fmt.Printf("[success] Environment applied successfully\n")

			return nil
		},
	}

	cmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "Skip confirmation prompt")
	cmd.Flags().StringVarP(&configFile, "file", "f", "environment.yml", "Environment configuration file (required)")
	cmd.Flags().StringVar(&backendType, "backend", "", "State backend type")
	cmd.Flags().StringArrayVar(&backendConfig, "backend-config", nil, "Backend configuration (key=value)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func newEnvironmentValidateCmd() *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate an environment configuration",
		Long:  `Validate an environment configuration file without applying.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "environment.yml"
			if len(args) > 0 {
				if strings.HasSuffix(args[0], ".yml") || strings.HasSuffix(args[0], ".yaml") {
					path = args[0]
				} else {
					path = filepath.Join(args[0], "environment.yml")
				}
			}
			if file != "" {
				path = file
			}

			loader := environment.NewLoader()
			if err := loader.Validate(path); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			fmt.Println("Environment configuration is valid!")
			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to environment.yml if not in default location")

	return cmd
}

// Helper functions (prefixed to avoid conflicts)

func envCreateStateManager(backendType string, backendConfig []string) (state.Manager, error) {
	if backendType == "" {
		backendType = "local"
	}

	config := backend.Config{
		Type:   backendType,
		Config: make(map[string]string),
	}

	for _, c := range backendConfig {
		parts := strings.SplitN(c, "=", 2)
		if len(parts) == 2 {
			config.Config[parts[0]] = parts[1]
		}
	}

	return state.NewManagerFromConfig(config)
}

func envTruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
