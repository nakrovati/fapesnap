package main

import (
	"fapesnap/cmd/fapesnap/providers"

	"github.com/spf13/cobra"
)

func initRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "root",
		Short: "Download photos from fapello/fapodrop",
	}

	rootCmd.AddCommand(providers.InitFapelloCmd())
	rootCmd.AddCommand(providers.InitFapodropCmd())
	rootCmd.AddCommand(providers.InitBunkrCmd())

	return rootCmd
}

func main() {
	rootCmd := initRootCmd()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
