package config

// AirgapManifest is the struct that represents the airgap manifest
type AirgapManifest struct {
	APIVersion float64 `yaml:"apiVersion"`
	Components struct {
		Kubernetes struct {
			Rke2 struct {
				Version string `yaml:"version"`
			} `yaml:"rke2"`
		} `yaml:"kubernetes"`
		Helm []struct {
			Name      string `yaml:"name"`
			Version   string `yaml:"version"`
			Location  string `yaml:"location"`
			Namespace string `yaml:"namespace"`
		} `yaml:"helm"`
		Images []struct {
			Name     string `yaml:"name"`
			Version  string `yaml:"version"`
			Location string `yaml:"location"`
		} `yaml:"images"`
	} `yaml:"components"`
}
