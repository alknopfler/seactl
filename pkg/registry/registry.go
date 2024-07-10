package registry

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
