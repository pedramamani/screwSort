package iters

func Permutations(vs []int) (ps [][]int) {
	var permuteUpTo func([]int, int)
	permuteUpTo = func(vs []int, n int) {
		if n == 1 {
			p := make([]int, len(vs))
			copy(p, vs)
			ps = append(ps, p)
		} else {
			for i := 0; i < n; i++ {
				permuteUpTo(vs, n-1)
				if n%2 == 1 {
					vs[i], vs[n-1] = vs[n-1], vs[i]
				} else {
					vs[0], vs[n-1] = vs[n-1], vs[0]
				}
			}
		}
	}
	permuteUpTo(vs, len(vs))
	return
}
