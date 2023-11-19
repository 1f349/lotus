package utils

func MapKeys[M ~map[K]any, K comparable](src M) []K {
	dst := make([]K, 0, len(src))
	for k := range src {
		dst = append(dst, k)
	}
	return dst
}
