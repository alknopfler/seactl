package registry

import (
	"context"
	"k8s.io/helm/pkg/helm"
)

// Registry is a struct that represents a registry, containing the following fields:
// - username: a string that represents the username of the registry
// - password: a string that represents the password of the registry
// - url: a string that represents the URL of the registry
// - cacert: a string that represents the CA certificate of the registry
type Registry struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	URL      string `json:"url" yaml:"url"`
	CACert   string `json:"cacert" yaml:"cacert"`
}

// New is a function that returns a new instance of the Registry struct
func New() *Registry {
	return &Registry{}
}

// Login is a interface method that returns an error
func (r *Registry) Login() (context.Context, error) {
	// Login using the registry credentials and podman bindings
	ctx := context.Background()
	helmClient := helm.NewClient()
	helmClient.Login(r.Username, r.Password, r.URL, r.CACert)
	return ctx,

}
