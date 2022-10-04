package geometry

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

// Line describes the 2D line ax+by+c=0 with the origin at the top-left of the image
// with x pointing right and y pointing down
//
// It is assumed that a²+b²=1 is always satisfied
type Line struct {
	a, b, c float64
}

// String returns a string representation of the Line in slope-intercept form
func (l Line) String() string {
	return fmt.Sprintf("LineMY0(%.2f, %.2f)", l.M(), l.Y0())
}

// LineMY0 constructs a Line given the slope and y-intercept y=mx+y₀
func LineMY0(m, y0 float64) Line {
	d := math.Hypot(1, m)
	return Line{-m / d, 1 / d, -y0 / d}
}

// LineAC constructs a Line given the a and c coefficients in ax+by+c=0 assuming a²+b²=1
//
// It panics if |a| > 1
func LineAC(a, c float64) Line {
	if math.Abs(a) > 1 {
		panic("failed to satisfy |a| <= 1")
	}
	return Line{a, math.Sqrt(1 - a*a), c}
}

// LineABC constructs a Line given the three coefficients in ax+by+c=0
//
// It normalizes a, b, and c to satisfy a²+b²=1
func LineABC(a, b, c float64) Line {
	d := math.Hypot(a, b)
	if b < 0 {
		a, b, c = -a, -b, -c
	}
	return Line{a / d, b / d, c / d}
}

// LinePTheta constructs a Line given a Point and an angle clockwise from +x
func LinePTheta(p Point, theta float64) Line {
	m := math.Tan(theta)
	y0 := p.y - m*p.x
	return LineMY0(m, y0)
}

// A returns the normalized a coefficient in ax+by+c=0
func (l Line) A() float64 {
	return l.a
}

// B returns the normalized b coefficient in ax+by+c=0
func (l Line) B() float64 {
	return l.b
}

// C returns the normalized c coefficient in ax+by+c=0
func (l Line) C() float64 {
	return l.c
}

// M returns the slope of the Line
func (l Line) M() float64 {
	return -l.a / l.b
}

// Y0 returns the y-intercept of the Line
func (l Line) Y0() float64 {
	return -l.c / l.b
}

// X returns the x-coordinate of a point on the Line given its y-coordinate
func (l Line) X(y float64) float64 {
	return -(l.b*y + l.c) / l.a
}

// Y returns the y-coordinate of a point on the Line given its x-coordinate
func (l Line) Y(x float64) float64 {
	return -(l.a*x + l.c) / l.b
}

// Scale returns a new Line scaled by the factor
func (l Line) Scale(f float64) Line {
	return Line{l.a, l.b, f * l.c}
}

// Translate returns a new translated Line
func (l Line) Translate(x, y float64) Line {
	return Line{l.a, l.b, l.c - l.a*x - l.b*y}
}

// PerpDistanceTo returns the perpendicular distance of the Point from the Line
func (l Line) PerpDistanceTo(p Point) float64 {
	return math.Abs(l.a*p.x + l.b*p.y + l.c)
}

// SideOf returns positive if the Point is under the Line, negative if it is above, and 0 if it is on
func (l Line) SideOf(p Point) float64 {
	return l.a*p.x + l.b*p.y + l.c
}

// AngleBetween returns the angle between the Line and the other Line
//
// The returned value is between 0 and Pi (inclusive)
func (l Line) AngleBetween(o Line) float64 {
	return math.Abs(math.Atan(l.M()) - math.Atan(o.M()))
}

// IntersectionWith returns the intersection Point of the Line with the other Line
//
// It panics if the two are parallel
func (l Line) IntersectionWith(o Line) Point {
	d := l.a*o.b - l.b*o.a
	if d == 0 {
		panic("parallel lines do not intersect")
	}
	return PointXY((l.b*o.c-l.c*o.b)/d, (l.c*o.a-l.a*o.c)/d)
}

// ToSegment returns the Segment of the Line that is inside the rectangular window
// represented by its top-left Point and bottom-right Point
//
// It panics if the Line does not overlap with the window
func (l Line) ToSegment(pTL, pBR Point) Segment {
	pTR, pBL := PointXY(pBR.x, pTL.y), PointXY(pTL.x, pBR.y)
	var ps []Point

	if l.SideOf(pTL)*l.SideOf(pTR) <= 0 {
		ps = append(ps, PointXY(l.X(pTL.y), pTL.y))
	}
	if l.SideOf(pBL)*l.SideOf(pBR) <= 0 {
		ps = append(ps, PointXY(l.X(pBL.y), pBL.y))
	}
	if l.SideOf(pTL)*l.SideOf(pBL) < 0 {
		ps = append(ps, PointXY(pTL.x, l.Y(pTL.x)))
	}
	if l.SideOf(pTR)*l.SideOf(pBR) < 0 {
		ps = append(ps, PointXY(pTR.x, l.Y(pTR.x)))
	}

	if len(ps) != 2 {
		panic("line does not pass through the window")
	}
	return SegmentPQ(ps[0], ps[1])
}

// Draw paints the Line on the image with the given color
func (l Line) Draw(im *image.RGBA, c color.RGBA) {
	b := im.Bounds()
	l.ToSegment(
		PointXY(float64(b.Min.X), float64(b.Min.Y)),
		PointXY(float64(b.Max.X), float64(b.Max.Y)),
	).Draw(im, c)
}
