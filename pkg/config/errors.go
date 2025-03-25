package config

import "errors"

var (
	ErrEmptyOperatorName       = errors.New("operator name cannot be empty")
	ErrUnsupportedOperatorType = errors.New("unsupported operator type, must be 'helm' or 'kubectl'")
	ErrMissingHelmConfig       = errors.New("helmConfig is required for helm type operators")
	ErrInvalidHelmRepo         = errors.New("helm repository name and URL must be specified")
	ErrMissingHelmChart        = errors.New("helm chart name must be specified")
	ErrMissingKubectlConfig    = errors.New("kubectlConfig is required for kubectl type operators")
	ErrMissingManifestFile     = errors.New("manifestFile is required for kubectl type operators")
	ErrFileNotFound            = errors.New("config file not found")
	ErrInvalidYAML             = errors.New("invalid YAML format")
)
