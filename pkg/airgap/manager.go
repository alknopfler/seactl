package airgap

import (
	"errors"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/seactl/pkg/config"
	"github.com/alknopfler/seactl/pkg/helm"
	"github.com/alknopfler/seactl/pkg/images"
	"github.com/alknopfler/seactl/pkg/registry"
	"github.com/alknopfler/seactl/pkg/rke2"
	"log"
	"os/exec"
	"sync"
)

type Manager interface {
	Download() error
	Verify() error
	Upload() error
}

func GenerateAirGapEnvironment(releaseVersion, releaseMode, registryURL, registryAuthFile, registryCACert, outputDirTarball string, insecure bool) error {
	fatalErrors := make(chan error)
	wgDone := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(3)

	releaseManifest, imagesManifest, err := config.ReadAirgapManifest(releaseVersion, releaseMode)
	if err != nil {
		return err
	}

	reg := registry.New(registryAuthFile, registryURL, registryCACert, insecure)

	// RKE2 Artifacts
	go func() {
		err = generateRKE2Artifacts(releaseManifest, outputDirTarball)
		if err != nil {
			fatalErrors <- err
		}
		wg.Done()
	}()

	// Helm Charts Artifacts to be uploaded to registry
	go func() {
		err = generateHelmArtifacts(releaseManifest, reg)
		if err != nil {
			fatalErrors <- err
		}
		wg.Done()
	}()

	// Images Artifacts to be uploaded to registry
	go func() {
		err = generateImagesArtifacts(imagesManifest, reg)
		if err != nil {
			fatalErrors <- err
		}
		wg.Done()
	}()

	// Wait until all the goroutines are done
	go func() {
		wg.Wait()
		close(wgDone)
	}()

	// Wait until either WaitGroup is done or an error is received through the channel
	select {
	case <-wgDone:
		// carry on
		break
	case err = <-fatalErrors:
		close(fatalErrors)
		log.Fatal("Error found running the program: ", err)
		return err
	}
	return nil
}

func generateRKE2Artifacts(airgapManifest *config.ReleaseManifest, outputDirTarball string) error {

	r := rke2.New(airgapManifest.Spec.Components.Kubernetes.Rke2.Version, outputDirTarball)

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
	for _, value := range releaseManifest.Spec.Components.Workloads.Helm {

		h := helm.New(value.ReleaseName, value.Version, value.Chart, reg)
		err := reg.RegistryHelmLogin()
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

		log.Printf("Helm chart %s prepared and uploaded successfully!\n", value)

	}
	log.Println(color.InGreen("Helm Chart artifacts pre-loaded in registry successfully!"))
	return nil
}

func generateImagesArtifacts(imagesManifest *config.ImagesManifest, reg *registry.Registry) error {
	// Images Artifacts to be uploaded to registry
	for _, value := range imagesManifest.Images {

		image := images.New(value.Name, reg)
		err := reg.RegistryLogin()
		if err != nil {
			return err
		}

		log.Printf("Starting to download images %s. This may take a while...", value)
		err = image.Download()
		if err != nil {
			return err
		}

		log.Printf("Starting to verify images %s. This may take a while...", value)
		err = image.Verify()
		if err != nil {
			return err
		}

		log.Printf("Starting to upload images %s to the registry %s This may take a while...", value, reg.RegistryURL)
		if reg.RegistryInsecure {
			image.Insecure = true
		}
		err = image.Upload()
		if err != nil {
			return err
		}

		log.Printf("Images %s prepared and uploaded successfully!\n", value)

	}
	log.Println(color.InGreen("Images artifacts pre-loaded in registry successfully!"))
	return nil

}

func CheckHelmCommand() error {
	if _, err := exec.LookPath("helm"); err != nil {
		return errors.New("Helm command not found in the system. You need to install it to continue")
	}
	return nil
}
