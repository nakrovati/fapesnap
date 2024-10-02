package providers

import (
	"testing"
)

func TestFapodropService_GetCollectionStringFromURL(t *testing.T) {
	provider := &FapodropProvider{}
	provider.InitProvider()

	tests := []struct {
		url          string
		expectedUser string
		expectError  bool
	}{
		{"https://fapodrop.com/username", "username", false},
		{"https://fapodrop.com/username/", "username", false},
		{"https://otherdomain.com/username", "", true},
		{"https://fapodrop.com/", "", true},
		{"https://fapodrop.com", "", true},
		{"https://sub.fapodrop.com/username", "", true},
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
