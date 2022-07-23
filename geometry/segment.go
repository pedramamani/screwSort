package geometry

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

type Segment struct {
	p, q, d *Point
}

func SegmentPQ(p, q *Point) *Segment {
	return &Segment{p, q, q.Copy().Subtract(p)}
}

func SegmentXpYpXqYq(xp, yp, xq, yq float64) *Segment {
	p := PointXY(xp, yp)
	q := PointXY(xq, yq)
	return SegmentPQ(p, q)
}

func (s *Segment) P() *Point {
	return s.p
}

func (s *Segment) Q() *Point {
	return s.q
}

func (s *Segment) D() *Point {
	return s.d
}

func (s Segment) Copy() *Segment {
	return SegmentPQ(s.p.Copy(), s.q.Copy())
}

func (s Segment) String() string {
	return fmt.Sprintf("SegmentXpYpXqYq(%.2f, %.2f, %.2f, %.2f)", s.p.x, s.p.y, s.q.x, s.q.y)
}

func (s *Segment) Center() *Point {
	return s.p.Add(s.q).Scale(0.5)
}

func (s *Segment) Length() float64 {
	return s.d.R()
}

func (s *Segment) Scale(f float64) *Segment {
	return s.ScaleAbout(PointXY(0, 0), f)
}

func (s *Segment) CenterScale(f float64) *Segment {
	return s.ScaleAbout(s.Center(), f)
}

func (s *Segment) PScale(f float64) *Segment {
	return s.ScaleAbout(s.p, f)
}

func (s *Segment) QScale(f float64) *Segment {
	return s.ScaleAbout(s.q, f)
}

func (s *Segment) ScaleAbout(p *Point, f float64) *Segment {
	s.p.ScaleAbout(p, f)
	s.q.ScaleAbout(p, f)
	s.d.ScaleAbout(p, f)
	return s
}

func (s *Segment) Translate(x, y float64) *Segment {
	s.p.Translate(x, y)
	s.q.Translate(x, y)
	return s
}

func (s *Segment) Rotate(a float64) *Segment {
	return s.RotateAbout(PointXY(0, 0), a)
}

func (s *Segment) CenterRotate(a float64) *Segment {
	return s.RotateAbout(s.Center(), a)
}

func (s *Segment) PRotate(a float64) *Segment {
	return s.RotateAbout(s.p, a)
}

func (s *Segment) QRotate(a float64) *Segment {
	return s.RotateAbout(s.q, a)
}

func (s *Segment) RotateAbout(p *Point, a float64) *Segment {
	s.p.RotateAbout(p, a)
	s.q.RotateAbout(p, a)
	s.d.RotateAbout(p, a)
	return s
}

func (s *Segment) PerpDistanceTo(p *Point) float64 {
	return math.Abs(s.d.Cross(p)+s.p.Cross(s.q)) / s.d.R()
}

// SideOf returns positive if p is above s, negative if under, and 0 if collinear
func (s *Segment) SideOf(p *Point) float64 {
	return s.d.x * p.OrientationOf(s.p, s.q)
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

func (s *Segment) IntersectionWith(t *Segment) *Point {
	sc, tc, dc := s.p.Cross(s.q), t.p.Cross(t.q), s.d.Cross(t.d)
	return PointXY((sc*t.d.x-tc*s.d.x)/dc, (sc*t.d.y-tc*s.d.y)/dc)
}

func (s *Segment) AngleBetween(t *Segment) float64 {
	return s.d.AngleBetween(t.d)
}

func (s *Segment) Draw(im *image.RGBA, c color.RGBA) *Segment {
	m := s.d.y / s.d.x

	if math.Abs(m) < 1 {
		for x := s.p.x - 0.5; x*s.d.x < (s.q.x-0.5)*s.d.x; x += math.Copysign(1, s.d.x) {
			im.Set(
				im.Rect.Min.X+int(math.Round(x)),
				im.Rect.Max.Y-int(math.Round(s.p.y+0.5+(x-s.p.x)*m)),
				c,
			)
		}
	} else {
		for y := s.p.y + 0.5; y*s.d.y < (s.q.y+0.5)*s.d.y; y += math.Copysign(1, s.d.y) {
			im.Set(
				im.Rect.Min.X+int(math.Round(s.p.x-0.5+(y-s.p.y)/m)),
				im.Rect.Max.Y-int(math.Round(y)),
				c,
			)
		}
	}
	return s
}
