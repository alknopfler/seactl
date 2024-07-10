package airgap

import (
	"github.com/TwiN/go-color"
	"github.com/alknopfler/seactl/pkg/config"
	"github.com/alknopfler/seactl/pkg/helm"
	"github.com/alknopfler/seactl/pkg/registry"
	"github.com/alknopfler/seactl/pkg/rke2"
	"log"
)

type Manager interface {
	Download() error
	Verify() error
	Upload() error
}

func GenerateAirGapEnvironment(releaseManifestFile, registryURL, registryUsername, registryPassword, registryCACert, outputDirTarball string, insecure bool) error {

	releaseManifest, err := config.ReadReleaseManifest(releaseManifestFile)
	if err != nil {
		return err
	}

	reg := registry.New(registryUsername, registryPassword, registryURL, registryCACert, insecure)

	// RKE2 Artifacts
	//err = generateRKE2Artifacts(releaseManifest, outputDirTarball)
	//if err != nil {
	//	return err
	//}

	// Helm Charts Artifacts to be uploaded to registry

	err = generateHelmArtifacts(releaseManifest, reg)

	return nil
}

func generateRKE2Artifacts(releaseManifest *config.ReleaseManifest, outputDirTarball string) error {

	r := rke2.New(releaseManifest.Components.Kubernetes.Rke2.Version, outputDirTarball)

	log.Printf("Starting to download RKE2 images to %s. This may take a while...", outputDirTarball)

	err := r.Download()
	if err != nil {
		return err
	}

	err = r.Verify()
	if err != nil {
		return err
	}

	log.Println(color.InGreen("RKE2 Images downloaded and verified successfully! you can find them in: " + outputDirTarball))
	return nil
}

func generateHelmArtifacts(releaseManifest *config.ReleaseManifest, reg *registry.Registry) error {
	// Helm Charts Artifacts to be uploaded to registr
	for _, value := range releaseManifest.Components.Helm {

		h := helm.New(value.Name, value.Version, value.Location, reg)
		err := h.RegistryLogin(reg)
		if err != nil {
			return err
		}

		log.Printf("Starting to download helm-chart %s. This may take a while...", value)
		err = h.Download()
		if err != nil {
			return err
		}

		log.Printf("Starting to verify helm-chart %s. This may take a while...", value)
		err = h.Verify()
		if err != nil {
			return err
		}

		log.Printf("Starting to upload helm-chart %s to the registry %s This may take a while...", value, reg.RegistryURL)
		if reg.RegistryInsecure {
			h.Insecure = true
		}
		err = h.Upload()
		if err != nil {
			return err
		}

		log.Println("Helm chart %s prepared and uploaded successfully!", value)

	}
	log.Println(color.InGreen("Helm Chart artifacts pre-loaded in registry successfully!"))
	return nil

}
