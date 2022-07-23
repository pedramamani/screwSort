package deprecated

import (
	"image"
	"image/color"
	"screwSort/geometry"
)

type Polygon struct {
	points    []*geometry.Point
	triangles []*Triangle
}

func NewPolygon(ps []*geometry.Point, fill color.RGBA) *Polygon {
	n := len(ps)
	if n <= 2 {
		panic("at least 3 points are required for a polygon")
	}

	ts := make([]*Triangle, n-2)
	for i := 0; i < n-2; i++ { // todo: fix logic for convex polygons
		ts[i] = NewTriangle(ps[0], ps[i+1], ps[i+2], fill)
	}
	return &Polygon{ps, ts}
}

func (p *Polygon) Points() []*geometry.Point {
	return p.points
}

func (p *Polygon) Triangles() []*Triangle {
	return p.triangles
}

func (p *Polygon) ApplyToPoints(f func(*geometry.Point) *geometry.Point) *Polygon {
	newPoints := make([]*geometry.Point, len(p.points))
	for _, p := range p.points {
		newPoints = append(newPoints, f(p))
	}
	return NewPolygon(newPoints, p.triangles[0].fill)
}

func (p *Polygon) Scale(a float64) *Polygon {
	return p.ScaleAbout(geometry.Origin, a)
}

func (p *Polygon) ScaleAbout(q *geometry.Point, a float64) *Polygon {
	return p.ApplyToPoints(func(x *geometry.Point) *geometry.Point { return x.ScaleAbout(q, a) })
}

func (p *Polygon) Translate(dx, dy float64) *Polygon {
	return p.ApplyToPoints(func(x *geometry.Point) *geometry.Point { return x.Translate(dx, dy) })
}

func (p *Polygon) Rotate(angle float64) *Polygon {
	return p.RotateAbout(geometry.Origin, angle)
}

func (p *Polygon) RotateAbout(q *geometry.Point, angle float64) *Polygon {
	return p.ApplyToPoints(func(x *geometry.Point) *geometry.Point { return x.RotateAbout(q, angle) })
}

// RegionOf returns positive if p is inside polygon, negative if outside, and 0 if on
func (p *Polygon) RegionOf(q *geometry.Point) float64 {
	for _, t := range p.triangles {
		if t.RegionOf(q) < 0 {
			return -1
		} else if t.RegionOf(q) == 0 {
			return 0
		}
	}
	return +1
}

func (p *Polygon) Draw(im *image.RGBA) *image.RGBA {
	for _, t := range p.triangles {
		t.Draw(im)
	}
	return im
}
