package utils

func Contains[T comparable](s []T, a T) bool {
	for _, b := range s {
		if b == a {
			return true
		}
	}
	return false
}