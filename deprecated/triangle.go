package deprecated

import (
	"gometry/geometry"
	"image"
	"image/color"
	"math"
)

type Triangle struct {
	a, b, c  *geometry.Point
	min, max *geometry.Point
	fill     color.RGBA
}

func NewTriangle(a, b, c *geometry.Point, fill color.RGBA) *Triangle {
	o := a.OrientationOf(b, c)
	if o == 0 {
		panic("the vertices of a triangle cannot be collinear")
	} else if o < 0 {
		b, c = c, b // order points counter-clockwise
	}
	minX, maxX := math.Inf(1), math.Inf(-1)
	minY, maxY := math.Inf(1), math.Inf(-1)
	for _, p := range []*geometry.Point{a, b, c} {
		if p.x() < minX {
			minX = p.x()
		}
		if p.x() > maxX {
			maxX = p.x()
		}
		if p.y() < minY {
			minY = p.y()
		}
		if p.y() > maxY {
			maxY = p.y()
		}
	}
	return &Triangle{a, b, c, geometry.NewPoint(minX, minY), geometry.NewPoint(maxX, maxY), fill}
}

// RegionOf returns positive if p is inside triangle, negative if outside, and 0 if on
func (t *Triangle) RegionOf(p *geometry.Point) float64 {
	a, b, c := p.OrientationOf(t.a, t.b), p.OrientationOf(t.b, t.c), p.OrientationOf(t.c, t.a)
	switch {
	case a > 0 && b > 0 && c > 0:
		return 1
	case a < 0 || b < 0 || c < 0:
		return -1
	default:
		return 0
	}
}

func (t *Triangle) Draw(im *image.RGBA) *image.RGBA {
	_, height := widthHeight(im)

	for y := math.Ceil(t.min.y() - 0.5); y <= math.Floor(t.max.y()-0.5); y++ {
		for x := math.Ceil(t.min.x() - 0.5); x <= math.Floor(t.max.x()-0.5); x++ {
			p := geometry.NewPoint(x, y)
			if t.RegionOf(p) >= 0 {
				im.SetRGBA(int(x), int(height-y), t.fill)
			}
		}
	}
	return im
}
