package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/zinrai/deck8sops-dev/pkg/utils"
)

type CommandExecutor struct {
	logger *utils.Logger
}

func NewCommandExecutor(logger *utils.Logger) *CommandExecutor {
	return &CommandExecutor{
		logger: logger,
	}
}

func (e *CommandExecutor) Execute(ctx context.Context, command string) (string, error) {
	e.logger.Debug("Executing command: %s", command)

	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		e.logger.Debug("Command failed: %v", err)
		if stderrStr != "" {
			e.logger.CommandOutput(command, stderrStr)
		}
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	if stdoutStr != "" {
		e.logger.CommandOutput(command, stdoutStr)
	}

	return stdoutStr, nil
}

func (e *CommandExecutor) CheckCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func (e *CommandExecutor) EnsureRequiredCommands() error {
	requiredCommands := []string{"kubectl", "helm"}

	for _, cmd := range requiredCommands {
		if !e.CheckCommandExists(cmd) {
			return fmt.Errorf("required command not found: %s", cmd)
		}
	}

	return nil
}

func (e *CommandExecutor) CheckKubeConnection(ctx context.Context) error {
	_, err := e.Execute(ctx, "kubectl cluster-info")
	if err != nil {
		return fmt.Errorf("kubernetes cluster connection failed: %w", err)
	}
	return nil
}

func (e *CommandExecutor) FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
