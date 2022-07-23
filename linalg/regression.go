package linalg

import (
	"math"
	"screwSort/geometry"
)

func calcStatisticParams(ps []*geometry.Point) (xBar, yBar, x2Bar, y2Bar, xyBar float64) {
	n := float64(len(ps))
	for _, p := range ps {
		xBar += p.X()
		yBar += p.Y()
		x2Bar += math.Pow(p.X(), 2)
		y2Bar += math.Pow(p.Y(), 2)
		xyBar += p.X() * p.Y()
	}
	xBar /= n
	yBar /= n
	x2Bar /= n
	y2Bar /= n
	xyBar /= n
	return
}

func LinearRegression(ps []*geometry.Point) (*geometry.Line, float64) {
	xBar, yBar, x2Bar, y2Bar, xyBar := calcStatisticParams(ps)
	m := (xyBar - xBar*yBar) / (x2Bar - math.Pow(xBar, 2))
	y0 := yBar - m*xBar
	r := y2Bar + math.Pow(m, 2)*x2Bar + math.Pow(y0, 2) - 2*m*xyBar - 2*y0*yBar + 2*m*y0*xBar
	return geometry.LineMY0(m, y0), r
}

func OrthogonalRegression(ps []*geometry.Point) (*geometry.Line, float64) {
	xBar, yBar, x2Bar, y2Bar, xyBar := calcStatisticParams(ps)
	covXY := xyBar - xBar*yBar
	f := ((x2Bar - math.Pow(xBar, 2)) - (y2Bar - math.Pow(yBar, 2))) / covXY
	f2 := math.Pow(f, 2)
	s := -covXY / math.Abs(covXY)
	a := s * math.Sqrt(f2/(f2+4)+s*f/math.Sqrt(f2+4)+4/(f2+4)) / math.Sqrt2
	b := math.Sqrt(1 - math.Pow(a, 2))
	c := -a*xBar - b*yBar
	l := geometry.LineAC(a, c)
	//r := math.Sqrt(math.Pow(a, 2)*x2Bar + math.Pow(b, 2)*y2Bar + math.Pow(c, 2) + 2*a*b*xyBar + 2*a*c*xBar + 2*b*c*yBar)
	_, r := Maximize(ps, func(p *geometry.Point) float64 { return l.PerpDistanceTo(p) })
	return l, r // todo: experiment with different definitions of r (eg. maximum residual)
}

func Maximize[T any](vs []T, f func(T) float64) (T, float64) {
	maxV := vs[0]
	maxF := math.Inf(-1)
	for _, v := range vs {
		fv := f(v)
		if fv > maxF {
			maxF = fv
			maxV = v
		}
	}
	return maxV, maxF
}
