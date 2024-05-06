package utils

import (
	"errors"
)

var (
	ErrMinMaxLessThanOne      = errors.New("min and max cannot be less than 1")
	ErrMinMaxGreaterThanLimit = errors.New("min and max cannot be greater than 100000")
	ErrMinGreaterThanMax      = errors.New("min cannot be greater than max")
)

func ValidateMinMax(min int, max int) error {
	if min < 1 || max < 1 {
		return ErrMinMaxLessThanOne
	}

	if min > 100000 || max > 100000 {
		return ErrMinMaxGreaterThanLimit
	}

	if min > max {
		return ErrMinGreaterThanMax
	}

	return nil
}
