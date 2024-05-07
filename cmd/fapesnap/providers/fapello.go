package providers

import (
	"fapesnap/internal/downloader"
	"fapesnap/pkg/providers/fapello"
	"log"

	"github.com/spf13/cobra"
)

func InitFapelloCmd() *cobra.Command {
	const (
		MaxPhotoID = 100000
		MinPhotoID = 1
	)

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
