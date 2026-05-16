package providers

type MediaType string

const (
	MediaTypeImage   MediaType = "image"
	MediaTypeVideo   MediaType = "video"
	MediaTypeFile    MediaType = "file"
	MediaTypeUnknown MediaType = "unknown"
)

type Media struct {
	Type         MediaType `json:"type"`
	URL          string    `json:"url"`
	ThumbnailURL string    `json:"thumbnailUrl,omitempty"`
	Name         string    `json:"name,omitempty"`
	Size         string    `json:"size,omitempty"`
}

type Provider interface {
	FetchMediaItems(collectionSlug string) ([]Media, error)
	GetCollectionFromURL(url string) (string, error)
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
