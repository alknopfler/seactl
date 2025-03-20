package main

import (
	"github.com/TwiN/go-color"
	"github.com/alknopfler/seactl/cmd"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {

}

func main() {
	command := newCommand()
	if err := command.Execute(); err != nil {
		log.Fatalf(color.InRed("[ERROR] %s"), err.Error())
	}
}

func newCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "seactl",
		Short: "SUSE Edge Air-gap tool enables to create an air-gap scenario using the suse-edge airgap manifest",
		Long: "SUSE Edge Air-gap tool enables to create an air-gap scenario using the suse-edge airgap manifest. The output could be a tarball, but also you could upload to a private registry.\n" +
			"Features: \n" +
			"- Read the SUSE Edge airgap manifest (from input file)\n" +
			"- Save artifacts to a tarball\n" +
			"- Login to a private registry\n" +
			"- Upload and preload the private registry with the artifacts\n" +
			"\n" +
			"Example of airgap manifest yaml file: \n" +
			"---\n" +
			"apiVersion: 1.0\n" +
			"components:\n" +
			"  kubernetes:\n" +
			"    rke2:\n" +
			"      version: v1.28.9+rke2r1\n" +
			"  helm:\n" +
			"    - name: sriov-crd-chart\n" +
			"      version: 1.2.2\n" +
			"      location: oci://registry.suse.com/edge/\n" +
			"      namespace: sriov-network-operator\n" +
			"  images:\n" +
			"    - name: hardened-sriov-network-operator\n" +
			"      version: v1.2.0-build20240327\n" +
			"      location: docker.io/rancher",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(0)
		},
	}

	c.AddCommand(cmd.NewAirGapCommand())

	return c
}
