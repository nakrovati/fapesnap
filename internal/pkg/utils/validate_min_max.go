package utils

import (
	"errors"
)

var (
	ErrMinMaxLessThanOne      = errors.New("min and max cannot be less than 1")
	ErrMinMaxGreaterThanLimit = errors.New("min and max cannot be greater than 100000")
	ErrMinGreaterThanMax      = errors.New("min cannot be greater than max")
)

func ValidateMinMax(minValue int, maxValue int) error {
	if minValue < 1 || maxValue < 1 {
		return ErrMinMaxLessThanOne
	}

	if minValue > 100000 || maxValue > 100000 {
		return ErrMinMaxGreaterThanLimit
	}

	if minValue > maxValue {
		return ErrMinGreaterThanMax
	}

	return nil
}
