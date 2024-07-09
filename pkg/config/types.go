package config

// ReleaseManifest is the struct that represents the release manifest
type ReleaseManifest struct {
	APIVersion        float64  `yaml:"apiVersion"`
	ReleaseVersion    string   `yaml:"releaseVersion"`
	SupoortedUpgrades []string `yaml:"supoortedUpgrades"`
	Components        struct {
		OperatingSystem struct {
			Upgrade struct {
				Version string `yaml:"version"`
			} `yaml:"upgrade"`
		} `yaml:"operatingSystem"`
		Kubernetes struct {
			Rke2 struct {
				Version string `yaml:"version"`
			} `yaml:"rke2"`
			K3S struct {
				Version string `yaml:"version"`
			} `yaml:"k3s"`
		} `yaml:"kubernetes"`
		Helm []struct {
			Longhorn struct {
				Version   string `yaml:"version"`
				Location  string `yaml:"location"`
				Namespace string `yaml:"namespace"`
			} `yaml:"longhorn"`
		} `yaml:"helm"`
		Images []struct {
			Name     string `yaml:"name"`
			Version  string `yaml:"version"`
			Location string `yaml:"location"`
		} `yaml:"images"`
	} `yaml:"components"`
}
