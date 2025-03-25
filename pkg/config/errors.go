package config

import "errors"

var (
	ErrEmptyOperatorName       = errors.New("operation name cannot be empty")
	ErrUnsupportedOperatorType = errors.New("unsupported operation type, must be 'helm', 'kubectl', or 'kustomize'")
	ErrMissingHelmConfig       = errors.New("helmConfig is required for helm type operations")
	ErrInvalidHelmRepo         = errors.New("helm repository name and URL must be specified")
	ErrMissingHelmChart        = errors.New("helm chart name must be specified")
	ErrMissingKubectlConfig    = errors.New("kubectlConfig is required for kubectl type operations")
	ErrMissingManifestFile     = errors.New("manifestFile is required for kubectl type operations")
	ErrMissingKustomizeConfig  = errors.New("kustomizeConfig is required for kustomize type operations")
	ErrMissingKustomizePath    = errors.New("path is required for kustomize type operations")
	ErrFileNotFound            = errors.New("config file not found")
	ErrInvalidYAML             = errors.New("invalid YAML format")
)
