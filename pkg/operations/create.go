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

	for i, operator := range cfg.Operations {
		logger.Info("[%d/%d] Processing operator: %s (type: %s)",
			i+1, len(cfg.Operations), operator.Name, operator.Type)

		var err error

		switch operator.Type {
		case "helm":
			err = helmExecutor.InstallChart(ctx, operator)
			if err == nil {
				err = helmExecutor.VerifyInstallation(ctx, operator)
			}
		case "kubectl":
			err = kubectlExecutor.ApplyManifest(ctx, operator)
			if err == nil {
				err = kubectlExecutor.VerifyInstallation(ctx, operator)
			}
		default:
			err = fmt.Errorf("unsupported operator type: %s", operator.Type)
		}

		if err != nil {
			return fmt.Errorf("failed to create operator %s: %w", operator.Name, err)
		}

		logger.Info("Successfully created operator: %s", operator.Name)
	}

	return nil
}
