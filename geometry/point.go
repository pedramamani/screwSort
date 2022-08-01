package geometry

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"screwSort/utility"
)

type Point struct {
	x, y float64
}

func (p Point) String() string {
	return fmt.Sprintf("Point{%.2f, %.2f}", p.x, p.y)
}

func PointXY(x, y float64) Point {
	return Point{x, y}
}

func PointImage(x, y int) Point {
	return Point{float64(x) + 0.5, float64(y) + 0.5}
}

func (p Point) X() float64 {
	return p.x
}

func (p Point) Y() float64 {
	return p.y
}

func (p Point) R() float64 {
	return math.Hypot(p.x, p.y)
}

func (p Point) Theta() float64 {
	return math.Atan2(p.y, p.x)
}

func (p Point) Translate(x, y float64) Point {
	return Point{p.x + x, p.y + y}
}

func (p Point) Scale(f float64) Point {
	return Point{f * p.x, f * p.y}
}

func (p Point) Rotate(a float64) Point {
	s, c := math.Sincos(a)
	return Point{c*p.x - s*p.y, s*p.x + c*p.y}
}

func (p Point) ScaleAbout(q Point, f float64) Point {
	return p.Subtract(q).Scale(f).Add(q)
}

func (p Point) RotateAbout(q Point, a float64) Point {
	return p.Subtract(q).Rotate(a).Add(q)
}

func (p Point) Add(q Point) Point {
	return Point{p.x + q.x, p.y + q.y}
}

func (p Point) Subtract(q Point) Point {
	return Point{p.x - q.x, p.y - q.y}
}

func (p Point) AverageWith(q Point) Point {
	return Point{(p.x + q.x) / 2, (p.y + q.y) / 2}
}

func (p Point) Dot(q Point) float64 {
	return p.x*q.x + p.y*q.y
}

func (p Point) Cross(q Point) float64 {
	return p.x*q.y - p.y*q.x
}

func (p Point) DistanceTo(q Point) float64 {
	return math.Hypot(q.x-p.x, q.y-p.y)
}

func (p Point) AngleBetween(q Point) float64 {
	return math.Acos(p.Dot(q) / (p.R() * q.R()))
}

func (p Point) AngleTo(q Point) float64 {
	return math.Atan2(q.y-p.y, q.x-p.x)
}

// OrientationOf returns positive if pab is counterclockwise, negative if clockwise, and 0 if collinear
func (p Point) OrientationOf(a, b Point) float64 {
	return p.Cross(a) + a.Cross(b) + b.Cross(p)
}

func (p Point) ToImage() (int, int) {
	return utility.IntRound(p.x - 0.5), utility.IntRound(p.y - 0.5)
}

func (p Point) Draw(im *image.RGBA, c color.RGBA) {
	x, y := p.ToImage()
	im.Set(x, y, c)
}
