package cmd

import (
	"github.com/alknopfler/seactl/pkg/airgap"
	"github.com/spf13/cobra"
)

var (
	releaseManifestFile string
	registryAuthFile    string
	registryURL         string
	registryCACert      string
	registryInsecure    bool
	outputDirTarball    string
)

func NewAirGapCommand() *cobra.Command {

	c := &cobra.Command{
		Use:   "generate",
		Short: "Command to generate the air-gap artifacts from the release manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			return airgap.GenerateAirGapEnvironment(releaseManifestFile, registryURL, registryAuthFile, registryCACert, outputDirTarball, registryInsecure)
		},
	}
	// Add flags
	flags := c.Flags()
	flags.StringVarP(&releaseManifestFile, "input", "i", "", "Release manifest file")
	flags.StringVarP(&registryURL, "registry-url", "r", "", "Registry URL")
	flags.StringVarP(&registryCACert, "registry-cacert", "c", "", "Registry CA Certificate file")
	flags.StringVarP(&registryAuthFile, "registry-authfile", "a", "", "Registry Auth file with username:password base64 encoded")
	flags.BoolVarP(&registryInsecure, "insecure", "k", false, "Skip TLS verification in registry")
	flags.StringVarP(&outputDirTarball, "output", "o", "", "Output directory to store the tarball files")
	// add options and required flags
	c.MarkFlagRequired("input")
	c.MarkFlagRequired("output")
	c.MarkFlagRequired("registry-url")

	return c
}
