package providers_test

import (
	"testing"

	"github.com/nakrovati/fapesnap/internal/providers"
)

func TestBunkrService_GetCollectionStringFromURL(t *testing.T) {
	provider := &providers.BunkrProvider{}
	provider.InitProvider()

	tests := []struct {
		url          string
		expectedUser string
		expectError  bool
	}{
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

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			result, err := provider.GetCollectionFromURL(test.url)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error for URL %s, but got nil", test.url)
				}
			} else {
				if err != nil {
					t.Errorf("Not expected error for URL %s, but got %v", test.url, err)
				}

				if result != test.expectedUser {
					t.Errorf("Expected %s, but got %s", test.expectedUser, result)
				}
			}
		})
	}
}
