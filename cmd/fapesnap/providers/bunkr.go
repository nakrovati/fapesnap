package providers

import (
	"fapesnap/internal/downloader"
	"fapesnap/pkg/providers/bunkr"
	"log"

	"github.com/spf13/cobra"
)

func InitBunkrCmd() *cobra.Command {
	bunkrCmd := &cobra.Command{
		Use:   "bunkr",
		Short: "Download photos from bunkr",
		Run: func(cmd *cobra.Command, _ []string) {
			album, _ := cmd.Flags().GetString("album")

			bunkrProvider := bunkr.Provider{Album: album}
			bunkrProvider.InitProvider()

			downloader := downloader.Downloader{PhotosProvider: &bunkrProvider}
			err := downloader.DownloadPhotos()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	bunkrCmd.Flags().StringP("album", "a", "", "Bunkr album ID")

	if err := bunkrCmd.MarkFlagRequired("album"); err != nil {
		log.Fatal(err)
	}

	return bunkrCmd
}
