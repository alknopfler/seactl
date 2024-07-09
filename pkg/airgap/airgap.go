package airgap

import (
	"fmt"
	"github.com/alknopfler/seactl/pkg/config"
)

var (
	rke2 = RKE2{}
	file string
)

type Downloader interface {
	Download() (string, error)
	Verify() error
}

func GenerateAirGap(releaseManifestFile, outputTarball string) error {
	releaseManifest, err := config.ReadReleaseManifest(releaseManifestFile)
	if err != nil {
		return err
	}

	rke2.Version = releaseManifest.Components.Kubernetes.Rke2.Version

	file, err = rke2.Download()
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
