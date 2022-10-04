package geometry

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"screwSort/utility"
)

// Segment describes a 2D line segment by its two end points
type Segment struct {
	p, q Point
}

// String returns a string representation of the Segment
func (s Segment) String() string {
	return fmt.Sprintf("Segment{%s, %s}", s.p, s.q)
}

// SegmentPQ constructs a segment from its end points
func SegmentPQ(p, q Point) Segment {
	return Segment{p, q}
}

// P returns the first end Point of the Segment
func (s Segment) P() Point {
	return s.p
}

// Q returns the second end Point of the Segment
func (s Segment) Q() Point {
	return s.q
}

// D returns the Vector from the first to the second end point of the Segment
func (s Segment) D() Point {
	return s.q.Subtract(s.p)
}

// Dx returns the x length of the Segment
func (s Segment) Dx() float64 {
	return s.q.x - s.p.x
}

// Dy returns the y length of the Segment
func (s Segment) Dy() float64 {
	return s.q.y - s.p.y
}

// Center returns the center Point of the Segment
func (s Segment) Center() Point {
	return s.p.AverageWith(s.q)
}

// Length returns the length of the Segment
func (s Segment) Length() float64 {
	return s.p.DistanceTo(s.q)
}

// Scale returns a new Segment scaled by the factor
func (s Segment) Scale(f float64) Segment {
	return Segment{s.p.Scale(f), s.q.Scale(f)}
}

// CenterScale returns a new Segment scaled about its center by the factor
func (s Segment) CenterScale(f float64) Segment {
	return s.ScaleAbout(s.Center(), f)
}

// ScaleAbout returns a new Segment scaled about the Point by the factor
func (s Segment) ScaleAbout(p Point, f float64) Segment {
	return Segment{s.p.ScaleAbout(p, f), s.q.ScaleAbout(p, f)}
}

// Translate returns a new translated Segment
func (s Segment) Translate(x, y float64) Segment {
	return Segment{s.p.Translate(x, y), s.q.Translate(x, y)}
}

// Rotate returns a new Segment rotated clockwise by the angle
func (s Segment) Rotate(a float64) Segment {
	return Segment{s.p.Rotate(a), s.q.Rotate(a)}
}

// CenterRotate returns a new Segment rotated clockwise about its center by the angle
func (s Segment) CenterRotate(a float64) Segment {
	return s.RotateAbout(s.Center(), a)
}

// RotateAbout returns a new Segment rotated clockwise about the point by the angle
func (s Segment) RotateAbout(p Point, a float64) Segment {
	return Segment{s.p.RotateAbout(p, a), s.q.RotateAbout(p, a)}
}

// PerpDistanceTo returns the perpendicular distance of the Point from the Segment
func (s Segment) PerpDistanceTo(p Point) float64 {
	d := s.D()
	return math.Abs(d.Cross(p)+s.p.Cross(s.q)) / d.R()
}

// SideOf returns positive if the Point is under the Segment, negative if it is above, and 0 if it is on
func (s Segment) SideOf(p Point) float64 {
	return s.Dx() * p.OrientationOf(s.p, s.q)
}

// IntersectionOf returns positive if the Segment intersects with another Segment,
// negative if they do not, and 0 if they touch but do not intersect
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

// IntersectionWith returns the intersection Point of the Segment with the other Segment
func (s Segment) IntersectionWith(t Segment) Point {
	sc, tc, dc := s.p.Cross(s.q), t.p.Cross(t.q), s.D().Cross(t.D())
	return PointXY((sc*t.Dx()-tc*s.Dx())/dc, (sc*t.Dy()-tc*s.Dy())/dc)
}

// AngleBetween returns the angle between the Segment and the other Segment
func (s Segment) AngleBetween(t Segment) float64 {
	return s.D().AngleBetween(t.D())
}

// Draw paints the Segment on the image with the given color
//
// It panics if the Segment is not well-defined
func (s Segment) Draw(im *image.RGBA, c color.RGBA) {
	if math.IsNaN(s.Length()) {
		panic("the segment is not well-defined")
	}

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
