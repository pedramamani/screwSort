package geometry

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

type Point struct {
	x, y float64
}

func PointXY(x, y float64) *Point {
	return &Point{x, y}
}

func PointRTheta(r, theta float64) *Point {
	s, c := math.Sincos(theta)
	return &Point{r * c, r * s}
}

func PointImage(x, y int, im image.Image) *Point {
	return PointXY(
		float64(x-im.Bounds().Min.X)+0.5,
		float64(im.Bounds().Max.Y-y)-0.5,
	)
}

func (p *Point) X() float64 {
	return p.x
}

func (p *Point) Y() float64 {
	return p.y
}

func (p *Point) R() float64 {
	return math.Hypot(p.x, p.y)
}

func (p *Point) Theta() float64 {
	return math.Atan2(p.y, p.x)
}

func (p *Point) Copy() *Point {
	return PointXY(p.x, p.y)
}

func (p *Point) String() string {
	return fmt.Sprintf("PointXY(%.2f, %.2f)", p.x, p.y)
}

func (p *Point) Translate(x, y float64) *Point {
	p.x += x
	p.y += y
	return p
}

func (p *Point) Scale(f float64) *Point {
	p.x *= f
	p.y *= f
	return p
}

func (p *Point) Rotate(a float64) *Point {
	s, c := math.Sincos(a)
	p.x, p.y = c*p.x-s*p.y, s*p.x+c*p.y
	return p
}

func (p *Point) ScaleAbout(q *Point, f float64) *Point {
	return p.Subtract(q).Scale(f).Add(q)
}

func (p *Point) RotateAbout(q *Point, a float64) *Point {
	return p.Subtract(q).Rotate(a).Add(q)
}

func (p *Point) Add(q *Point) *Point {
	return p.Translate(q.x, q.y)
}

func (p *Point) Subtract(q *Point) *Point {
	return p.Translate(-q.x, -q.y)
}

func (p *Point) Dot(q *Point) float64 {
	return p.x*q.x + p.y*q.y
}

func (p *Point) Cross(q *Point) float64 {
	return p.x*q.y - p.y*q.x
}

func (p *Point) DistanceTo(q *Point) float64 {
	return math.Hypot(q.x-p.x, q.y-p.y)
}

func (p *Point) AngleBetween(q *Point) float64 {
	return math.Acos(p.Dot(q) / (p.R() * q.R()))
}

func (p *Point) AngleTo(q *Point) float64 {
	return math.Atan2(q.y-p.y, q.x-p.x)
}

// OrientationOf returns positive if pab is counterclockwise, negative if clockwise, and 0 if collinear
func (p *Point) OrientationOf(a, b *Point) float64 {
	return p.Cross(a) + a.Cross(b) + b.Cross(p)
}

func (p *Point) ToImage(im image.Image) (int, int) {
	return im.Bounds().Min.X + int(math.Round(p.x-0.5)), im.Bounds().Max.Y - int(math.Round(p.y+0.5))
}

func (p *Point) Draw(im *image.RGBA, c color.RGBA) *Point {
	x, y := p.ToImage(im)
	im.Set(x, y, c)
	return p
}
