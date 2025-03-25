package operations

import (
	"context"
	"fmt"

	"github.com/zinrai/deck8sops-dev/pkg/config"
	"github.com/zinrai/deck8sops-dev/pkg/executor"
	"github.com/zinrai/deck8sops-dev/pkg/utils"
)

func Create(ctx context.Context, cfg *config.Config, logger *utils.Logger) error {
	cmdExecutor := executor.NewCommandExecutor(logger)

	if err := cmdExecutor.EnsureRequiredCommands(); err != nil {
		return fmt.Errorf("command check failed: %w", err)
	}

	if err := cmdExecutor.CheckKubeConnection(ctx); err != nil {
		return fmt.Errorf("kubernetes connection check failed: %w", err)
	}

	helmExecutor := executor.NewHelmExecutor(cmdExecutor, logger)

	kubectlExecutor := executor.NewKubectlExecutor(cmdExecutor, logger)

	logger.Info("Starting to create %d operations", len(cfg.Operations))

	for i, operation := range cfg.Operations {
		logger.Info("[%d/%d] Processing operation: %s (type: %s)",
			i+1, len(cfg.Operations), operation.Name, operation.Type)

		var err error

		switch operation.Type {
		case "helm":
			err = helmExecutor.InstallChart(ctx, operation)
			if err == nil {
				err = helmExecutor.VerifyInstallation(ctx, operation)
			}
		case "kubectl":
			err = kubectlExecutor.ApplyManifest(ctx, operation)
			if err == nil {
				err = kubectlExecutor.VerifyInstallation(ctx, operation)
			}
		default:
			err = fmt.Errorf("unsupported operation type: %s", operation.Type)
		}

		if err != nil {
			return fmt.Errorf("failed to create operation %s: %w", operation.Name, err)
		}

		logger.Info("Successfully created operation: %s", operation.Name)
	}

	return nil
}
