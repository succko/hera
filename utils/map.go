package utils

func MapKeys[t any](m map[string]t) []string {
	keys := make([]string, len(m))
	i := 0
	for k, _ := range m {
		keys[i] = k
		i++
	}
	return keys
}
