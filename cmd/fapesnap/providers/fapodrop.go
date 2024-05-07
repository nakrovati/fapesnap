package providers

import (
	"fapesnap/internal/downloader"
	"fapesnap/pkg/providers/fapodrop"
	"log"

	"github.com/spf13/cobra"
)

func InitFapodropCmd() *cobra.Command {
	const (
		MaxPhotoID = 100000
		MinPhotoID = 1
	)

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
