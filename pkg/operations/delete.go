package operations

import (
	"context"
	"fmt"

	"github.com/zinrai/deck8sops-dev/pkg/config"
	"github.com/zinrai/deck8sops-dev/pkg/executor"
	"github.com/zinrai/deck8sops-dev/pkg/utils"
)

func Delete(ctx context.Context, cfg *config.Config, logger *utils.Logger) error {
	cmdExecutor := executor.NewCommandExecutor(logger)

	if err := cmdExecutor.EnsureRequiredCommands(); err != nil {
		return fmt.Errorf("command check failed: %w", err)
	}

	if err := cmdExecutor.CheckKubeConnection(ctx); err != nil {
		return fmt.Errorf("kubernetes connection check failed: %w", err)
	}

	helmExecutor := executor.NewHelmExecutor(cmdExecutor, logger)

	kubectlExecutor := executor.NewKubectlExecutor(cmdExecutor, logger)

	operations := cfg.Operations
	logger.Info("Starting to delete %d operations in reverse order", len(operations))

	for i := len(operations) - 1; i >= 0; i-- {
		operation := operations[i]
		logger.Info("[%d/%d] Deleting operation: %s (type: %s)",
			len(operations)-i, len(operations), operation.Name, operation.Type)

		var err error

		switch operation.Type {
		case "helm":
			err = helmExecutor.UninstallChart(ctx, operation)
		case "kubectl":
			err = kubectlExecutor.DeleteManifest(ctx, operation)
		default:
			err = fmt.Errorf("unsupported operation type: %s", operation.Type)
		}

		if err != nil {
			logger.Error("Failed to delete operation %s: %v", operation.Name, err)
			continue
		}

		logger.Info("Successfully deleted operation: %s", operation.Name)
	}

	return nil
}
