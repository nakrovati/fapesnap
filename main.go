package main

import (
	"fapesnap/pkg/providers/fapodrop"

	"github.com/spf13/cobra"
)

var (
	userName string
	rootCmd  = &cobra.Command{
		Use:   "fapodrop",
		Short: "Download photos from fapodrop",
		Run: func(cmd *cobra.Command, args []string) {
			fapodrop.DownloadPhotos(userName)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&userName, "username", "u", "", "Profile name in fapodrop/fapello")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
