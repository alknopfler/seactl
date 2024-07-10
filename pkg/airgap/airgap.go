package airgap

import (
	"fmt"
	"github.com/alknopfler/seactl/pkg/config"
)

type Downloader interface {
	Download() (string, error)
	Verify() error
}

func GenerateAirGap(releaseManifestFile, outputDirTarball string) error {
	releaseManifest, err := config.ReadReleaseManifest(releaseManifestFile)
	if err != nil {
		return err
	}
	rke2 := &RKE2{
		Version:          releaseManifest.Components.Kubernetes.Rke2.Version,
		OutputDirTarball: outputDirTarball,
	}

	err = rke2.Download()
	if err != nil {
		return err
	}

	err = rke2.Verify()
	if err != nil {
		return err
	}

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
