package utility

import (
	"golang.org/x/exp/constraints"
	"math"
)

func IntRound(x float64) int {
	return int(math.Round(x))
}

func AbsInt(x int) int {
	if x >= 0 {
		return x
	} else {
		return -x
	}
}

func Min[T constraints.Ordered](vs ...T) T {
	vm := vs[0]
	for _, v := range vs {
		if v < vm {
			vm = v
		}
	}
	return vm
}

func Max[T constraints.Ordered](vs ...T) T {
	vm := vs[0]
	for _, v := range vs {
		if v > vm {
			vm = v
		}
	}
	return vm
}

func Minimize[T any, R constraints.Ordered](vs []T, f func(T) R) (T, R) {
	vm := vs[0]
	fm := f(vm)
	for _, v := range vs {
		fv := f(v)
		if fv < fm {
			fm = fv
			vm = v
		}
	}
	return vm, fm
}

func Maximize[T any, R constraints.Ordered](vs []T, f func(T) R) (T, R) {
	vm := vs[0]
	fm := f(vm)
	for _, v := range vs {
		fv := f(v)
		if fv > fm {
			fm = fv
			vm = v
		}
	}
	return vm, fm
}

func Mean(vs ...float64) float64 {
	var vm float64
	for _, v := range vs {
		vm += v
	}
	return vm / float64(len(vs))
}
