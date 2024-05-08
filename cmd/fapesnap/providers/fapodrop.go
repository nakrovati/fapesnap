package providers

import (
	"fapesnap/internal/downloader"
	"fapesnap/internal/pkg/utils"
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
			username, _ := cmd.Flags().GetString("username")
			min, _ := cmd.Flags().GetInt("min")
			max, _ := cmd.Flags().GetInt("max")

			err := utils.ValidateMinMax(min, max)
			if err != nil {
				log.Fatal(err)
			}

			fapodropProvider := fapodrop.Provider{MaxPhotoID: max, MinPhotoID: min, Username: username}
			fapodropProvider.InitProvider()

			downloader := downloader.Downloader{PhotosProvider: &fapodropProvider}
			err = downloader.DownloadPhotos()
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
