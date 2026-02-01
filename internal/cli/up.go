package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/architect-io/arcctl/pkg/schema/component"
	"github.com/architect-io/arcctl/pkg/schema/component/inference"
	"github.com/architect-io/arcctl/pkg/state"
	"github.com/architect-io/arcctl/pkg/state/types"
	"github.com/spf13/cobra"
)

// resourceID generates a unique ID for a resource.
func resourceID(component, resourceType, name string) string {
	return fmt.Sprintf("%s/%s/%s", component, resourceType, name)
}

func newUpCmd() *cobra.Command {
	var (
		file      string
		name      string
		variables []string
		varFile   string
		detach    bool
		noOpen    bool
		port      int
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

			// Get absolute path for reliable operations
			absPath, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}

			// Determine architect.yml location
			componentFile := file
			if componentFile == "" {
				// Check if path is a file or directory
				info, err := os.Stat(absPath)
				if err != nil {
					return fmt.Errorf("failed to access path: %w", err)
				}
				if info.IsDir() {
					// Look for architect.yml in the directory
					componentFile = filepath.Join(absPath, "architect.yml")
					if _, err := os.Stat(componentFile); os.IsNotExist(err) {
						componentFile = filepath.Join(absPath, "architect.yaml")
					}
				} else {
					// Path is a file, use it directly
					componentFile = absPath
					absPath = filepath.Dir(absPath)
				}
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

			componentName := filepath.Base(absPath)
			fmt.Printf("Component: %s\n", componentName)
			fmt.Printf("Environment: %s (local)\n", envName)
			fmt.Println()

			// Create state manager (always local for 'up')
			mgr, err := upCreateStateManager()
			if err != nil {
				return fmt.Errorf("failed to create state manager: %w", err)
			}

			// Create cancellable context that responds to Ctrl+C
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Set up signal handling for graceful shutdown during provisioning
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigChan
				fmt.Println("\nInterrupted, cancelling...")
				cancel()
			}()

			// Ensure local datacenter exists
			dcName := "local"
			dc, err := mgr.GetDatacenter(ctx, dcName)
			if err != nil {
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
				env = &types.EnvironmentState{
					Name:       envName,
					Datacenter: dcName,
					Status:     types.EnvironmentStatusPending,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					Components: make(map[string]*types.ComponentState),
				}

				dc.Environments = append(dc.Environments, envName)
				dc.UpdatedAt = time.Now()
				if err := mgr.SaveDatacenter(ctx, dc); err != nil {
					return fmt.Errorf("failed to update datacenter: %w", err)
				}
			}

			// Initialize Docker provisioner
			basePort := port
			if basePort == 0 {
				basePort = 8080
			}
			provisioner := NewDockerProvisioner(envName, basePort)

			// Create progress table
			progress := NewProgressTable(os.Stdout)

			// Add resources to progress table with dependencies
			// Databases have no dependencies
			for _, db := range comp.Databases() {
				id := resourceID(componentName, "database", db.Name())
				progress.AddResource(id, db.Name(), "database", componentName, nil)
			}

			// Buckets have no dependencies
			for _, bucket := range comp.Buckets() {
				id := resourceID(componentName, "bucket", bucket.Name())
				progress.AddResource(id, bucket.Name(), "bucket", componentName, nil)
			}

			// Functions depend on databases
			var dbDeps []string
			for _, db := range comp.Databases() {
				dbDeps = append(dbDeps, resourceID(componentName, "database", db.Name()))
			}
			for _, fn := range comp.Functions() {
				id := resourceID(componentName, "function", fn.Name())
				progress.AddResource(id, fn.Name(), "function", componentName, dbDeps)
			}

			// Deployments depend on databases
			for _, depl := range comp.Deployments() {
				id := resourceID(componentName, "deployment", depl.Name())
				progress.AddResource(id, depl.Name(), "deployment", componentName, dbDeps)
			}

			// Services depend on deployments/functions
			var workloadDeps []string
			for _, fn := range comp.Functions() {
				workloadDeps = append(workloadDeps, resourceID(componentName, "function", fn.Name()))
			}
			for _, depl := range comp.Deployments() {
				workloadDeps = append(workloadDeps, resourceID(componentName, "deployment", depl.Name()))
			}
			for _, svc := range comp.Services() {
				id := resourceID(componentName, "service", svc.Name())
				progress.AddResource(id, svc.Name(), "service", componentName, workloadDeps)
			}

			// Routes depend on services/functions
			for _, route := range comp.Routes() {
				id := resourceID(componentName, "route", route.Name())
				// Routes can depend on services or directly on functions
				var routeDeps []string
				if route.Service() != "" {
					routeDeps = append(routeDeps, resourceID(componentName, "service", route.Service()))
				} else if route.Function() != "" {
					routeDeps = append(routeDeps, resourceID(componentName, "function", route.Function()))
				}
				progress.AddResource(id, route.Name(), "route", componentName, routeDeps)
			}

			// Print initial progress table
			progress.PrintInitial()

			// Ensure Docker network exists (not tracked as a resource)
			if err := provisioner.EnsureNetwork(ctx); err != nil {
				return fmt.Errorf("failed to create network: %w", err)
			}

			// Collect database connection info for deployments
			dbConnections := make(map[string]*DatabaseConnection)

			// Provision databases
			for _, db := range comp.Databases() {
				id := resourceID(componentName, "database", db.Name())
				progress.UpdateStatus(id, StatusInProgress, "")
				progress.PrintUpdate(id)

				conn, err := provisioner.ProvisionDatabase(ctx, db, componentName)
				if err != nil {
					progress.SetError(id, err)
					progress.PrintUpdate(id)
					return fmt.Errorf("failed to provision database %q: %w", db.Name(), err)
				}
				dbConnections[db.Name()] = conn
				progress.UpdateStatus(id, StatusCompleted, fmt.Sprintf("localhost:%d", conn.Port))
				progress.PrintUpdate(id)
			}

			// Provision buckets (placeholder - would need MinIO or similar)
			for _, bucket := range comp.Buckets() {
				id := resourceID(componentName, "bucket", bucket.Name())
				progress.UpdateStatus(id, StatusSkipped, "not yet implemented")
				progress.PrintUpdate(id)
			}

			// Build and deploy containers (deployments and functions)
			var appPort int

			// Helper function to build environment variables
			buildEnv := func() map[string]string {
				envVars := make(map[string]string)
				// Add database URLs
				for dbName, conn := range dbConnections {
					containerDBHost := fmt.Sprintf("%s-%s-%s", envName, componentName, dbName)
					containerURL := fmt.Sprintf("postgres://app:%s@%s:5432/%s?sslmode=disable",
						conn.Password, containerDBHost, conn.Database)
					envVars["DATABASE_URL"] = containerURL
					envVars[fmt.Sprintf("DB_%s_URL", strings.ToUpper(dbName))] = containerURL
					envVars[fmt.Sprintf("DB_%s_HOST", strings.ToUpper(dbName))] = containerDBHost
					envVars[fmt.Sprintf("DB_%s_PORT", strings.ToUpper(dbName))] = "5432"
					envVars[fmt.Sprintf("DB_%s_USER", strings.ToUpper(dbName))] = conn.Username
					envVars[fmt.Sprintf("DB_%s_PASSWORD", strings.ToUpper(dbName))] = conn.Password
					envVars[fmt.Sprintf("DB_%s_NAME", strings.ToUpper(dbName))] = conn.Database
				}
				// Add user-provided variables
				for k, v := range vars {
					envVars[k] = v
				}
				return envVars
			}

			// Helper function to build and run a workload (deployment or function)
			buildAndRun := func(name string, build component.Build, image string, resourceType string, resID string) error {
				progress.UpdateStatus(resID, StatusInProgress, "building")
				progress.PrintUpdate(resID)

				workloadEnv := buildEnv()

				var imageTag string
				if build != nil {
					buildCtx := build.Context()
					if !filepath.IsAbs(buildCtx) {
						buildCtx = filepath.Join(absPath, buildCtx)
					}

					dockerfile := ""
					explicitDockerfile := build.Dockerfile()
					if explicitDockerfile != "" && explicitDockerfile != "Dockerfile" {
						if !filepath.IsAbs(explicitDockerfile) {
							dockerfile = filepath.Join(buildCtx, explicitDockerfile)
						} else {
							dockerfile = explicitDockerfile
						}
					}

					// Collect build args (NEXT_PUBLIC_* vars need to be available at build time)
					buildArgs := make(map[string]string)
					for k, v := range vars {
						if strings.HasPrefix(k, "NEXT_PUBLIC_") {
							buildArgs[k] = v
						}
					}

					var err error
					imageTag, err = provisioner.BuildImage(ctx, name, buildCtx, dockerfile, buildArgs)
					if err != nil {
						progress.SetError(resID, err)
						progress.PrintUpdate(resID)
						return fmt.Errorf("failed to build %s %q: %w", resourceType, name, err)
					}
				} else if image != "" {
					imageTag = image
				} else {
					progress.UpdateStatus(resID, StatusSkipped, "no build context or image")
					progress.PrintUpdate(resID)
					return nil
				}

				containerPort := 3000 // Default for Next.js/Node apps
				ports := map[int]int{containerPort: 0}

				_, hostPort, err := provisioner.RunContainer(ctx, name, imageTag, componentName, workloadEnv, ports)
				if err != nil {
					progress.SetError(resID, err)
					progress.PrintUpdate(resID)
					return fmt.Errorf("failed to run %s %q: %w", resourceType, name, err)
				}
				appPort = hostPort
				progress.UpdateStatus(resID, StatusCompleted, fmt.Sprintf("port %d", hostPort))
				progress.PrintUpdate(resID)
				return nil
			}

			// Track running processes for cleanup
			var runningProcesses []*exec.Cmd

			// Process functions (preferred for Next.js apps)
			for _, fn := range comp.Functions() {
				fnID := resourceID(componentName, "function", fn.Name())

				if fn.IsSourceBased() {
					progress.UpdateStatus(fnID, StatusInProgress, "starting")
					progress.PrintUpdate(fnID)

					// Source-based functions run as local processes
					src := fn.Src()
					srcPath := src.Path()
					if !filepath.IsAbs(srcPath) {
						srcPath = filepath.Join(absPath, srcPath)
					}

					// Use inference to fill in missing values
					inferredInfo, err := inference.InferProjectWithOverrides(srcPath, src.Language(), src.Framework())
					if err != nil {
						// Warning only, continue
						_ = err
					}

					// Determine dev command (explicit > inferred)
					devCommand := inference.FirstNonEmpty(src.Dev(), inferredInfo.DevCommand)
					if devCommand == "" {
						progress.UpdateStatus(fnID, StatusSkipped, "no dev command")
						progress.PrintUpdate(fnID)
						continue
					}

					// Determine install command
					installCommand := inference.FirstNonEmpty(src.Install(), inferredInfo.InstallCommand)

					// Determine port
					fnPort := inference.FirstNonZero(fn.Port(), inferredInfo.Port, 3000)

					// Run install command if specified
					if installCommand != "" {
						installCmd := exec.CommandContext(ctx, "sh", "-c", installCommand)
						installCmd.Dir = srcPath
						installCmd.Stdout = os.Stdout
						installCmd.Stderr = os.Stderr
						if err := installCmd.Run(); err != nil {
							progress.SetError(fnID, fmt.Errorf("install failed: %w", err))
							progress.PrintUpdate(fnID)
							return fmt.Errorf("install command failed for %q: %w", fn.Name(), err)
						}
					}

					// Start the dev server
					devCmd := exec.CommandContext(ctx, "sh", "-c", devCommand)
					devCmd.Dir = srcPath

					// Merge environment variables
					devCmd.Env = os.Environ()
					for k, v := range buildEnv() {
						devCmd.Env = append(devCmd.Env, fmt.Sprintf("%s=%s", k, v))
					}
					for k, v := range fn.Environment() {
						devCmd.Env = append(devCmd.Env, fmt.Sprintf("%s=%s", k, v))
					}
					// Set PORT environment variable
					devCmd.Env = append(devCmd.Env, fmt.Sprintf("PORT=%d", fnPort))

					// Create pipes for output
					stdout, _ := devCmd.StdoutPipe()
					stderr, _ := devCmd.StderrPipe()

					if err := devCmd.Start(); err != nil {
						progress.SetError(fnID, err)
						progress.PrintUpdate(fnID)
						return fmt.Errorf("failed to start function %q: %w", fn.Name(), err)
					}

					runningProcesses = append(runningProcesses, devCmd)

					// Stream output with prefix
					go streamWithPrefix(stdout, fmt.Sprintf("[%s] ", fn.Name()))
					go streamWithPrefix(stderr, fmt.Sprintf("[%s] ", fn.Name()))

					appPort = fnPort
					progress.UpdateStatus(fnID, StatusCompleted, fmt.Sprintf("port %d", fnPort))
					progress.PrintUpdate(fnID)
					continue
				}
				// Container-based functions
				var build component.Build
				var image string
				if fn.Container() != nil {
					build = fn.Container().Build()
					image = fn.Container().Image()
				}
				if err := buildAndRun(fn.Name(), build, image, "function", fnID); err != nil {
					return err
				}
			}

			// Cleanup function for processes
			cleanupProcesses := func() {
				for _, cmd := range runningProcesses {
					if cmd.Process != nil {
						_ = cmd.Process.Signal(syscall.SIGTERM)
					}
				}
				// Give processes time to terminate gracefully
				time.Sleep(500 * time.Millisecond)
				for _, cmd := range runningProcesses {
					if cmd.Process != nil {
						_ = cmd.Process.Kill()
					}
				}
			}
			defer cleanupProcesses()

			// Process deployments
			for _, depl := range comp.Deployments() {
				deplID := resourceID(componentName, "deployment", depl.Name())
				if err := buildAndRun(depl.Name(), depl.Build(), depl.Image(), "deployment", deplID); err != nil {
					return err
				}
			}

			// Expose routes
			routeURLs := make(map[string]string)
			for _, route := range comp.Routes() {
				routeID := resourceID(componentName, "route", route.Name())
				progress.UpdateStatus(routeID, StatusInProgress, "")
				progress.PrintUpdate(routeID)

				url := fmt.Sprintf("http://localhost:%d", appPort)
				routeURLs[route.Name()] = url

				progress.UpdateStatus(routeID, StatusCompleted, url)
				progress.PrintUpdate(routeID)
			}

			// Print final progress summary
			progress.PrintFinalSummary()

			if len(routeURLs) > 0 {
				var primaryURL string
				for _, url := range routeURLs {
					primaryURL = url
					break
				}
				fmt.Printf("\nApplication running at %s\n", primaryURL)

				// Open browser unless --no-open is specified
				if !noOpen && !detach {
					openBrowser(primaryURL)
				}
			}

			// Update environment state
			env.Status = types.EnvironmentStatusReady
			env.UpdatedAt = time.Now()
			env.Variables = vars

			compState := &types.ComponentState{
				Name:       componentName,
				Version:    "local",
				Source:     path,
				Status:     types.ResourceStatusReady,
				Variables:  vars,
				DeployedAt: time.Now(),
				UpdatedAt:  time.Now(),
				Resources:  provisioner.ProvisionedResources(),
			}
			env.Components[componentName] = compState

			if err := mgr.SaveEnvironment(ctx, env); err != nil {
				return fmt.Errorf("failed to save environment state: %w", err)
			}

			if detach {
				fmt.Println()
				fmt.Println("Running in background. To stop:")
				fmt.Printf("  arcctl destroy environment %s\n", envName)
			} else {
				fmt.Println()
				fmt.Println("Press Ctrl+C to stop...")

				// Wait for context cancellation (already set up above)
				<-ctx.Done()

				fmt.Println()
				fmt.Println("Shutting down...")

				// Use a fresh context for cleanup since the original is cancelled
				cleanupCtx := context.Background()

				// Cleanup containers
				if err := CleanupByEnvName(cleanupCtx, envName); err != nil {
					fmt.Printf("Warning: failed to cleanup containers: %v\n", err)
				}

				// Update state
				env.Status = types.EnvironmentStatusPending
				env.UpdatedAt = time.Now()
				_ = mgr.SaveEnvironment(cleanupCtx, env)

				fmt.Println("Stopped.")
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

// openBrowser opens the default browser to the given URL.
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return
	}
	_ = cmd.Start()
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

// streamWithPrefix reads from a reader and prints each line with a prefix.
func streamWithPrefix(r io.Reader, prefix string) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			lines := strings.Split(string(buf[:n]), "\n")
			for _, line := range lines {
				if line != "" {
					fmt.Printf("%s%s\n", prefix, line)
				}
			}
		}
		if err != nil {
			break
		}
	}
}
