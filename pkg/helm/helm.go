package helm

import (
	"github.com/alknopfler/seactl/pkg/registry"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	tempDir = ".temporal"
)

type Helm struct {
	Name     string
	Version  string
	URL      string
	TmpDir   string
	Insecure bool
	reg      *registry.Registry
}

func New(name, version, url string, reg *registry.Registry) *Helm {
	return &Helm{
		Name:    name,
		Version: version,
		URL:     url,
		reg:     reg,
	}

}

func (h *Helm) Download() error {
	var args []string
	os.Mkdir(tempDir, os.FileMode(os.ModeDir))

	args = append(args, "pull", h.URL+h.Name, "--version", h.Version, "-d", tempDir)
	cmd := exec.Command("helm", args...)
	err := cmd.Run()

	if err != nil {
		log.Printf("failed to login to the registry: %s", err)
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
	}
	cmd := exec.Command("helm", args...)
	err := cmd.Run()
	if err != nil {
		log.Printf("failed to push to the registry: %s", err)
		return err
	}
	defer os.RemoveAll(tempDir)
	return nil
}
