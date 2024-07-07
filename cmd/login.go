package cmd

import (
	"errors"
	"github.com/alknopfler/seactl/pkg/registry"
	"github.com/spf13/cobra"
)

func NewRegistryLoginCommand() *cobra.Command {
	var username, password, url, cacert string
	c := &cobra.Command{
		Use:   "login",
		Short: "Command to Login into a private registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := registry.New()
			r.URL = url
			if username != "" {
				if password == "" {
					return errors.New("password is required when username is provided")
				}
				r.Username = username
				r.Password = password
			}
			if cacert != "" {
				r.CACert = cacert
			}
			return r.Login()
		},
	}
	flags := c.Flags()

	flags.StringVarP(&username, "user", "u", "", "Registry Username")
	flags.StringVarP(&password, "password", "p", "", "Registry Password")
	flags.StringVarP(&url, "url", "r", "", "Registry URL")
	flags.StringVarP(&cacert, "cacert", "c", "", "Registry CA Certificate file")
	c.MarkFlagRequired("url")

	return c
}
