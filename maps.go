package functional

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func Keys[K comparable, V any](in map[K]V) []K {
	arr := make([]K, 0, len(in))
	for k := range in {
		arr = append(arr, k)
	}
	return arr
}

func Values[K comparable, V any](in map[K]V) []V {
	arr := make([]V, 0, len(in))
	for _, v := range in {
		arr = append(arr, v)
	}
	return arr
}

func Entries[K comparable, V any](in map[K]V) []Pair[K, V] {
	arr := make([]Pair[K, V], 0, len(in))
	for k, v := range in {
		arr = append(arr, Pair[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return arr
}
