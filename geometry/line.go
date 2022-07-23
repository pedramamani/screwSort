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

func LineMY0(m, y0 float64) *Line {
	d := math.Hypot(1, m)
	return &Line{-m / d, 1 / d, -y0 / d}
}

func LineAC(a, c float64) *Line {
	if math.Abs(a) > 1 {
		panic("a must be between -1 and 1 (inclusive)")
	}
	return &Line{a, math.Sqrt(1 - math.Pow(a, 2)), c}
}

func LineABC(a, b, c float64) *Line {
	d := math.Hypot(a, b)
	if b < 0 {
		a, b, c = -a, -b, -c
	}
	return &Line{a / d, b / d, c / d}
}

func LinePTheta(p *Point, theta float64) *Line {
	m := math.Tan(theta)
	y0 := p.y - m*p.x
	return LineMY0(m, y0)
}

func (l *Line) A() float64 {
	return l.a
}

func (l *Line) B() float64 {
	return l.b
}

func (l *Line) C() float64 {
	return l.c
}

func (l *Line) M() float64 {
	return -l.a / l.b
}

func (l *Line) Y0() float64 {
	return -l.c / l.b
}

func (l Line) Copy() *Line {
	return LineAC(l.a, l.c)
}

func (l *Line) String() string {
	return fmt.Sprintf("LineMY0(%.2f, %.2f)", l.M(), l.Y0())
}

func (l *Line) X(y float64) float64 {
	return -(l.b*y + l.c) / l.a
}

func (l *Line) Y(x float64) float64 {
	return -(l.a*x + l.c) / l.b
}

func (l *Line) Scale(a float64) *Line {
	l.c *= a
	return l
}

func (l *Line) Translate(p *Point) *Line {
	l.c = l.c - l.a*p.x - l.b*p.y
	return l
}

func (l *Line) PerpDistanceTo(p *Point) float64 {
	return math.Abs(l.a*p.x + l.b*p.y + l.c)
}

// SideOf returns positive if p is above l, negative if under, and 0 if on
func (l *Line) SideOf(p *Point) float64 {
	return l.a*p.x + l.b*p.y + l.c
}

func (l *Line) IntersectionWith(o *Line) *Point {
	d := l.a*o.b - l.b*o.a
	return PointXY((l.b*o.c-l.c*o.b)/d, (l.c*o.a-l.a*o.c)/d)
}

func (l Line) ToSegment(pMin, pMax *Point) *Segment {
	pBr, pTl := PointXY(pMin.x, pMax.y), PointXY(pMax.x, pMin.y)
	var ps []*Point

	if l.SideOf(pMin)*l.SideOf(pBr) <= 0 {
		ps = append(ps, PointXY(l.X(pMin.y), pMin.y))
	}
	if l.SideOf(pTl)*l.SideOf(pMax) <= 0 {
		ps = append(ps, PointXY(l.X(pMax.y), pMax.y))
	}
	if l.SideOf(pMin)*l.SideOf(pTl) < 0 {
		ps = append(ps, PointXY(pMin.x, l.Y(pMin.x)))
	}
	if l.SideOf(pBr)*l.SideOf(pMax) < 0 {
		ps = append(ps, PointXY(pMax.x, l.Y(pMax.x)))
	}

	if len(ps) >= 2 {
		return SegmentPQ(ps[0], ps[1])
	}
	return nil
}

func (l *Line) Draw(im *image.RGBA, c color.RGBA) *Line {
	l.ToSegment(
		PointXY(0, 0),
		PointXY(float64(im.Rect.Max.X-im.Rect.Min.X), float64(im.Rect.Max.Y-im.Rect.Min.Y)),
	).Draw(im, c)
	return l
}
