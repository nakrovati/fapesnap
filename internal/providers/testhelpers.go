package providers

import (
	"testing"
)

type TestCase struct {
	URL          string
	ExpectedUser string
	ExpectError  bool
}

func RunProviderTests(t *testing.T, provider Provider, tests []TestCase) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.URL, func(t *testing.T) {
			t.Parallel()

			result, err := provider.GetCollectionFromURL(test.URL)

			if test.ExpectError {
				if err == nil {
					t.Errorf("Expected error for URL %s, but got nil", test.URL)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for URL %s: %v", test.URL, err)
				}

				if result != test.ExpectedUser {
					t.Errorf("Expected %s, got %s", test.ExpectedUser, result)
				}
			}
		})
	}
}
