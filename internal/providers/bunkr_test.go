package providers_test

import (
	"testing"

	"github.com/nakrovati/fapesnap/internal/providers"
)

func TestBunkrService_GetCollectionStringFromURL(t *testing.T) {
	t.Parallel()

	provider := &providers.BunkrProvider{}
	provider.InitProvider()

	tests := []providers.TestCase{
		{"https://bunkrrr.org/a/album", "album", false},
		{"https://bunkrrr.org/a/album/", "album", false},
		{"https://bunkr.si/a/album", "album", false},
		{"https://otherdomain.com/album", "", true},
		{"https://otherdomain.com/a/album", "", true},
		{"https://bunkrrr.org/", "", true},
		{"https://bunkrrr.org", "", true},
		{"https://sub.bunkrrr.org/album", "", true},
		{"https://sub.bunkrrr.org/a/album", "album", false},
	}

	providers.RunProviderTests(t, provider, tests)
}
