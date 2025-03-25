package executor

import (
	"context"
	"fmt"

	"github.com/zinrai/deck8sops-dev/pkg/config"
	"github.com/zinrai/deck8sops-dev/pkg/utils"
)

type KubectlExecutor struct {
	cmdExecutor *CommandExecutor
	logger      *utils.Logger
}

func NewKubectlExecutor(cmdExecutor *CommandExecutor, logger *utils.Logger) *KubectlExecutor {
	return &KubectlExecutor{
		cmdExecutor: cmdExecutor,
		logger:      logger,
	}
}

func (k *KubectlExecutor) ApplyManifest(ctx context.Context, operation config.Operator) error {
	if operation.KubectlConfig == nil {
		return fmt.Errorf("kubectl config is nil for operation %s", operation.Name)
	}

	k.logger.Info("Applying manifest for %s from %s",
		operation.Name, operation.KubectlConfig.ManifestFile)

	cmd := fmt.Sprintf("kubectl apply -f %s", operation.KubectlConfig.ManifestFile)
	_, err := k.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to apply manifest: %w", err)
	}

	k.logger.Info("Successfully applied manifest for %s", operation.Name)
	return nil
}

func (k *KubectlExecutor) DeleteManifest(ctx context.Context, operation config.Operator) error {
	if operation.KubectlConfig == nil {
		return fmt.Errorf("kubectl config is nil for operation %s", operation.Name)
	}

	k.logger.Info("Deleting manifest for %s from %s",
		operation.Name, operation.KubectlConfig.ManifestFile)

	cmd := fmt.Sprintf("kubectl delete -f %s", operation.KubectlConfig.ManifestFile)
	_, err := k.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to delete manifest: %w", err)
	}

	k.logger.Info("Successfully deleted manifest for %s", operation.Name)
	return nil
}

func (k *KubectlExecutor) VerifyInstallation(ctx context.Context, operation config.Operator) error {
	if operation.Namespace == "" {
		k.logger.Info("Namespace not specified for %s, skipping verification", operation.Name)
		return nil
	}

	k.logger.Info("Verifying installation for %s in namespace %s", operation.Name, operation.Namespace)

	cmd := fmt.Sprintf("kubectl get all -n %s", operation.Namespace)
	_, err := k.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("installation verification failed: %w", err)
	}

	return nil
}
