package deprecated

import (
	"gometry/geometry"
	"image"
	"image/color"
	"math"
)

type Circle struct {
	center   *geometry.Point
	radius   float64
	min, max *geometry.Point
	fill     color.RGBA
}

func NewCircle(c *geometry.Point, r float64, fill color.RGBA) *Circle {
	if r <= 0 {
		panic("circle must have positive radius")
	}
	return &Circle{c, r, c.Translate(-r, -r), c.Translate(r, r), fill}
}

func (c *Circle) Center() *geometry.Point {
	return c.center
}

func (c *Circle) Radius() float64 {
	return c.radius
}

func (c *Circle) Scale(a float64) *Circle {
	return c.ScaleAbout(geometry.Origin, a)
}

func (c *Circle) CenterScale(a float64) *Circle {
	return c.ScaleAbout(c.center, a)
}

func (c *Circle) ScaleAbout(p *geometry.Point, a float64) *Circle {
	if a == 0 {
		panic("scaling circle by 0 leads to zero radius")
	}
	return NewCircle(c.center.ScaleAbout(p, a), math.Abs(a)*c.radius, c.fill)
}

func (c *Circle) Translate(dx, dy float64) *Circle {
	return NewCircle(c.center.Translate(dx, dy), c.radius, c.fill)
}

func (c *Circle) Rotate(angle float64) *Circle {
	return c.RotateAbout(geometry.Origin, angle)
}

func (c *Circle) RotateAbout(p *geometry.Point, angle float64) *Circle {
	return NewCircle(c.center.RotateAbout(p, angle), c.radius, c.fill)
}

// RegionOf returns positive if p is inside circle, negative if outside, and 0 if on
func (c *Circle) RegionOf(p *geometry.Point) float64 {
	return c.radius - c.center.DistanceTo(p)
}

func (c *Circle) Draw(im *image.RGBA) *image.RGBA {
	_, height := widthHeight(im)

	for y := math.Ceil(c.min.y() - 0.5); y <= math.Floor(c.max.y()-0.5); y++ {
		for x := math.Ceil(c.min.x() - 0.5); x <= math.Floor(c.max.x()-0.5); x++ {
			p := geometry.NewPoint(x, y)
			if c.RegionOf(p) >= 0 {
				im.SetRGBA(int(x), int(height-y), c.fill)
			}
		}
	}
	return im
}
