package geometry

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"screwSort/utility"
)

type Segment struct {
	p, q Point
}

func (s Segment) String() string {
	return fmt.Sprintf("Segment{%s, %s}", s.p, s.q)
}

func SegmentPQ(p, q Point) Segment {
	return Segment{p, q}
}

func (s Segment) P() Point {
	return s.p
}

func (s Segment) Q() Point {
	return s.q
}

func (s Segment) D() Point {
	return s.q.Subtract(s.p)
}

func (s Segment) Dx() float64 {
	return s.q.x - s.p.x
}

func (s Segment) Dy() float64 {
	return s.q.y - s.p.y
}

func (s Segment) Center() Point {
	return s.p.AverageWith(s.q)
}

func (s Segment) Length() float64 {
	return s.p.DistanceTo(s.q)
}

func (s Segment) Scale(f float64) Segment {
	return Segment{s.p.Scale(f), s.q.Scale(f)}
}

func (s Segment) CenterScale(f float64) Segment {
	return s.ScaleAbout(s.Center(), f)
}

func (s Segment) ScaleAbout(p Point, f float64) Segment {
	return Segment{s.p.ScaleAbout(p, f), s.q.ScaleAbout(p, f)}
}

func (s Segment) Translate(x, y float64) Segment {
	return Segment{s.p.Translate(x, y), s.q.Translate(x, y)}
}

func (s Segment) Rotate(a float64) Segment {
	return Segment{s.p.Rotate(a), s.q.Rotate(a)}
}

func (s Segment) CenterRotate(a float64) Segment {
	return s.RotateAbout(s.Center(), a)
}

func (s Segment) RotateAbout(p Point, a float64) Segment {
	return Segment{s.p.RotateAbout(p, a), s.q.RotateAbout(p, a)}
}

func (s Segment) PerpDistanceTo(p Point) float64 {
	d := s.D()
	return math.Abs(d.Cross(p)+s.p.Cross(s.q)) / d.R()
}

// SideOf returns positive if p is above s, negative if under, and 0 if collinear
func (s Segment) SideOf(p Point) float64 {
	return s.Dx() * p.OrientationOf(s.p, s.q)
}

// IntersectionOf returns positive if s and t intersect, negative if they do not, and 0 if they touch
func (s Segment) IntersectionOf(t Segment) float64 {
	a := s.p.OrientationOf(t.p, t.q) * s.q.OrientationOf(t.p, t.q)
	b := t.p.OrientationOf(s.p, s.q) * t.q.OrientationOf(s.p, s.q)
	switch {
	case a < 0 && b < 0:
		return 1
	case a < 0 || b < 0:
		return -1
	default:
		return 0
	}
}

func (s Segment) IntersectionWith(t Segment) Point {
	sc, tc, dc := s.p.Cross(s.q), t.p.Cross(t.q), s.D().Cross(t.D())
	return PointXY((sc*t.Dx()-tc*s.Dx())/dc, (sc*t.Dy()-tc*s.Dy())/dc)
}

func (s Segment) AngleBetween(t Segment) float64 {
	return s.D().AngleBetween(t.D())
}

func (s Segment) Draw(im *image.RGBA, c color.RGBA) {
	dx, dy := s.Dx(), s.Dy()
	m := dy / dx
	x0, y0 := s.p.x-0.5, s.p.y-0.5

	if math.Abs(m) < 1 {
		n := utility.AbsInt(utility.IntRound(s.q.x-0.5) - utility.IntRound(x0))
		for i := 0; i <= n; i++ {
			x := math.Round(x0) + math.Copysign(float64(i), dx)
			im.Set(
				utility.IntRound(x),
				utility.IntRound(y0+(x-x0)*m),
				c,
			)
		}
	} else {
		n := utility.AbsInt(utility.IntRound(s.q.y-0.5) - utility.IntRound(y0))
		for i := 0; i <= n; i++ {
			y := math.Round(y0) + math.Copysign(float64(i), dy)
			im.Set(
				utility.IntRound(x0+(y-y0)/m),
				utility.IntRound(y),
				c,
			)
		}
	}
}
