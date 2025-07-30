package cmd

import (
	"fmt"
	"github.com/alknopfler/seactl/pkg/airgap"
	"github.com/spf13/cobra"
	"os"
)

var (
	releaseVersion   string
	releaseMode      string // This variable is provided to use `factory` or `production` mode
	registryAuthFile string
	registryURL      string
	registryCACert   string
	registryInsecure bool
	outputDirTarball string
	dryRun           bool // If true, do not perform any actions like upload files to registry or download nothing, just print what would be done
)

func NewAirGapCommand() *cobra.Command {

	c := &cobra.Command{
		Use:   "generate",
		Short: "Command to generate the air-gap artifacts from the airgap manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if helm command is available
			if err := airgap.CheckHelmCommand(); err != nil {
				return err
			}
			if releaseMode != "factory" && releaseMode != "production" {
				fmt.Fprintf(os.Stderr, "Invalid value for --release-mode: %s\n", releaseMode)
				fmt.Fprintf(os.Stderr, "Allowed values are 'factory' or 'production'\n")
				os.Exit(1)
			}
			if releaseVersion == "" || len(releaseVersion) < 5 || releaseVersion[1] != '.' || releaseVersion[3] != '.' {
				fmt.Printf("invalid release version format: %s, expected format is X.Y.Z", releaseVersion)
				os.Exit(1)
			}
			return airgap.GenerateAirGapEnvironment(dryRun, releaseVersion, releaseMode, registryURL, registryAuthFile, registryCACert, outputDirTarball, registryInsecure)
		},
	}
	// Add flags
	flags := c.Flags()
	flags.StringVarP(&releaseVersion, "release-version", "v", "", "SUSE Edge release version (e.g. 3.4.0 with X.Y.Z format)")
	flags.StringVarP(&releaseMode, "release-mode", "m", "factory", "Release mode, either 'factory' or 'production'")
	flags.StringVarP(&registryURL, "registry-url", "r", "", "Registry URL")
	flags.StringVarP(&registryCACert, "registry-cacert", "c", "", "Registry CA Certificate file")
	flags.StringVarP(&registryAuthFile, "registry-authfile", "a", "", "Registry Auth file with username:password base64 encoded")
	flags.BoolVarP(&registryInsecure, "insecure", "k", false, "Skip TLS verification in registry")
	flags.StringVarP(&outputDirTarball, "output", "o", "", "Output directory to store the tarball files")
	flags.BoolVarP(&dryRun, "dry-run", "d", false, "Dry run mode, do not perform any actions, just print what would be done")
	// add options and required flags
	c.MarkFlagRequired("release-version")
	c.MarkFlagRequired("output")
	c.MarkFlagRequired("registry-url")

	return c
}
