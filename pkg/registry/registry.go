package registry

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Registry struct {
	RegistryAuthFile string
	RegistryURL      string
	RegistryCACert   string
	RegistryInsecure bool
}

func New(registryAuthFile, registryURL, registryCACert string, insecure bool) *Registry {
	return &Registry{
		RegistryAuthFile: registryAuthFile,
		RegistryURL:      registryURL,
		RegistryCACert:   registryCACert,
		RegistryInsecure: insecure,
	}
}

func (r *Registry) RegistryHelmLogin() error {
	var args, auth []string
	args = append(args, "registry", "login", r.RegistryURL)

	auth, err := r.GetUserFromAuthFile(r.RegistryAuthFile)
	if err == nil {
		if auth[0] != "" && auth[1] != "" {
			args = append(args, "--username", auth[0], "--password", auth[1])
		}
	}

	if r.RegistryInsecure {
		args = append(args, "--insecure")
	} else if r.RegistryCACert != "" {
		args = append(args, "--ca-file", r.RegistryCACert)
	}
	cmd := exec.Command("helm", args...)
	err = cmd.Run()

	if err != nil {
		log.Printf("failed to login to the registry: %s", err)
		return err
	}
	return nil
}

func (r *Registry) RegistryPodmanLogin() error {
	//Using bindings for podman login

}

func (r *Registry) GetUserFromAuthFile(filePath string) ([]string, error) {
	// Read the content of the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []string{}, fmt.Errorf("failed to read auth file: %w", err)
	}

	// Decode the base64 content
	decodedData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return []string{}, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	// Split the decoded string by ":" to get user and pass
	parts := strings.SplitN(string(decodedData), ":", 2)
	if len(parts) != 2 {
		return []string{}, fmt.Errorf("decoded data does not contain user:pass format")
	}

	// Return the user part
	return parts, nil
}
