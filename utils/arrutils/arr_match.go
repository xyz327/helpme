package arrutils

func AllMatch[T interface{}](arr []T, predicate FilterFunc[T]) bool {
	ch := make(chan bool, len(arr))
	for _, v := range arr {
		go func(val T) {
			ch <- predicate(val)
		}(v)
	}
	for range arr {
		if !<-ch {
			return false
		}
	}
	return true
}

func AnyMatch[T interface{}](arr []T, predicate FilterFunc[T]) bool {
	ch := make(chan bool, len(arr))
	for _, v := range arr {
		go func(val T) {
			ch <- predicate(val)
		}(v)
	}
	for range arr {
		if <-ch {
			return true
		}
	}
	return false
}
