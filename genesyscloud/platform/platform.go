package platform

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/process"
)

// The Platform package provides information as to which platform is executing the provider, namely, Terraform, OpenTofu, or a Debug Server.
// It also provides a mechanism to validate and execute a command against the appropriate binary for the platform

type Platform int

type platformConfig struct {
	platform     Platform
	binaryPath   string
	providerAddr string
}

var platformConfigSingleton *platformConfig

const (
	PlatformUnknown Platform = iota
	PlatformTerraform
	PlatformOpenTofu
	PlatformDebugServer
)

func (p Platform) String() string {
	switch p {
	case PlatformTerraform:
		return "terraform"
	case PlatformOpenTofu:
		return "tofu"
	case PlatformDebugServer:
		return "debug-server"
	default:
		return "unknown"
	}
}

func (p Platform) BinaryPath() string {
	return platformConfigSingleton.binaryPath
}

func (p Platform) Binary() string {
	if platformConfigSingleton.binaryPath == "" {
		return ""
	}
	pathSegments := strings.Split(platformConfigSingleton.binaryPath, string(os.PathSeparator))
	return pathSegments[len(pathSegments)-1]
}

func (p Platform) IsDebugServer() bool {
	return p == PlatformDebugServer
}

func (p Platform) GetProviderRegistry() string {
	switch p {
	case PlatformTerraform:
		return "registry.terraform.io"
	case PlatformOpenTofu:
		return "registry.opentofu.org"
	default:
		return ""
	}
}

func (p Platform) ExecuteCommand(ctx context.Context, args ...string) (commandOutput *CommandOutput, err error) {
	// Validate platform
	if p == PlatformDebugServer {
		return nil, fmt.Errorf("cannot execute platform command against debug server")
	}
	// Validate binary path
	if platformConfigSingleton.binaryPath == "" {
		return nil, fmt.Errorf("binary path is empty")
	}
	return executePlatformCommand(ctx, platformConfigSingleton.binaryPath, args)
}

func IsValidPlatform(p Platform) bool {
	switch p {
	case PlatformTerraform, PlatformOpenTofu, PlatformDebugServer:
		return true
	default:
		return false
	}
}

func (p Platform) Validate() error {
	if !IsValidPlatform(p) {
		return fmt.Errorf("Invalid platform value detected: %v. This is an error of the terraform-provider-genesyscloud provider. This may indicate the provider is running in an unsupported environment. Please ensure you're using a supported operating system and architecture.", p)
	}
	if platformConfigSingleton == nil {
		return fmt.Errorf("Platform configuration is not initialized. This is likely an internal provider error. Please file a bug report if this persists in the terraform-provider-genesyscloud issues list.")
	}
	if platformConfigSingleton.binaryPath == "" {
		return fmt.Errorf("Unable to determine provider binary path. This may indicate incorrect provider installation or an unsupported execution environment. Please verify your provider installation is complete.")
	}
	return nil
}

func GetPlatform() Platform {
	return platformConfigSingleton.platform
}

func init() {
	// Initialize the config once
	platformConfigSingleton = &platformConfig{}

	path, err := detectExecutingBinary()
	if err != nil {
		log.Printf(`Error detecting binary: %v`, err)
		platformConfigSingleton.platform = PlatformUnknown
		return
	}

	platformConfigSingleton.binaryPath = path
	defer detectedPlatformLog()

	// Verify binary exists and has proper permissions
	if err := verifyBinary(platformConfigSingleton.binaryPath); err != nil {
		log.Printf("binary verification failed: %v", err)
		return
	}

	debugPatterns := []string{
		"dlv",   // Delve debugger
		"debug", // Debug Server
	}

	for _, pattern := range debugPatterns {
		if strings.Contains(platformConfigSingleton.binaryPath, pattern) {
			platformConfigSingleton.platform = PlatformDebugServer
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	versionOutput, err := executePlatformCommand(ctx, platformConfigSingleton.binaryPath, []string{"version"})
	if err != nil {
		log.Printf("Failed to execute version command: %v", err)
		return
	}

	if strings.Contains(strings.ToLower(versionOutput.Stdout), "tofu") {
		platformConfigSingleton.platform = PlatformOpenTofu
	} else {
		platformConfigSingleton.platform = PlatformTerraform
	}

}

func detectedPlatformLog() {
	platform := GetPlatform()
	log.Printf("Detected executing platform is: %v", platform.String())
}

// detectExecutingBinary returns the path of the currently executing binary by finding
// the parent process and determining its executable path (either `terraform“ or `tofu`)
//
// Returns:
//   - string: The path to the executing binary
//   - error: An error if the process cannot be found or if the executable path cannot be determined
func detectExecutingBinary() (string, error) {
	ppid, err := os.FindProcess(os.Getppid())
	if err != nil {
		return "", err
	}
	tfProcess, err := process.NewProcess(int32(ppid.Pid))
	if err != nil {
		return "", err
	}

	exe, err := tfProcess.Exe()
	if err != nil {
		return "", err
	}

	return exe, nil
}

// verifyBinary performs basic security checks on the provided binary path to ensure
// it exists, is a regular file (not a symlink or directory), and has proper execute permissions.
//
// Parameters:
//   - path: The filesystem path to the binary to verify
//
// Returns:
//   - error: An error if any verification check fails, nil if all checks pass
func verifyBinary(path string) error {
	// Basic existence check
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat binary: %w", err)
	}

	// Ensure it's a regular file, not a symlink or directory
	if !info.Mode().IsRegular() {
		return fmt.Errorf("binary path is not a regular file")
	}

	// Check if we have execute permission
	if info.Mode().Perm()&0111 == 0 {
		return fmt.Errorf("binary is not executable")
	}

	return nil
}

// validateCommandArgs uses HashiCorp's flags parser to validate command arguments
// before they are passed to the platform binary (terraform/tofu).
//
// Parameters:
//   - args: Slice of string arguments to validate
//
// Returns:
//   - error: An error if any argument fails validation, nil if all arguments are valid
func validateCommandArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided")
	}

	// Additional custom validation if needed
	command := args[0]
	if !isAllowedCommand(command) {
		return fmt.Errorf("command %q is not allowed", command)
	}

	// TODO: If sub commands are intended to be used, consider
	// adding extra validation for these commands.
	return nil
}

// isAllowedCommand checks if the given command is in the allowed list
func isAllowedCommand(cmd string) bool {
	allowedCommands := map[string]bool{
		"init":         true,
		"plan":         true,
		"apply":        true,
		"destroy":      true,
		"validate":     true,
		"output":       true,
		"show":         true,
		"state":        true,
		"import":       true,
		"version":      true,
		"fmt":          true,
		"force-unlock": true,
		"providers":    true,
		"login":        true,
		"logout":       true,
		"refresh":      true,
		"graph":        true,
		"taint":        true,
		"untaint":      true,
		"workspace":    true,
		"metadata":     true,
		"test":         true,
		"console":      true,
	}

	return allowedCommands[strings.TrimPrefix(cmd, "-")]
}

type CommandOutput struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// ExecutePlatformCommand executes a command against the platform binary (`terraform` or `tofu`) with
// the provided arguments within the given context. It captures both stdout and stderr output from the
// command execution.
//
// Parameters:
//   - ctx: Context for command execution and timeout control
//   - args: Slice of string arguments to pass to the command
//
// Returns:
//   - stdoutString: The stdout output from the command execution
//   - stderrString: The stderr output from the command execution
//   - error: An error if the command fails, times out, or if the platform binary cannot be detected
//
// The function will return an error if it cannot detect the executing binary path
func executePlatformCommand(ctx context.Context, binaryPath string, args []string) (commandOutput *CommandOutput, err error) {
	var stdout, stderr bytes.Buffer

	// Validate context
	if ctx == nil {
		return nil, fmt.Errorf("nil context provided")
	}

	// Validate arguments
	if err := validateCommandArgs(args); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}

	// Verify binary exists and has proper permissions
	if err := verifyBinary(binaryPath); err != nil {
		return nil, fmt.Errorf("binary verification failed: %w", err)
	}

	cmd := exec.CommandContext(ctx, binaryPath)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Args = append(cmd.Args, args...)

	log.Printf("Running command against platform binary: %s", cmd.String())
	err = cmd.Run()
	output := &CommandOutput{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if cmd.ProcessState != nil {
		output.ExitCode = cmd.ProcessState.ExitCode()
	} else {
		output.ExitCode = -1
	}

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return output, ctx.Err()
		}
		return output, err
	}

	return output, nil

}