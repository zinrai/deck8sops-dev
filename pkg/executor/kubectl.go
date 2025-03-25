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

func (k *KubectlExecutor) ApplyManifest(ctx context.Context, operator config.Operator) error {
	if operator.KubectlConfig == nil {
		return fmt.Errorf("kubectl config is nil for operator %s", operator.Name)
	}

	k.logger.Info("Applying manifest for %s from %s",
		operator.Name, operator.KubectlConfig.ManifestFile)

	cmd := fmt.Sprintf("kubectl apply -f %s", operator.KubectlConfig.ManifestFile)
	_, err := k.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to apply manifest: %w", err)
	}

	k.logger.Info("Successfully applied manifest for %s", operator.Name)
	return nil
}

func (k *KubectlExecutor) DeleteManifest(ctx context.Context, operator config.Operator) error {
	if operator.KubectlConfig == nil {
		return fmt.Errorf("kubectl config is nil for operator %s", operator.Name)
	}

	k.logger.Info("Deleting manifest for %s from %s",
		operator.Name, operator.KubectlConfig.ManifestFile)

	cmd := fmt.Sprintf("kubectl delete -f %s", operator.KubectlConfig.ManifestFile)
	_, err := k.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to delete manifest: %w", err)
	}

	k.logger.Info("Successfully deleted manifest for %s", operator.Name)
	return nil
}

func (k *KubectlExecutor) VerifyInstallation(ctx context.Context, operator config.Operator) error {
	if operator.Namespace == "" {
		k.logger.Info("Namespace not specified for %s, skipping verification", operator.Name)
		return nil
	}

	k.logger.Info("Verifying installation for %s in namespace %s", operator.Name, operator.Namespace)

	cmd := fmt.Sprintf("kubectl get all -n %s", operator.Namespace)
	_, err := k.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("installation verification failed: %w", err)
	}

	return nil
}
