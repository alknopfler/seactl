package airgap

import (
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/seactl/pkg/config"
	"log"
)

type Downloader interface {
	Download() (string, error)
	Verify() error
}

func GenerateAirGapToTarball(releaseManifestFile, outputDirTarball string) error {
	releaseManifest, err := config.ReadReleaseManifest(releaseManifestFile)
	if err != nil {
		return err
	}
	rke2 := &RKE2{
		Version:          releaseManifest.Components.Kubernetes.Rke2.Version,
		OutputDirTarball: outputDirTarball,
	}

	log.Printf("Starting to download RKE2 images to the output directory %s. This may take a while...", outputDirTarball)
	err = rke2.Download()
	if err != nil {
		return err
	}

	err = rke2.Verify()
	if err != nil {
		return err
	}
	log.Println(color.InGreen("RKE2 Images downloaded and verified successfully! you can find them in the output directory: " + outputDirTarball))

	return nil
}

func PreloadAirGapToRegistry(releaseManifestFile, registryUsername, registryPassword, registryURL, registryCACert string) error {
	releaseManifest, err := config.ReadReleaseManifest(releaseManifestFile)
	if err != nil {
		return err
	}
	fmt.Println(releaseManifest)
	return nil
}
