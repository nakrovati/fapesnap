package main

import (
	"fapesnap/pkg/providers/fapello"
	"fapesnap/pkg/providers/fapodrop"
	"fapesnap/pkg/utils"
	"log"

	"github.com/spf13/cobra"
)

type Provider interface {
	DownloadPhotos(userName string, min int, max int) error
}

var (
	userName string
	min      int
	max      int
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

			err := utils.ValidateMinMax(min, max)
			if err != nil {
				log.Fatal(err)
			}

			var fapodropProvider Provider = &fapodrop.FapodropProvider{}
			fapodropProvider.DownloadPhotos(userName, min, max)
		},
	}
	fapelloCmd = &cobra.Command{
		Use:   "fapello",
		Short: "Download photos from fapello",
		Run: func(cmd *cobra.Command, args []string) {
			if userName == "" {
				log.Fatal("You must specify a username")
			}

			err := utils.ValidateMinMax(min, max)
			if err != nil {
				log.Fatal(err)
			}

			var fapelloProvider Provider = &fapello.FapelloProvider{}
			fapelloProvider.DownloadPhotos(userName, min, max)
		},
	}
)

func init() {
	rootCmd.AddCommand(fapodropCmd)
	rootCmd.AddCommand(fapelloCmd)

	fapodropCmd.PersistentFlags().StringVarP(&userName, "username", "u", "", "Profile name in fapodrop")
	fapodropCmd.PersistentFlags().IntVarP(&min, "min", "", 1, "Minimum photo ID")
	fapodropCmd.PersistentFlags().IntVarP(&max, "max", "", 100000, "Maximum photo ID")

	fapelloCmd.PersistentFlags().StringVarP(&userName, "username", "u", "", "Profile name in fapello")
	fapelloCmd.PersistentFlags().IntVarP(&min, "min", "", 1, "Minimum photo ID")
	fapelloCmd.PersistentFlags().IntVarP(&max, "max", "", 100000, "Maximum photo ID")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
