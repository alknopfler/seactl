package registry

type Registry struct{}

func NewRegistry() Registry {
	return Registry{}
}

func NewRegistryLogin() Registry {
	return Registry{}
}

func (*Registry) Login() error {
	return nil
}
