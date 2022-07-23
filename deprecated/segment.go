package deprecated

import (
	"image"
	"image/color"
	"math"
	"screwSort/geometry"
)

type Segment struct {
	p, q, center, delta *geometry.Point
	triangles           []*Triangle
	thickness, length   float64
}

func NewSegment(p, q *geometry.Point, thickness float64, fill color.RGBA) *Segment {
	if p.DistanceTo(q) == 0 {
		panic("segment must have positive length")
	}
	delta := q.Subtract(p)
	posOffset := geometry.NewPoint(0, thickness/2).Rotate(delta.Angle())
	negOffset := geometry.NewPoint(0, -thickness/2).Rotate(delta.Angle())
	ts := []*Triangle{
		NewTriangle(p.Add(negOffset), q.Add(negOffset), q.Add(posOffset), fill),
		NewTriangle(q.Add(posOffset), p.Add(negOffset), p.Add(posOffset), fill),
	}
	return &Segment{p, q, p.AverageWith(q), delta, ts, thickness, p.DistanceTo(q)}
}

func (s *Segment) P() *geometry.Point {
	return s.p
}

func (s *Segment) Q() *geometry.Point {
	return s.q
}

func (s *Segment) Center() *geometry.Point {
	return s.center
}

func (s *Segment) Delta() *geometry.Point {
	return s.delta
}

func (s *Segment) Length() float64 {
	return s.length
}

func (s *Segment) Scale(a float64) *Segment {
	return s.ScaleAbout(geometry.Origin, a)
}

func (s *Segment) CenterScale(a float64) *Segment {
	return s.ScaleAbout(s.center, a)
}

func (s *Segment) PScale(a float64) *Segment {
	return s.ScaleAbout(s.p, a)
}

func (s *Segment) QScale(a float64) *Segment {
	return s.ScaleAbout(s.q, a)
}

func (s *Segment) ScaleAbout(p *geometry.Point, a float64) *Segment {
	if a == 0 {
		panic("scaling segment by 0 leads to zero length")
	}
	return NewSegment(s.p.ScaleAbout(p, a), s.q.ScaleAbout(p, a), s.thickness, s.triangles[0].fill)
}

func (s *Segment) Translate(dx, dy float64) *Segment {
	return NewSegment(s.p.Translate(dx, dy), s.q.Translate(dx, dy), s.thickness, s.triangles[0].fill)
}

func (s *Segment) Rotate(angle float64) *Segment {
	return s.RotateAbout(geometry.Origin, angle)
}

func (s *Segment) CenterRotate(angle float64) *Segment {
	return s.RotateAbout(s.center, angle)
}

func (s *Segment) PRotate(angle float64) *Segment {
	return s.RotateAbout(s.p, angle)
}

func (s *Segment) QRotate(angle float64) *Segment {
	return s.RotateAbout(s.q, angle)
}

func (s *Segment) RotateAbout(p *geometry.Point, angle float64) *Segment {
	return NewSegment(s.p.RotateAbout(p, angle), s.q.RotateAbout(p, angle), s.thickness, s.triangles[0].fill)
}

func (s *Segment) PerpDistanceTo(p *geometry.Point) float64 {
	return math.Abs(s.delta.Cross(p)+s.p.Cross(s.q)) / s.length
}

// SideOf returns positive if p is above s, negative if under, and 0 if collinear
func (s *Segment) SideOf(p *geometry.Point) float64 {
	return s.delta.x() * p.OrientationOf(s.p, s.q)
}

func (s *Segment) AngleBetween(t *Segment) float64 {
	return s.delta.AngleBetween(t.delta)
}

// IntersectionOf returns positive if s and t intersect, negative if they do not, and 0 if they touch
func (s *Segment) IntersectionOf(t *Segment) float64 {
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

func (s *Segment) IntersectionWith(t *Segment) *geometry.Point {
	sc, tc, dc := s.p.Cross(s.q), t.p.Cross(t.q), s.delta.Cross(t.delta)
	return geometry.NewPoint((sc*t.delta.x()-tc*s.delta.x())/dc, (sc*t.delta.y()-tc*s.delta.y())/dc)
}

func (s *Segment) Draw(im *image.RGBA) *image.RGBA {
	for _, t := range s.triangles {
		t.Draw(im)
	}
	return im
}
