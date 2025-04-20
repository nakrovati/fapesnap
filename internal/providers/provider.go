package providers

type Provider interface {
	FetchPhotoURLs(collection string) ([]string, error)
	GetCollectionFromURL(url string) (string, error)
}

//nolint:ireturn
func GetProvider(providerName string) Provider {
	switch providerName {
	case "fapello":
		fapelloProvider := &FapelloProvider{MinPhotoID: 1, MaxPhotoID: 100000}
		fapelloProvider.InitProvider()

		return fapelloProvider
	case "fapodrop":
		fapodropProvider := &FapodropProvider{MinPhotoID: 1, MaxPhotoID: 100000}
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
