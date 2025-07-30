package helm

import (
	"fmt"
	"github.com/alknopfler/seactl/pkg/registry"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	tempDir = "./"
)

type Helm struct {
	Name     string // release name (e.g., "rancher")
	Chart    string // chart name or full OCI reference
	Version  string
	URL      string // optional repo URL (for HTTPS charts)
	TmpDir   string
	Insecure bool
	reg      *registry.Registry
}

func New(name, version, chart, url string, reg *registry.Registry) *Helm {
	return &Helm{
		Name:    name,
		Version: version,
		Chart:   chart,
		URL:     url,
		reg:     reg,
	}
}

func (h *Helm) Download() error {
	var args []string

	// Determine chart reference
	if strings.HasPrefix(h.Chart, "oci://") {
		// OCI chart: full reference is already in h.Chart
		args = []string{"pull", h.Chart, "--version", h.Version, "-d", tempDir}
	} else {
		if h.URL == "" {
			return fmt.Errorf("repository URL is missing for chart %s", h.Name)
		}

		// Regular Helm repo chart
		args = []string{
			"pull", h.Chart,
			"--repo", strings.TrimSuffix(h.URL, "/"),
			"--version", h.Version,
			"-d", tempDir,
		}
	}
	// Execute the command
	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		log.Printf("failed to download chart: %v", err)
		return err
	}

	return nil
}

func (h *Helm) Verify() error {
	if _, err := os.Stat(filepath.Join(tempDir, h.Name+"-"+h.Version+".tgz")); os.IsNotExist(err) {
		log.Printf("file does not exist to be verified %s", err.Error())
		return err
	}
	return nil
}

func (h *Helm) Upload() error {
	var args []string
	args = append(args, "push", filepath.Join(tempDir, h.Name+"-"+h.Version+".tgz"), "oci://"+h.reg.RegistryURL)

	if h.Insecure {
		args = append(args, "--insecure-skip-tls-verify")
	} else if h.reg.RegistryCACert != "" {
		args = append(args, "--ca-file", h.reg.RegistryCACert)
	}

	cmd := exec.Command("helm", args...)
	err := cmd.Run()
	if err != nil {
		log.Printf("failed to push to the registry: %s", err)
		return err
	}
	defer os.Remove(filepath.Join(tempDir, h.Name+"-"+h.Version+".tgz"))
	return nil
}
