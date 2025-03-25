package config

type Config struct {
	Operations []Operator `yaml:"operations"`
}

type Operator struct {
	Name          string         `yaml:"name"`
	Type          string         `yaml:"type"` // "helm" or "kubectl"
	Namespace     string         `yaml:"namespace"`
	HelmConfig    *HelmConfig    `yaml:"helmConfig,omitempty"`
	KubectlConfig *KubectlConfig `yaml:"kubectlConfig,omitempty"`
}

type HelmConfig struct {
	Repo       RepoInfo `yaml:"repo"`
	Chart      string   `yaml:"chart"`
	Version    string   `yaml:"version,omitempty"`
	ValuesFile string   `yaml:"valuesFile,omitempty"`
}

type RepoInfo struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type KubectlConfig struct {
	ManifestFile string `yaml:"manifestFile"`
}

func (o *Operator) Validate() error {
	if o.Name == "" {
		return ErrEmptyOperatorName
	}

	switch o.Type {
	case "helm":
		if o.HelmConfig == nil {
			return ErrMissingHelmConfig
		}
		if o.HelmConfig.Repo.Name == "" || o.HelmConfig.Repo.URL == "" {
			return ErrInvalidHelmRepo
		}
		if o.HelmConfig.Chart == "" {
			return ErrMissingHelmChart
		}
	case "kubectl":
		if o.KubectlConfig == nil {
			return ErrMissingKubectlConfig
		}
		if o.KubectlConfig.ManifestFile == "" {
			return ErrMissingManifestFile
		}
	default:
		return ErrUnsupportedOperatorType
	}

	return nil
}
