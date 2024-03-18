package utils

import "fmt"

func ValidateMinMax(min int, max int) error {
	if min > max {
		return fmt.Errorf("min cannot be greater than max")
	}

	if min < 1 || max < 1 {
		return fmt.Errorf("min and max cannot be less than 1")
	}

	if min > 100000 || max > 100000 {
		return fmt.Errorf("max cannot be greater than 100000")
	}

	return nil
}
