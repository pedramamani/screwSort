package fit

import (
	"math"
	"screwSort/geometry"
	"screwSort/utility"
)

// Fit describes a 2D fit line along with errors associated with each point
type Fit struct {
	line   geometry.Line
	errors []float64
}

// FitLineErrors constructs a new Fit from the line and the errors
func FitLineErrors(line geometry.Line, errors []float64) Fit {
	return Fit{line, errors}

}

// Line returns the Line of the Fit
func (f Fit) Line() geometry.Line {
	return f.line
}

// Errors returns the errors of the Fit
func (f Fit) Errors() []float64 {
	return f.errors
}

// MeanError returns the mean error of the Fit
func (f Fit) MeanError() float64 {
	return utility.Mean(f.errors...)
}

// MaxError returns the maximum error of the Fit
func (f Fit) MaxError() float64 {
	return utility.Max(f.errors...)
}

// Moments returns the central moments of the points up to order 2
func Moments(ps []geometry.Point) (xBar, yBar, xyBar, x2Bar, y2Bar float64) {
	for _, p := range ps {
		xBar += p.X()
		yBar += p.Y()
		xyBar += p.X() * p.Y()
		x2Bar += p.X() * p.X()
		y2Bar += p.Y() * p.Y()
	}
	n := float64(len(ps))
	xBar /= n
	yBar /= n
	xyBar /= n
	x2Bar /= n
	y2Bar /= n
	return
}

// LinearFit returns a linear Fit that minimizes the y residual of the points
func LinearFit(ps []geometry.Point) Fit {
	xBar, yBar, xyBar, x2Bar, _ := Moments(ps)
	m := (xyBar - xBar*yBar) / (x2Bar - xBar*xBar)
	y0 := yBar - m*xBar
	l := geometry.LineMY0(m, y0)

	es := make([]float64, len(ps))
	for i, p := range ps {
		es[i] = math.Abs(p.Y() - l.Y(p.X()))
	}
	return FitLineErrors(l, es)
}

// OrthogonalFit returns an orthogonal Fit that minimizes the perpendicular distance to the points
func OrthogonalFit(ps []geometry.Point) Fit {
	xBar, yBar, xyBar, x2Bar, y2Bar := Moments(ps)
	covXY := xyBar - xBar*yBar
	vX, vY := x2Bar-xBar*xBar, y2Bar-yBar*yBar

	var a float64
	if covXY != 0 {
		f := (vX - vY) / covXY
		f2 := f * f
		s := -covXY / math.Abs(covXY)
		a = s * math.Sqrt(f2/(f2+4)+s*f/math.Sqrt(f2+4)+4/(f2+4)) / math.Sqrt2
	} else if vX < vY {
		a = 1
	} else {
		a = 0
	}
	l := geometry.LineAC(a, -a*xBar-math.Sqrt(1-a*a)*yBar)

	es := make([]float64, len(ps))
	for i, p := range ps {
		es[i] = l.PerpDistanceTo(p)
	}
	return FitLineErrors(l, es)
}
