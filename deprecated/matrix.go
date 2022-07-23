package deprecated

func altSign(n int) float64 {
	if n%2 == 0 {
		return 1
	}
	return -1
}

func Shape(mat [][]float64) (int, int) {
	return len(mat), len(mat[0])
}

func Zeros(m, n int) [][]float64 {
	ret := make([][]float64, m)
	for j := 0; j < m; j++ {
		ret[j] = make([]float64, n)
	}
	return ret
}

func Mean(mat [][]float64) float64 {
	var sum float64
	m, n := Shape(mat)
	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			sum += mat[j][i]
		}
	}
	return sum / float64(m*n)
}

func Transpose(mat [][]float64) [][]float64 {
	m, n := Shape(mat)
	ret := Zeros(n, m)
	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
			ret[j][i] = mat[i][j]
		}
	}
	return ret
}

func Add(mat, matO [][]float64) [][]float64 {
	m, n := Shape(mat)
	mo, no := Shape(matO)
	if m != mo || n != no {
		panic("incompatible dims")
	}
	ret := Zeros(m, n)
	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			ret[j][i] = mat[j][i] + matO[j][i]
		}
	}
	return ret
}

func Subtract(mat, matO [][]float64) [][]float64 {
	m, n := Shape(mat)
	mo, no := Shape(matO)
	if m != mo || n != no {
		panic("incompatible dims")
	}
	ret := Zeros(m, n)
	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			ret[j][i] = mat[j][i] - matO[j][i]
		}
	}
	return ret
}

func Multiply(mat, matO [][]float64) [][]float64 {
	m, n := Shape(mat)
	mo, no := Shape(matO)
	if n != mo {
		panic("incompatible dims")
	}
	ret := Zeros(m, no)
	for j := 0; j < m; j++ {
		for i := 0; i < no; i++ {
			for k := 0; k < n; k++ {
				ret[j][i] += mat[j][k] * matO[k][i]
			}
		}
	}
	return ret
}

func Scale(mat [][]float64, a float64) [][]float64 {
	m, n := Shape(mat)
	ret := Zeros(m, n)
	for j := 0; j < m; j++ {
		for i := 0; i < n; i++ {
			ret[j][i] = a * mat[j][i]
		}
	}
	return ret
}

func Determinant(mat [][]float64) float64 {
	m, n := Shape(mat)
	if n != m {
		panic("[][]float64 must be square")
	}
	switch n {
	case 0:
		return 1
	case 1:
		return mat[0][0]
	case 2:
		return mat[0][0]*mat[1][1] - mat[0][1]*mat[1][0]
	}
	sub := Zeros(n-1, n-1)
	var det float64
	var k, j, i, z int

	for k = 0; k < n; k++ {
		for j = 0; j < n-1; j++ {
			z = 0
			for i = 0; i < n; i++ {
				if i == k {
					continue
				}
				sub[j][z] = mat[j+1][i]
				z++
			}
		}
		det += altSign(k) * mat[0][k] * Determinant(sub)
	}
	return det
}

func Cofactor(mat [][]float64) [][]float64 {
	m, n := Shape(mat)
	if n != m {
		panic("[][]float64 must be square")
	}

	ret := Zeros(n, n)
	sub := Zeros(n-1, n-1)
	var i, j, p, q, x, y int

	for i = 0; i < n; i++ {
		for j = 0; j < n; j++ {
			p = 0
			for x = 0; x < n; x++ {
				if x == i {
					continue
				}
				q = 0
				for y = 0; y < n; y++ {
					if y == j {
						continue
					}
					sub[p][q] = mat[x][y]
					q++
				}
				p++
			}
			ret[i][j] = altSign(i+j) * Determinant(sub)
		}
	}
	return ret
}

func Inverse(mat [][]float64) [][]float64 {
	det := Determinant(mat)
	if det == 0 {
		panic("[][]float64 non-invertible since determinant is zero")
	}
	return Scale(Cofactor(Transpose(mat)), 1/det)
}
