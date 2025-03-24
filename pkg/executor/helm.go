package executor

import (
	"context"
	"fmt"
	"strings"

	"github.com/zinrai/deck8sops-dev/pkg/config"
	"github.com/zinrai/deck8sops-dev/pkg/utils"
)

type HelmExecutor struct {
	cmdExecutor *CommandExecutor
	logger      *utils.Logger
}

func NewHelmExecutor(cmdExecutor *CommandExecutor, logger *utils.Logger) *HelmExecutor {
	return &HelmExecutor{
		cmdExecutor: cmdExecutor,
		logger:      logger,
	}
}

func (h *HelmExecutor) AddRepository(ctx context.Context, repo config.RepoInfo) error {
	output, err := h.cmdExecutor.Execute(ctx, fmt.Sprintf("helm repo list"))
	if err == nil && strings.Contains(output, repo.Name) {
		h.logger.Info("Helm repository %s already exists, updating", repo.Name)
	} else {
		_, err := h.cmdExecutor.Execute(ctx, fmt.Sprintf("helm repo add %s %s", repo.Name, repo.URL))
		if err != nil {
			return fmt.Errorf("failed to add helm repository: %w", err)
		}
		h.logger.Info("Added Helm repository %s", repo.Name)
	}

	_, err = h.cmdExecutor.Execute(ctx, "helm repo update")
	if err != nil {
		return fmt.Errorf("failed to update helm repositories: %w", err)
	}

	return nil
}

func (h *HelmExecutor) InstallChart(ctx context.Context, operator config.Operator) error {
	if operator.HelmConfig == nil {
		return fmt.Errorf("helm config is nil for operator %s", operator.Name)
	}

	err := h.AddRepository(ctx, operator.HelmConfig.Repo)
	if err != nil {
		return err
	}

	cmd := fmt.Sprintf("helm upgrade --install %s %s/%s --namespace %s --create-namespace",
		operator.Name,
		operator.HelmConfig.Repo.Name,
		operator.HelmConfig.Chart,
		operator.Namespace)

	if operator.HelmConfig.Version != "" {
		cmd += fmt.Sprintf(" --version %s", operator.HelmConfig.Version)
	}

	if operator.HelmConfig.ValuesFile != "" {
		if !h.cmdExecutor.FileExists(operator.HelmConfig.ValuesFile) {
			return fmt.Errorf("values file not found: %s", operator.HelmConfig.ValuesFile)
		}
		cmd += fmt.Sprintf(" --values %s", operator.HelmConfig.ValuesFile)
	}

	h.logger.Info("Installing Helm chart %s/%s as %s in namespace %s",
		operator.HelmConfig.Repo.Name,
		operator.HelmConfig.Chart,
		operator.Name,
		operator.Namespace)

	_, err = h.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to install helm chart: %w", err)
	}

	h.logger.Info("Successfully installed Helm chart %s", operator.Name)
	return nil
}

func (h *HelmExecutor) UninstallChart(ctx context.Context, operator config.Operator) error {
	h.logger.Info("Uninstalling Helm chart %s from namespace %s",
		operator.Name, operator.Namespace)

	cmd := fmt.Sprintf("helm uninstall %s --namespace %s", operator.Name, operator.Namespace)
	_, err := h.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to uninstall helm chart: %w", err)
	}

	h.logger.Info("Successfully uninstalled Helm chart %s", operator.Name)
	return nil
}

func (h *HelmExecutor) VerifyInstallation(ctx context.Context, operator config.Operator) error {
	h.logger.Info("Verifying installation for %s in namespace %s", operator.Name, operator.Namespace)

	releaseCmd := fmt.Sprintf("helm status %s --namespace %s", operator.Name, operator.Namespace)
	_, err := h.cmdExecutor.Execute(ctx, releaseCmd)
	if err != nil {
		return fmt.Errorf("helm release verification failed: %w", err)
	}

	podCmd := fmt.Sprintf("kubectl get all -n %s", operator.Namespace)
	_, err = h.cmdExecutor.Execute(ctx, podCmd)
	if err != nil {
		return fmt.Errorf("pod verification failed: %w", err)
	}

	return nil
}
