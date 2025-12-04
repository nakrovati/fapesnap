package utils

func RoundUp(num int) int {
	if num < 1000 {
		return 1000
	}

	return ((num + 999) / 1000) * 1000
}
