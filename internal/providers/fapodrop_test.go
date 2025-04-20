package providers_test

import (
	"testing"

	"github.com/nakrovati/fapesnap/internal/providers"
)

func TestFapodropService_GetCollectionStringFromURL(t *testing.T) {
	t.Parallel()

	provider := &providers.FapodropProvider{}
	provider.InitProvider()

	tests := []providers.TestCase{
		{"https://fapodrop.com/username", "username", false},
		{"https://fapodrop.com/username/", "username", false},
		{"https://otherdomain.com/username", "", true},
		{"https://fapodrop.com/", "", true},
		{"https://fapodrop.com", "", true},
		{"https://sub.fapodrop.com/username", "", true},
	}

	providers.RunProviderTests(t, provider, tests)
}
