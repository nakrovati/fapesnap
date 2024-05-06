package cmd

import (
	"fapesnap/pkg/downloader"
	"fapesnap/pkg/providers/fapello"
	"fapesnap/pkg/providers/fapodrop"
	"log"

	"github.com/spf13/cobra"
)

const (
	MaxPhotoID = 100000
	MinPhotoID = 1
)

func initFapelloCmd() *cobra.Command {
	fapelloCmd := &cobra.Command{
		Use:   "fapello",
		Short: "Download photos from fapello",
		Run: func(cmd *cobra.Command, _ []string) {
			userName, _ := cmd.Flags().GetString("username")
			min, _ := cmd.Flags().GetInt("min")
			max, _ := cmd.Flags().GetInt("max")

			fapelloProvider := fapello.Provider{}
			fapelloProvider.InitProvider()

			downloader := downloader.Downloader{PhotosProvider: &fapelloProvider}
			err := downloader.DownloadPhotos(userName, min, max)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	fapelloCmd.Flags().StringP("username", "u", "", "Profile name in fapello")
	fapelloCmd.Flags().IntP("min", "", MinPhotoID, "Minimum photo ID")
	fapelloCmd.Flags().IntP("max", "", MaxPhotoID, "Maximum photo ID")

	if err := fapelloCmd.MarkFlagRequired("username"); err != nil {
		log.Fatal(err)
	}

	return fapelloCmd
}

func initFapodropCmd() *cobra.Command {
	fapodropCmd := &cobra.Command{
		Use:   "fapodrop",
		Short: "Download photos from fapodrop",
		Run: func(cmd *cobra.Command, _ []string) {
			userName, _ := cmd.Flags().GetString("username")
			min, _ := cmd.Flags().GetInt("min")
			max, _ := cmd.Flags().GetInt("max")

			fapodropProvider := fapodrop.Provider{}
			fapodropProvider.InitProvider()

			downloader := downloader.Downloader{PhotosProvider: &fapodropProvider}
			err := downloader.DownloadPhotos(userName, min, max)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	fapodropCmd.Flags().StringP("username", "u", "", "Profile name in fapodrop")
	fapodropCmd.Flags().IntP("min", "", MinPhotoID, "Minimum photo ID")
	fapodropCmd.Flags().IntP("max", "", MaxPhotoID, "Maximum photo ID")

	if err := fapodropCmd.MarkFlagRequired("username"); err != nil {
		log.Fatal(err)
	}

	return fapodropCmd
}

func initRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "root",
		Short: "Download photos from fapello/fapodrop",
	}

	rootCmd.AddCommand(initFapodropCmd())
	rootCmd.AddCommand(initFapelloCmd())

	return rootCmd
}

func Execute() {
	rootCmd := initRootCmd()

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
