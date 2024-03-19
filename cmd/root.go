package cmd

import (
	"fapesnap/pkg/downloader"
	"fapesnap/pkg/providers/fapello"
	"fapesnap/pkg/providers/fapodrop"
	"log"

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
			fapodropProvider := fapodrop.FapodropProvider{}
			fapodropProvider.InitProvider()

			downloader := downloader.Downloader{PhotosProvider: &fapodropProvider}
			err := downloader.DownloadPhotos(userName, min, max)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	fapelloCmd = &cobra.Command{
		Use:   "fapello",
		Short: "Download photos from fapello",
		Run: func(cmd *cobra.Command, args []string) {
			fapelloProvider := fapello.FapelloProvider{}
			fapelloProvider.InitProvider()

			downloader := downloader.Downloader{PhotosProvider: &fapelloProvider}
			err := downloader.DownloadPhotos(userName, min, max)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(fapodropCmd)
	rootCmd.AddCommand(fapelloCmd)

	fapodropCmd.Flags().StringVarP(&userName, "username", "u", "", "Profile name in fapodrop")
	fapodropCmd.Flags().IntVarP(&min, "min", "", 1, "Minimum photo ID")
	fapodropCmd.Flags().IntVarP(&max, "max", "", 100000, "Maximum photo ID")
	if err := fapodropCmd.MarkFlagRequired("username"); err != nil {
		log.Fatal(err)
	}

	fapelloCmd.Flags().StringVarP(&userName, "username", "u", "", "Profile name in fapello")
	fapelloCmd.Flags().IntVarP(&min, "min", "", 1, "Minimum photo ID")
	fapelloCmd.Flags().IntVarP(&max, "max", "", 100000, "Maximum photo ID")
	if err := fapelloCmd.MarkFlagRequired("username"); err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
