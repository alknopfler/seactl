package config

import (
	"os"
	"testing"
)

func TestReadReleaseManifest(t *testing.T) {
	validYAML := `
apiVersion: 1.0
releaseVersion: "v1.0.0"
supoortedUpgrades:
  - "v1.1.0"
components:
  operatingSystem:
    upgrade:
      version: "v1.0.0"
  kubernetes:
    rke2:
      version: "v1.20.0"
  longhorn:
    version: "v1.0.0"
    location: "http://example.com"
    namespace: "longhorn-system"
  images:
    - name: "nginx"
      version: "1.19.6"
      location: "docker.io/nginx:1.19.6"
`

	invalidYAML := `
apiVersion: 1.0
releaseVersion: "v1.0.0"
supoortedUpgrades:
  - "v1.1.0"
components:
  operatingSystem:
    upgrade:
      version: "v1.0.0"
  kubernetes:
    rke2:
      version: "v1.20.0"
  longhorn:
    version: "v1.0.0"
    location: "http://example.com"
    namespace: "longhorn-system"
  images:
    - name: "nginx"
      version: "1.19.6"
`

	tmpFile, err := os.CreateTemp("", "valid_release_manifest_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(validYAML)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	// Test reading the valid YAML file
	t.Run("Valid YAML", func(t *testing.T) {
		manifest, err := ReadReleaseManifest(tmpFile.Name())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if manifest.APIVersion != 1.0 {
			t.Errorf("expected apiVersion 1.0, got %v", manifest.APIVersion)
		}
		if manifest.ReleaseVersion != "v1.0.0" {
			t.Errorf("expected releaseVersion v1.0.0, got %v", manifest.ReleaseVersion)
		}
		if len(manifest.SupoortedUpgrades) != 1 || manifest.SupoortedUpgrades[0] != "v1.1.0" {
			t.Errorf("expected supoortedUpgrades [v1.1.0], got %v", manifest.SupoortedUpgrades)
		}
		if manifest.Components.OperatingSystem.Upgrade.Version != "v1.0.0" {
			t.Errorf("expected operatingSystem upgrade version v1.0.0, got %v", manifest.Components.OperatingSystem.Upgrade.Version)
		}
		if manifest.Components.Kubernetes.Rke2.Version != "v1.20.0" {
			t.Errorf("expected kubernetes rke2 version v1.20.0, got %v", manifest.Components.Kubernetes.Rke2.Version)
		}
		if manifest.Components.Helm[0].Version != "v1.0.0" || manifest.Components.Helm[0].Location != "http://example.com" || manifest.Components.Helm[0].Namespace != "longhorn-system" {
			t.Errorf("expected longhorn with version v1.0.0, location http://example.com, and namespace longhorn-system, got %v", manifest.Components.Helm)
		}
		if len(manifest.Components.Images) != 1 || manifest.Components.Images[0].Name != "nginx" || manifest.Components.Images[0].Version != "1.19.6" || manifest.Components.Images[0].Location != "docker.io/nginx:1.19.6" {
			t.Errorf("expected images with name nginx, version 1.19.6, and location docker.io/nginx:1.19.6, got %v", manifest.Components.Images)
		}
	})

	// Create a temporary invalid YAML file for testing
	tmpInvalidFile, err := os.CreateTemp("", "invalid_release_manifest_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpInvalidFile.Name())

	if _, err := tmpInvalidFile.Write([]byte(invalidYAML)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	if err := tmpInvalidFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	// Test reading the invalid YAML file
	t.Run("Invalid YAML", func(t *testing.T) {
		_, err := ReadReleaseManifest(tmpInvalidFile.Name())
		if err == nil {
			t.Fatal("expected an error, got none")
		}
	})

	// Test reading a non-existent file
	t.Run("Non-existent file", func(t *testing.T) {
		_, err := ReadReleaseManifest("non_existent_file.yaml")
		if err == nil {
			t.Fatal("expected an error, got none")
		}
	})
}
