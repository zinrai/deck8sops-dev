package executor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zinrai/deck8sops-dev/pkg/config"
	"github.com/zinrai/deck8sops-dev/pkg/utils"
)

type KustomizeExecutor struct {
	cmdExecutor *CommandExecutor
	logger      *utils.Logger
}

func NewKustomizeExecutor(cmdExecutor *CommandExecutor, logger *utils.Logger) *KustomizeExecutor {
	return &KustomizeExecutor{
		cmdExecutor: cmdExecutor,
		logger:      logger,
	}
}

func (k *KustomizeExecutor) ApplyKustomize(ctx context.Context, operation config.Operator) error {
	if operation.KustomizeConfig == nil {
		return fmt.Errorf("kustomize config is nil for operation %s", operation.Name)
	}

	path := operation.KustomizeConfig.Path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("kustomize path does not exist: %s", path)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	k.logger.Info("Applying kustomize manifests for %s from %s",
		operation.Name, absPath)

	var cmd string
	cmd = fmt.Sprintf("kubectl apply -k %s", absPath)

	_, err = k.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to apply kustomize manifests: %w", err)
	}

	k.logger.Info("Successfully applied kustomize manifests for %s", operation.Name)
	return nil
}

func (k *KustomizeExecutor) DeleteKustomize(ctx context.Context, operation config.Operator) error {
	if operation.KustomizeConfig == nil {
		return fmt.Errorf("kustomize config is nil for operation %s", operation.Name)
	}

	path := operation.KustomizeConfig.Path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		k.logger.Info("Kustomize path does not exist: %s, skipping deletion", path)
		return nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	k.logger.Info("Deleting kustomize manifests for %s from %s",
		operation.Name, absPath)

	var cmd string
	cmd = fmt.Sprintf("kubectl delete -k %s", absPath)

	_, err = k.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to delete kustomize manifests: %w", err)
	}

	k.logger.Info("Successfully deleted kustomize manifests for %s", operation.Name)
	return nil
}

func (k *KustomizeExecutor) VerifyInstallation(ctx context.Context, operation config.Operator) error {
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
