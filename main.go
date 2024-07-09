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
		Short: "SUSE Edge Air-gap tool enables to create an air-gap scenario using the suse-edge release manifest",
		Long: "SUSE Edge Air-gap tool enables to create an air-gap scenario using the suse-edge release manifest. The output could be a tarball, but also you could upload to a private registry.\n" +
			"Features: \n" +
			"- Read the SUSE Edge release manifest (from input file)\n" +
			"- Save artifacts to a tarball\n" +
			"- Login to a private registry\n" +
			"- Upload and preload the private registry with the artifcats\n",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(0)
		},
	}

	c.AddCommand(cmd.NewAirGapCommand())

	return c
}
