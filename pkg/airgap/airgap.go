package airgap

import (
	"fmt"
	"github.com/alknopfler/seactl/pkg/config"
)

func GenerateAirGap(releaseManifestFile, outputTarball string) error {
	releaseManifest, err := config.ReadReleaseManifest(releaseManifestFile)
	if err != nil {
		return err
	}
	fmt.Println(releaseManifest)
	return nil
}

func PreloadAirGapToRegistry(releaseManifestFile, registryUsername, registryPassword, registryURL, registryCACert string) error {
	config.ReadReleaseManifest(releaseManifestFile)
	return nil
}
