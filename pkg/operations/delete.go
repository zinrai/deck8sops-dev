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
		operator := operations[i]
		logger.Info("[%d/%d] Deleting operator: %s (type: %s)",
			len(operations)-i, len(operations), operator.Name, operator.Type)

		var err error

		switch operator.Type {
		case "helm":
			err = helmExecutor.UninstallChart(ctx, operator)
		case "kubectl":
			err = kubectlExecutor.DeleteManifest(ctx, operator)
		default:
			err = fmt.Errorf("unsupported operator type: %s", operator.Type)
		}

		if err != nil {
			logger.Error("Failed to delete operator %s: %v", operator.Name, err)
			continue
		}

		logger.Info("Successfully deleted operator: %s", operator.Name)
	}

	return nil
}
