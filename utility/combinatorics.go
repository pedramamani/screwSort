package utility

// Permutations returns all possible permutations of the integers 0 through n-1 (inclusive)
func Permutations(n int) [][]int {
	var ps [][]int
	var permuteUpTo func([]int, int)
	permuteUpTo = func(xs []int, m int) {
		if m == 1 {
			p := make([]int, len(xs))
			copy(p, xs)
			ps = append(ps, p)
		} else {
			for i := 0; i < m; i++ {
				permuteUpTo(xs, m-1)
				if m%2 == 1 {
					xs[i], xs[m-1] = xs[m-1], xs[i]
				} else {
					xs[0], xs[m-1] = xs[m-1], xs[0]
				}
			}
		}
	}
	permuteUpTo(Range(n), n)
	return ps
}

// Range returns the ascending list of integers 0 through n-1 (inclusive)
func Range(n int) []int {
	xs := make([]int, n)
	for i := range xs {
		xs[i] = i
	}
	return xs
}
