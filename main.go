package main

import (
	"fapesnap/pkg/downloader"
	"fapesnap/pkg/providers/fapello"
	"fapesnap/pkg/providers/fapodrop"

	"github.com/spf13/cobra"
)

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
			fapodropProvider := fapodrop.FapodropProvider{
				ProviderName: "fapodrop",
				BaseURL:      "https://fapodrop.com",
			}

			downloader := downloader.Downloader{ProviderName: fapodropProvider.ProviderName, PhotosProvider: &fapodropProvider}
			downloader.DownloadPhotos(userName, min, max)
		},
	}
	fapelloCmd = &cobra.Command{
		Use:   "fapello",
		Short: "Download photos from fapello",
		Run: func(cmd *cobra.Command, args []string) {
			fapelloProvider := fapello.FapelloProvider{
				ProviderName: "fapello",
				BaseURL:      "https://fapello.com",
			}

			downloader := downloader.Downloader{ProviderName: fapelloProvider.ProviderName, PhotosProvider: &fapelloProvider}
			downloader.DownloadPhotos(userName, min, max)
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
