package registry

import (
	"log"
	"os/exec"
)

type Registry struct {
	RegistryUsername string
	RegistryPassword string
	RegistryURL      string
	RegistryCACert   string
	RegistryInsecure bool
}

func New(registryUsername, registryPassword, registryURL, registryCACert string, insecure bool) *Registry {
	return &Registry{
		RegistryUsername: registryUsername,
		RegistryPassword: registryPassword,
		RegistryURL:      registryURL,
		RegistryCACert:   registryCACert,
		RegistryInsecure: insecure,
	}
}

func (r *Registry) RegistryHelmLogin() error {
	var args []string
	args = append(args, "registry", "login", r.RegistryURL)

	if r.RegistryUsername != "" && r.RegistryPassword != "" {
		args = append(args, "--username", r.RegistryUsername, "--password", r.RegistryPassword)
	}

	if r.RegistryInsecure {
		args = append(args, "--insecure")
	} else if r.RegistryCACert != "" {
		args = append(args, "--ca-file", r.RegistryCACert)
	}
	cmd := exec.Command("helm", args...)
	err := cmd.Run()

	if err != nil {
		log.Printf("failed to login to the registry: %s", err)
		return err
	}
	return nil
}

func (r *Registry) RegistryPodmanLogin() error {
	//Using bindings for podman login

}
