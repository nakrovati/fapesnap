package main

import (
	"fapesnap/pkg/providers/fapello"
	"fapesnap/pkg/providers/fapodrop"
	"log"

	"github.com/spf13/cobra"
)

type Provider interface {
	DownloadPhotos(userName string) error
}

var (
	userName string
	rootCmd  = &cobra.Command{
		Use:   "root",
		Short: "Download photos from fapello/fapodrop",
	}
	fapodropCmd = &cobra.Command{
		Use:   "fapodrop",
		Short: "Download photos from fapodrop",
		Run: func(cmd *cobra.Command, args []string) {
			if userName == "" {
				log.Fatal("You must specify a username")
			}
			var fapodropProvider Provider = &fapodrop.FapodropProvider{}
			fapodropProvider.DownloadPhotos(userName)
		},
	}
	fapelloCmd = &cobra.Command{
		Use:   "fapello",
		Short: "Download photos from fapello",
		Run: func(cmd *cobra.Command, args []string) {
			if userName == "" {
				log.Fatal("You must specify a username")
			}
			var fapelloProvider Provider = &fapello.FapelloProvider{}
			fapelloProvider.DownloadPhotos(userName)
		},
	}
)

func init() {
	rootCmd.AddCommand(fapodropCmd)
	rootCmd.AddCommand(fapelloCmd)
	fapodropCmd.PersistentFlags().StringVarP(&userName, "username", "u", "", "Profile name in fapodrop")
	fapelloCmd.PersistentFlags().StringVarP(&userName, "username", "u", "", "Profile name in fapello")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
