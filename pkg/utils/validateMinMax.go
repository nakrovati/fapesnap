package utils

import "fmt"

func ValidateMinMax(min int, max int) error {
	if min < 1 || max < 1 {
		return fmt.Errorf("min and max cannot be less than 1")
	}

	if min > 100000 || max > 100000 {
		return fmt.Errorf("min abd max cannot be greater than 100000")
	}

	if min > max {
		return fmt.Errorf("min (%d) cannot be greater than max (%d)", min, max)
	}

	return nil
}
