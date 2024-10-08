package providers_test

import (
	"testing"

	"github.com/nakrovati/fapesnap/internal/providers"
)

func TestFapelloService_GetCollectionStringFromURL(t *testing.T) {
	provider := &providers.FapelloProvider{}
	provider.InitProvider()

	tests := []struct {
		url          string
		expectedUser string
		expectError  bool
	}{
		{"https://fapello.com/username", "username", false},
		{"https://fapello.com/username/", "username", false},
		{"https://otherdomain.com/username", "", true},
		{"https://fapello.com/", "", true},
		{"https://fapello.com", "", true},
		{"https://sub.fapello.com/username", "", true},
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
