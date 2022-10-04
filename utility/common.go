package utility

import (
	"golang.org/x/exp/constraints"
	"math"
)

// IntRound returns the integer resulting from rounding of x
func IntRound(x float64) int {
	return int(math.Round(x))
}

// AbsInt returns the absolute value of an integer
func AbsInt(x int) int {
	if x >= 0 {
		return x
	} else {
		return -x
	}
}

// Min returns the minimum of its arguments
func Min[T constraints.Ordered](vs ...T) T {
	vm := vs[0]
	for _, v := range vs {
		if v < vm {
			vm = v
		}
	}
	return vm
}

// Max returns the maximum of its arguments
func Max[T constraints.Ordered](vs ...T) T {
	vm := vs[0]
	for _, v := range vs {
		if v > vm {
			vm = v
		}
	}
	return vm
}

// Minimize returns the value that minimizes the function
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

// Maximize returns the value that maximizes the function
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

// Mean returns the mean of its arguments
func Mean(vs ...float64) float64 {
	var vm float64
	for _, v := range vs {
		vm += v
	}
	return vm / float64(len(vs))
}
