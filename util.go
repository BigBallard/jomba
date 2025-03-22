package jomba

type Comparater[T any]  func(T)bool

func SliceElementWhere[T any](slice []*T, comp Comparater[T]) (interface{}, bool) {
	for _, element := range slice {
		if match := comp(*element); match {
			return element, true
		}
	}
	return nil, false
}