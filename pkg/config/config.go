package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// Func ReadAirgapManifest from a input file, and return a AirgapManifest struct or an error if something goes wrong
func ReadAirgapManifest(input string) (*AirgapManifest, error) {

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
	var manifest AirgapManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		log.Printf("failed to unmarshal YAML: %v", err)
		return nil, err
	}

	// Validate the manifest
	if err := validateAirgapManifest(&manifest); err != nil {
		return nil, err
	}

	return &manifest, nil

}

func validateAirgapManifest(manifest *AirgapManifest) error {
	// Validate the manifest
	if manifest.APIVersion == 0 {
		log.Printf("apiVersion is missing or invalid")
		return errors.New("apiVersion is missing or invalid")
	}
	if manifest.Components.Kubernetes.Rke2.Version == "" {
		log.Printf("kubernetes rke2 version is missing")
		return errors.New("kubernetes rke2 version is missing")
	}
	for i, helm := range manifest.Components.Helm {
		if helm.Name == "" || helm.Version == "" || helm.Location == "" || helm.Namespace == "" {
			log.Printf("helm %d has missing fields", i+1)
			return errors.New("helm has missing fields")
		}
		if helm.Location != "" && helm.Location[:3] != "oci" {
			log.Printf("helm %d location is not an OCI registry", i+1)
			return errors.New("helm location is not an OCI registry")
		}
	}
	for i, img := range manifest.Components.Images {
		if img.Name == "" || img.Version == "" || img.Location == "" {
			log.Printf("images %d has missing fields", i+1)
			return errors.New("images has missing fields")
		}
	}
	return nil
}
