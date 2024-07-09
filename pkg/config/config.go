package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// Func ReadReleaseManifest from a input file, and return a ReleaseManifest struct or an error if something goes wrong
func ReadReleaseManifest(input string) (*ReleaseManifest, error) {

	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Printf("file does not exist: %s", input)
		return nil, err
	}

	// Read file content
	data, err := os.ReadFile(input)
	if err != nil {
		log.Printf("failed to read file: %v", err)
		return nil, err
	}

	// Unmarshal YAML into struct
	var manifest ReleaseManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		log.Printf("failed to unmarshal YAML: %v", err)
		return nil, err
	}

	// Validate the manifest
	if err := validateReleaseManifest(&manifest); err != nil {
		return nil, err
	}

	return &manifest, nil

}

func validateReleaseManifest(manifest *ReleaseManifest) error {
	// Validate the manifest
	if manifest.APIVersion == 0 {
		log.Printf("apiVersion is missing or invalid")
		return errors.New("apiVersion is missing or invalid")
	}
	if manifest.ReleaseVersion == "" {
		log.Printf("releaseVersion is missing")
		return errors.New("releaseVersion is missing")
	}
	if len(manifest.SupoortedUpgrades) == 0 {
		log.Printf("supoortedUpgrades is missing or empty")
		return errors.New("supoortedUpgrades is missing or empty")
	}
	if manifest.Components.OperatingSystem.Upgrade.Version == "" {
		log.Printf("operatingSystem upgrade version is missing")
		return errors.New("operatingSystem upgrade version is missing")
	}
	if manifest.Components.Kubernetes.Rke2.Version == "" {
		log.Printf("kubernetes rke2 version is missing")
		return errors.New("kubernetes rke2 version is missing")
	}
	for i, img := range manifest.Components.Images {
		if img.Name == "" || img.Version == "" || img.Location == "" {
			log.Printf("image %d has missing fields", i+1)
			return errors.New("image has missing fields")
		}
	}
	return nil
}
