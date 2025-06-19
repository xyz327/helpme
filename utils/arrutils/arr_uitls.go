package arrutils

import (
	"math/rand"
	"time"
)

type MapFuc[T any, R any] func(in T) R

// Map 转换集合对象 T[] -> R[]
func Map[T any, R any](list []T, mapper MapFuc[T, R]) []R {
	if list == nil {
		return nil
	}
	newList := make([]R, len(list))
	for i, _in := range list {
		newList[i] = mapper(_in)
	}
	return newList
}

type FlatMapFunc[T any, R any] func(in T) []R

// FlatMap 转换集合对象 T[] -> R[]
func FlatMap[T any, R any](list []T, mapper FlatMapFunc[T, R]) []R {
	newList := make([]R, 0)
	for _, _in := range list {
		newList = append(newList, mapper(_in)...)
	}
	return newList
}
func InSlice[E comparable](slice []E, item E) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type ToMapFunc[T any, K comparable] func(in T) K

// ToMap arr convert to Map
func ToMap[T any, K comparable](arr []T, keyMapper ToMapFunc[T, K]) map[K]T {
	res := make(map[K]T, len(arr))
	for _, _in := range arr {
		res[keyMapper(_in)] = _in
	}
	return res
}

func ChunkSlice[E uint32 | uint64](slice []E, chunkSize int) [][]E {
	var chunks [][]E
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func Shuffle[E any](slice []E) []E {
	shuffled := make([]E, len(slice))
	perm := rand.Perm(len(slice))
	for i, v := range perm {
		shuffled[v] = slice[i]
	}
	return shuffled
}

type GroupByFunc[T any, R comparable] func(item T) R

func GroupBy[T any, R comparable](list []T, mapper GroupByFunc[T, R]) map[R][]T {
	_new := make(map[R][]T)
	for _, _in := range list {
		k := mapper(_in)
		v, ok := _new[k]
		if !ok {
			v = make([]T, 0)
		}
		_new[k] = append(v, _in)
	}
	return _new
}

type FilterFunc[T any] func(item T) bool

func Filter[T any](list []T, filter FilterFunc[T]) []T {
	var res []T = make([]T, 0)
	for _, _in := range list {
		if filter(_in) {
			res = append(res, _in)
		}
	}
	return res
}

func RandPick[E any](slice []E) E {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(slice))
	return slice[index]

}

// Reduce
func Reduce[T any, R any](list []T, reducer func(acc R, in T) R, init R) R {
	var acc R = init
	for _, _in := range list {
		acc = reducer(acc, _in)
	}
	return acc
}

// Reverse 倒排集合
func Reverse[T any](list []T) []T {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}

// FindFirst 查找第一个满足条件的元素
func FindFirst[T any](list []T, filter FilterFunc[T]) (T, bool) {
	for _, _in := range list {
		if filter(_in) {
			return _in, true
		}
	}
	var empty T
	return empty, false
}

func Merge[T any](dest []T, source ...[]T) []T {
	for _, s := range source {
		dest = append(dest, s...)
	}
	return dest
}

// Distinct 去重
func Distinct[T comparable](list []T) []T {
	var res []T = make([]T, 0)
	m := make(map[T]bool)
	for _, _in := range list {
		if _, ok := m[_in]; !ok {
			m[_in] = true
			res = append(res, _in)
		}
	}
	return res
}

// Partition 分块
func Partition[T any](list []T, chunkSize int) [][]T {
	var chunks [][]T
	for i := 0; i < len(list); i += chunkSize {
		end := i + chunkSize
		if end > len(list) {
			end = len(list)
		}
		chunks = append(chunks, list[i:end])
	}
	return chunks
}

// ForEach
func ForEach[T any](list []T, f ConsumerFunc[T]) {
	for _, _in := range list {
		f(_in)
	}
}
