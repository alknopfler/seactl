package images

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/alknopfler/seactl/pkg/registry"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"log"
	"net/http"
	"os"
)

type Images struct {
	Name     string
	Insecure bool // If true, skip TLS verification
	reg      *registry.Registry
	ImageRef v1.Image
}

var (
	remoteImage = remote.Image
	remoteWrite = remote.Write
)

func New(name string, reg *registry.Registry) *Images {
	return &Images{
		Name: name,

		reg: reg,
	}
}

func (i *Images) Download() error {
	ref, err := name.ParseReference(i.Name)
	if err != nil {
		log.Printf("failed to parse image reference %v", err)
		return err
	}

	fmt.Println(ref)

	img, err := remoteImage(ref)
	if err != nil {
		log.Printf("pulling image %q: %v", img, err)
		return err
	}

	i.ImageRef = img
	log.Printf("successfully pulled image %q", img)
	return nil
}

func (i *Images) Verify() error {
	// Verify the image
	return nil
}

func (i *Images) Upload() error {

	ref, err := name.ParseReference(i.Name)
	if err != nil {
		return fmt.Errorf("parsing reference %q: %v", i.Name, err)
	}

	opts, err := i.getRemoteOpts()
	if err != nil {
		return fmt.Errorf("getting remote options: %v", err)
	}

	err = remoteWrite(ref, i.ImageRef, opts...)
	if err != nil {
		return fmt.Errorf("pushing image %q: %v", i.ImageRef, err)
	}

	log.Printf("successfully pushed image %q", i.ImageRef)
	return nil
}

func (i *Images) getRemoteOpts() ([]remote.Option, error) {
	// Create a custom HTTP transport
	tlsConfig := &tls.Config{}

	if i.Insecure {
		tlsConfig.InsecureSkipVerify = true
	} else if i.reg.RegistryCACert != "" {
		// Load CA certificate
		caCert, err := os.ReadFile(i.reg.RegistryCACert)
		if err != nil {
			return nil, fmt.Errorf("reading CA certificate: %v", err)
		}

		// Create a CA certificate pool
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	authFile, err := i.reg.GetUserFromAuthFile()
	if err != nil {
		return nil, fmt.Errorf("reading auth file: %v", err)
	}

	// Create a registry authenticator
	auth := &authn.Basic{
		Username: authFile[0],
		Password: authFile[1],
	}

	remoteOpts := []remote.Option{
		remote.WithTransport(transport),
		remote.WithAuth(auth),
	}

	// Remote options with custom HTTP client and authenticator
	return remoteOpts, nil
}
