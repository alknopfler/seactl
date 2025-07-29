package airgap

import (
	"errors"
	"fmt"
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

func GenerateAirGapEnvironment(dryrun bool, releaseVersion, releaseMode, registryURL, registryAuthFile, registryCACert, outputDirTarball string, insecure bool) error {
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
		err := generateRKE2Artifacts(dryrun, releaseManifest, outputDirTarball)
		if err != nil {
			fatalErrors <- err
		}
		wg.Done()
	}()

	// Helm Charts Artifacts to be uploaded to registry
	go func() {
		err = generateHelmArtifacts(dryrun, releaseManifest, reg)
		if err != nil {
			fatalErrors <- err
		}
		wg.Done()
	}()

	// Images Artifacts to be uploaded to registry
	go func() {
		err = generateImagesArtifacts(dryrun, imagesManifest, reg)
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

func generateRKE2Artifacts(dryrun bool, airgapManifest *config.ReleaseManifest, outputDirTarball string) error {

	r := rke2.New(airgapManifest.Spec.Components.Kubernetes.Rke2.Version, outputDirTarball)

	if !dryrun {
		err := r.Download()
		if err != nil {
			return err
		}

		err = r.Verify()
		if err != nil {
			return err
		}
	} else {
		log.Println("Dry run mode enabled, skipping download and verification of RKE2 images.")
	}
	log.Println(color.InGreen("RKE2 Images downloaded and verified successfully! you can find them in: " + outputDirTarball))
	return nil
}

func generateHelmArtifacts(dryrun bool, releaseManifest *config.ReleaseManifest, reg *registry.Registry) error {
	// Helm Charts Artifacts to be uploaded to registr
	for _, value := range releaseManifest.Spec.Components.Workloads.Helm {

		h := helm.New(value.ReleaseName, value.Version, value.Chart, value.Repository, reg)
		if !dryrun {
			err := reg.RegistryHelmLogin()
			if err != nil {
				return err
			}

			err = h.Download()
			if err != nil {
				return err
			}

			err = h.Verify()
			if err != nil {
				return err
			}

			if reg.RegistryInsecure {
				h.Insecure = true
			}

			err = h.Upload()
			if err != nil {
				return err
			}

			log.Printf(color.InGreen("Helm chart %s prepared and uploaded successfully!\n"), value)
		} else {
			// list all info about the helm chart instead of uploading it
			log.Println("DryRun mode - Helm Chart Info:")
			log.Printf("\nName: %s\nVersion: %s\nURL: %s\nChart: %s\n", h.Name, h.Version, h.URL, h.Chart)
		}
	}
	log.Println(color.InGreen("Helm Chart artifacts pre-loaded in registry successfully!"))
	return nil
}

func generateImagesArtifacts(dryrun bool, imagesManifest *config.ImagesManifest, reg *registry.Registry) error {
	// Images Artifacts to be uploaded to registry
	for _, value := range imagesManifest.Images {

		image := images.New(value.Name, reg)
		if !dryrun {
			err := reg.RegistryLogin()
			if err != nil {
				return err
			}

			err = image.Download()
			if err != nil {
				return err
			}

			err = image.Verify()
			if err != nil {
				return err
			}

			if reg.RegistryInsecure {
				image.Insecure = true
			}

			// list all info about the image instead of uploading it
			fmt.Println("Image Info:")
			fmt.Printf("Name: %s\n", image.Name)

			err = image.Upload()
			if err != nil {
				return err
			}
		} else {
			// list all info about the image instead of uploading it
			log.Println("DryRun mode - Image Info:")
			log.Printf("\nName: %s\n", image.Name)
		}
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
