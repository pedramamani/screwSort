package geometry

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

type Line struct {
	a, b, c float64
}

func (l Line) String() string {
	return fmt.Sprintf("LineMY0(%.2f, %.2f)", l.M(), l.Y0())
}

func LineMY0(m, y0 float64) Line {
	d := math.Hypot(1, m)
	return Line{-m / d, 1 / d, -y0 / d}
}

func LineAC(a, c float64) Line {
	if math.Abs(a) > 1 {
		panic("failed to satisfy |a| <= 1")
	}
	return Line{a, math.Sqrt(1 - a*a), c}
}

func LineABC(a, b, c float64) Line {
	d := math.Hypot(a, b)
	if b < 0 {
		a, b, c = -a, -b, -c
	}
	return Line{a / d, b / d, c / d}
}

func LinePTheta(p Point, theta float64) Line {
	m := math.Tan(theta)
	y0 := p.y - m*p.x
	return LineMY0(m, y0)
}

func (l Line) A() float64 {
	return l.a
}

func (l Line) B() float64 {
	return l.b
}

func (l Line) C() float64 {
	return l.c
}

func (l Line) M() float64 {
	return -l.a / l.b
}

func (l Line) Y0() float64 {
	return -l.c / l.b
}

func (l Line) X(y float64) float64 {
	return -(l.b*y + l.c) / l.a
}

func (l Line) Y(x float64) float64 {
	return -(l.a*x + l.c) / l.b
}

func (l Line) Scale(a float64) Line {
	return Line{l.a, l.b, a * l.c}
}

func (l Line) Translate(p Point) Line {
	return Line{l.a, l.b, l.c - l.a*p.x - l.b*p.y}
}

func (l Line) PerpDistanceTo(p Point) float64 {
	return math.Abs(l.a*p.x + l.b*p.y + l.c)
}

// SideOf returns positive if p is above l, negative if under, and 0 if on
func (l Line) SideOf(p Point) float64 {
	return l.a*p.x + l.b*p.y + l.c
}

func (l Line) AngleBetween(o Line) float64 {
	return math.Abs(math.Atan(l.M()) - math.Atan(o.M()))
}

func (l Line) IntersectionWith(o Line) Point {
	d := l.a*o.b - l.b*o.a
	if d == 0 { // todo: handle this case properly
		fmt.Println("parallel lines", l.String(), o.String(), "do not intersect")
	}
	return PointXY((l.b*o.c-l.c*o.b)/d, (l.c*o.a-l.a*o.c)/d)
}

func (l Line) ToSegment(p Point) Segment {
	pTL, pTR, pBR, pBL := PointXY(0, 0), PointXY(p.x, 0), p, PointXY(0, p.y)
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
		panic("failed to satisfy len(ps) == 2")
	}
	return SegmentPQ(ps[0], ps[1])
}

func (l Line) Draw(im *image.RGBA, c color.RGBA) {
	l.ToSegment(PointXY(float64(im.Rect.Dx()), float64(im.Rect.Dy()))).Draw(im, c)
}
