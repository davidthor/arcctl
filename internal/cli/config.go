package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/architect-io/arcctl/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(newConfigCmd())
}

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage arcctl configuration",
		Long:  `Manage arcctl configuration settings stored in ~/.arcctl/config.yaml`,
	}

	cmd.AddCommand(newConfigGetCmd())
	cmd.AddCommand(newConfigSetCmd())
	cmd.AddCommand(newConfigListCmd())
	cmd.AddCommand(newConfigInitCmd())
	cmd.AddCommand(newConfigProfileCmd())

	return cmd
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			value, err := getConfigValue(cfg, args[0])
			if err != nil {
				return err
			}

			fmt.Println(value)
			return nil
		},
	}
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := config.NewManager()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if err := mgr.Set(args[0], args[1]); err != nil {
				return fmt.Errorf("failed to set config: %w", err)
			}

			fmt.Printf("Set %s = %s\n", args[0], args[1])
			return nil
		},
	}
}

func newConfigListCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			switch outputFormat {
			case "yaml":
				data, err := yaml.Marshal(cfg)
				if err != nil {
					return err
				}
				fmt.Print(string(data))
			case "json":
				// Use yaml with JSON compatible output
				data, err := yaml.Marshal(cfg)
				if err != nil {
					return err
				}
				fmt.Print(string(data))
			default:
				// Table format
				printConfigTable(cfg)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (table, yaml, json)")

	return cmd
}

func newConfigInitCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := config.DefaultConfigPath()
			if err != nil {
				return err
			}

			// Check if config already exists
			if _, err := os.Stat(configPath); err == nil && !force {
				return fmt.Errorf("config file already exists at %s (use --force to overwrite)", configPath)
			}

			cfg := config.DefaultConfig()
			if err := cfg.SaveToFile(configPath); err != nil {
				return fmt.Errorf("failed to create config file: %w", err)
			}

			fmt.Printf("Created config file at %s\n", configPath)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing config file")

	return cmd
}

func newConfigProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage configuration profiles",
	}

	cmd.AddCommand(newConfigProfileListCmd())
	cmd.AddCommand(newConfigProfileUseCmd())
	cmd.AddCommand(newConfigProfileCreateCmd())
	cmd.AddCommand(newConfigProfileDeleteCmd())

	return cmd
}

func newConfigProfileListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if len(cfg.Profiles) == 0 {
				fmt.Println("No profiles configured")
				return nil
			}

			fmt.Println("Available profiles:")
			for name := range cfg.Profiles {
				marker := "  "
				if name == cfg.ActiveProfile {
					marker = "* "
				}
				fmt.Printf("%s%s\n", marker, name)
			}

			return nil
		},
	}
}

func newConfigProfileUseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use <name>",
		Short: "Switch to a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := config.NewManager()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			cfg := mgr.Get()
			if _, ok := cfg.Profiles[args[0]]; !ok {
				return fmt.Errorf("profile %s not found", args[0])
			}

			if err := mgr.Set("active_profile", args[0]); err != nil {
				return fmt.Errorf("failed to set profile: %w", err)
			}

			fmt.Printf("Switched to profile: %s\n", args[0])
			return nil
		},
	}
}

func newConfigProfileCreateCmd() *cobra.Command {
	var datacenter, environment string
	var stateBackend string

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := config.NewManager()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			cfg := mgr.Get()
			if cfg.Profiles == nil {
				cfg.Profiles = make(map[string]config.ProfileConfig)
			}

			profile := config.ProfileConfig{
				Datacenter:  datacenter,
				Environment: environment,
			}

			if stateBackend != "" {
				profile.State = config.StateConfig{
					Backend: stateBackend,
				}
			}

			cfg.Profiles[args[0]] = profile

			if err := mgr.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Created profile: %s\n", args[0])
			return nil
		},
	}

	cmd.Flags().StringVarP(&datacenter, "datacenter", "d", "", "Default datacenter for the profile")
	cmd.Flags().StringVarP(&environment, "environment", "e", "", "Default environment for the profile")
	cmd.Flags().StringVar(&stateBackend, "state-backend", "", "State backend for the profile")

	return cmd
}

func newConfigProfileDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := config.NewManager()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			cfg := mgr.Get()
			if _, ok := cfg.Profiles[args[0]]; !ok {
				return fmt.Errorf("profile %s not found", args[0])
			}

			delete(cfg.Profiles, args[0])

			// Clear active profile if it was deleted
			if cfg.ActiveProfile == args[0] {
				cfg.ActiveProfile = ""
			}

			if err := mgr.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Deleted profile: %s\n", args[0])
			return nil
		},
	}
}

func getConfigValue(cfg *config.Config, key string) (string, error) {
	parts := strings.Split(key, ".")

	switch parts[0] {
	case "default_datacenter":
		return cfg.DefaultDatacenter, nil
	case "default_environment":
		return cfg.DefaultEnvironment, nil
	case "active_profile":
		return cfg.ActiveProfile, nil
	case "state":
		if len(parts) > 1 {
			switch parts[1] {
			case "backend":
				return cfg.State.Backend, nil
			default:
				if val, ok := cfg.State.Config[parts[1]]; ok {
					return val, nil
				}
			}
		}
	case "logging":
		if len(parts) > 1 {
			switch parts[1] {
			case "level":
				return cfg.Logging.Level, nil
			case "format":
				return cfg.Logging.Format, nil
			case "file":
				return cfg.Logging.File, nil
			}
		}
	case "plugins":
		if len(parts) > 1 && parts[1] == "default" {
			return cfg.Plugins.Default, nil
		}
	case "registry":
		if len(parts) > 1 && parts[1] == "default" {
			return cfg.Registry.Default, nil
		}
	case "secrets":
		if len(parts) > 1 && parts[1] == "provider" {
			return cfg.Secrets.Provider, nil
		}
	}

	return "", fmt.Errorf("unknown config key: %s", key)
}

func printConfigTable(cfg *config.Config) {
	fmt.Println("Configuration:")
	fmt.Println()

	if cfg.DefaultDatacenter != "" {
		fmt.Printf("  default_datacenter: %s\n", cfg.DefaultDatacenter)
	}
	if cfg.DefaultEnvironment != "" {
		fmt.Printf("  default_environment: %s\n", cfg.DefaultEnvironment)
	}
	if cfg.ActiveProfile != "" {
		fmt.Printf("  active_profile: %s\n", cfg.ActiveProfile)
	}

	fmt.Println()
	fmt.Println("State:")
	fmt.Printf("  backend: %s\n", cfg.State.Backend)
	for k, v := range cfg.State.Config {
		fmt.Printf("  %s: %s\n", k, v)
	}

	fmt.Println()
	fmt.Println("Secrets:")
	fmt.Printf("  provider: %s\n", cfg.Secrets.Provider)

	fmt.Println()
	fmt.Println("Logging:")
	fmt.Printf("  level: %s\n", cfg.Logging.Level)
	fmt.Printf("  format: %s\n", cfg.Logging.Format)
	if cfg.Logging.File != "" {
		fmt.Printf("  file: %s\n", cfg.Logging.File)
	}

	fmt.Println()
	fmt.Println("Plugins:")
	fmt.Printf("  default: %s\n", cfg.Plugins.Default)

	if len(cfg.Profiles) > 0 {
		fmt.Println()
		fmt.Println("Profiles:")
		for name, profile := range cfg.Profiles {
			marker := ""
			if name == cfg.ActiveProfile {
				marker = " (active)"
			}
			fmt.Printf("  %s%s:\n", name, marker)
			if profile.Datacenter != "" {
				fmt.Printf("    datacenter: %s\n", profile.Datacenter)
			}
			if profile.Environment != "" {
				fmt.Printf("    environment: %s\n", profile.Environment)
			}
		}
	}
}
