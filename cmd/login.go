package cmd

import (
	"github.com/alknopfler/seactl/pkg/registry"
	"github.com/spf13/cobra"
)

func NewRegistryLoginCommand() *cobra.Command {
	var username, password string
	c := &cobra.Command{
		Use:   "login",
		Short: "Commands to Login into a private registry",
		RunE: func(cmd *cobra.Command, args []string) error {

			r := registry.NewRegistryLogin()
			return r.Login()
		},
	}
	flags := c.Flags()

	flags.StringVarP(&username, "user", "u", "", "Registry Username")
	flags.StringVarP(&password, "password", "p", "", "Registry Password")
	c.MarkFlagRequired("user")
	c.MarkFlagRequired("password")
	return c
}
