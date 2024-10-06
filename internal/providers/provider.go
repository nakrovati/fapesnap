package providers

import "context"

type Provider interface {
	SetContext(ctx context.Context)
	FetchPhotoURLs(collection string) ([]string, error)
	GetCollectionFromURL(url string) (string, error)
}

func GetProvider(providerName string) Provider {
	switch providerName {
	case "fapello":
		fapelloProvider := &FapelloProvider{MaxPhotoID: 100000, MinPhotoID: 1}
		fapelloProvider.InitProvider()
		return fapelloProvider
	case "fapodrop":
		fapodropProvider := &FapodropProvider{MaxPhotoID: 100000, MinPhotoID: 1}
		fapodropProvider.InitProvider()
		return fapodropProvider
	case "bunkr":
		bunkrProvider := &BunkrProvider{}
		bunkrProvider.InitProvider()
		return bunkrProvider
	default:
		return nil
	}
}
