package cmd

import (
	"fmt"

	"github.com/alknopfler/seactl/pkg/airgap"
	"github.com/alknopfler/seactl/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	releaseVersion      string
	releaseMode         string
	registryAuthFile    string
	rancherAppsAuthFile string
	susePrivateRegistryAuthFile string
	registryURL         string
	registryCACert      string
	registryInsecure    bool
	outputDirTarball    string
	dryRun              bool
	debug               bool
	skipRke2            bool
	skipHelm            bool
	skipImages          bool
)

func NewAirGapCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "mirror",
		Short: "Command to mirror the air-gap artifacts from the airgap manifest",
		Example: `  # Re-run only the Helm charts after fixing a failed chart, skipping RKE2 and images
  seactl mirror -v 3.4.0 -m production -o /tmp/airgap -a registry-auth.txt \
    -r myregistry:5000 --rancher-apps-authfile rancher-auth.txt \
    --skip-rke2 --skip-images`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set global debug flag
			logger.Debug = debug

			// Check helm
			if err := airgap.CheckHelmCommand(); err != nil {
				return err
			}

			// Validate flags only if release-version is present (Binary mode)
			if releaseVersion != "" {
				// Validate release mode
				if releaseMode != "factory" && releaseMode != "production" {
					return fmt.Errorf("invalid value for --release-mode: %s, allowed: 'factory' or 'production'", releaseMode)
				}
				// Validate release version format X.Y.Z
				if len(releaseVersion) < 5 || releaseVersion[1] != '.' || releaseVersion[3] != '.' {
					return fmt.Errorf("invalid release version format: %s, expected format X.Y.Z", releaseVersion)
				}
			} else {
				// Container mode: validation is implicit (files must exist)
				// We clear releaseMode to ensure it's empty so ReadAirgapManifest knows to use local files
				releaseMode = ""
			}

			// Call airgap generation
			return airgap.GenerateAirGapEnvironment(
				dryRun, releaseVersion, releaseMode,
				registryURL, registryAuthFile, rancherAppsAuthFile, susePrivateRegistryAuthFile, registryCACert,
				outputDirTarball, registryInsecure,
				skipRke2, skipHelm, skipImages,
			)
		},
	}

	flags := c.Flags()
	flags.StringVarP(&releaseVersion, "release-version", "v", "", "SUSE Edge release version (X.Y.Z)")
	flags.StringVarP(&releaseMode, "release-mode", "m", "factory", "Release mode: factory or production")
	flags.StringVarP(&registryURL, "registry-url", "r", "", "Registry URL")
	flags.StringVarP(&registryCACert, "registry-cacert", "c", "", "Registry CA Certificate")
	flags.StringVarP(&registryAuthFile, "registry-authfile", "a", "", "Registry Auth file (base64 user:pass)")
	flags.StringVar(&rancherAppsAuthFile, "rancher-apps-authfile", "", "Rancher Apps registry auth file (base64 user:pass)")
	flags.StringVar(&susePrivateRegistryAuthFile, "suse-private-registry-authfile", "", "SUSE Private Registry auth file (base64 user:pass)")
	flags.BoolVarP(&registryInsecure, "insecure", "k", false, "Skip TLS verification")
	flags.StringVarP(&outputDirTarball, "output", "o", "", "Output directory for tarball files")
	flags.BoolVarP(&dryRun, "dry-run", "d", false, "Dry run mode")
	flags.BoolVar(&debug, "debug", false, "Enable debug logging")
	flags.BoolVar(&skipRke2, "skip-rke2", false, "Skip downloading and verifying RKE2 artifacts")
	flags.BoolVar(&skipHelm, "skip-helm", false, "Skip downloading and uploading Helm chart artifacts")
	flags.BoolVar(&skipImages, "skip-images", false, "Skip downloading and uploading image artifacts")

	// Required flags
	c.MarkFlagRequired("output")
	c.MarkFlagRequired("registry-url")
	c.MarkFlagRequired("rancher-apps-authfile")

	return c
}
