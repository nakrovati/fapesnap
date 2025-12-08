package providers

type Provider interface {
	FetchPhotoURLs(collection string) ([]Photo, error)
	GetCollectionFromURL(url string) (string, error)
}

type Photo struct {
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
}

//nolint:ireturn
func GetProvider(providerName string) Provider {
	switch providerName {
	case "fapello":
		return NewFapelloProvider()
	case "fapodrop":
		return NewFapodropProvider()
	case "bunkr":
		return NewBunkrProvider()
	default:
		return nil
	}
}
