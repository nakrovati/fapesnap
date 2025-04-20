package providers_test

import (
	"testing"

	"github.com/nakrovati/fapesnap/internal/providers"
)

func TestFapelloService_GetCollectionStringFromURL(t *testing.T) {
	t.Parallel()

	provider := &providers.FapelloProvider{}
	provider.InitProvider()

	tests := []providers.TestCase{
		{"https://fapello.com/username", "username", false},
		{"https://fapello.com/username/", "username", false},
		{"https://otherdomain.com/username", "", true},
		{"https://fapello.com/", "", true},
		{"https://fapello.com", "", true},
		{"https://sub.fapello.com/username", "", true},
	}

	providers.RunProviderTests(t, provider, tests)
}
