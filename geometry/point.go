package geometry

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"screwSort/utility"
)

// Point describes the 2D point (x, y) with the origin at the top-left of the image
// with x pointing right and y pointing down
type Point struct {
	x, y float64
}

// Vector is a type alias for Point
type Vector = Point

// String returns a string representation of the Point
func (p Point) String() string {
	return fmt.Sprintf("Point{%.2f, %.2f}", p.x, p.y)
}

// PointXY constructs a Point from its x and y coordinates
func PointXY(x, y float64) Point {
	return Point{x, y}
}

// PointImage constructs a Point from the x and y coordinates of a pixel in an image
func PointImage(x, y int) Point {
	return Point{float64(x) + 0.5, float64(y) + 0.5}
}

// X returns the x-coordinate of the Point
func (p Point) X() float64 {
	return p.x
}

// Y returns the y-coordinate of the Point
func (p Point) Y() float64 {
	return p.y
}

// R returns the length of the Vector
func (p Point) R() float64 {
	return math.Hypot(p.x, p.y)
}

// Theta returns the angle of the Vector clockwise from +x
func (p Point) Theta() float64 {
	return math.Atan2(p.y, p.x)
}

// Translate returns a new translated Point
func (p Point) Translate(x, y float64) Point {
	return Point{p.x + x, p.y + y}
}

// Scale returns a new Point scaled by the factor
func (p Point) Scale(f float64) Point {
	return Point{f * p.x, f * p.y}
}

// Rotate returns a new Point rotated clockwise by the angle
func (p Point) Rotate(a float64) Point {
	s, c := math.Sincos(a)
	return Point{c*p.x - s*p.y, s*p.x + c*p.y}
}

// ScaleAbout returns a new Point scaled about another Point by the factor
func (p Point) ScaleAbout(q Point, f float64) Point {
	return p.Subtract(q).Scale(f).Add(q)
}

// RotateAbout returns a new Point rotated clockwise about another Point by the angle
func (p Point) RotateAbout(q Point, a float64) Point {
	return p.Subtract(q).Rotate(a).Add(q)
}

// Add returns a new Vector resulting from addition of another Vector to the Vector
func (p Point) Add(q Point) Point {
	return Point{p.x + q.x, p.y + q.y}
}

// Subtract returns a new Vector resulting from subtraction of another Vector from the Vector
func (p Point) Subtract(q Point) Point {
	return Point{p.x - q.x, p.y - q.y}
}

// AverageWith returns a new Vector resulting from averaging of the Vector with another Vector
func (p Point) AverageWith(q Point) Point {
	return Point{(p.x + q.x) / 2, (p.y + q.y) / 2}
}

// Dot returns the dot product of the Vector with another Vector
func (p Point) Dot(q Point) float64 {
	return p.x*q.x + p.y*q.y
}

// Cross returns the cross product of the Vector with another Vector
func (p Point) Cross(q Point) float64 {
	return p.x*q.y - p.y*q.x
}

// DistanceTo returns the distance between the Point and another Point
func (p Point) DistanceTo(q Point) float64 {
	return math.Hypot(q.x-p.x, q.y-p.y)
}

// AngleBetween returns the angle between the Vector and another Vector
func (p Point) AngleBetween(q Point) float64 {
	return math.Acos(p.Dot(q) / (p.R() * q.R()))
}

// AngleTo returns the angle between +x and the line connecting the Point with another Point
func (p Point) AngleTo(q Point) float64 {
	return math.Atan2(q.y-p.y, q.x-p.x)
}

// OrientationOf returns positive if pab is clockwise, negative if it is counterclockwise, and 0 if it is collinear
func (p Point) OrientationOf(a, b Point) float64 {
	return p.Cross(a) + a.Cross(b) + b.Cross(p)
}

// ToImage returns the x and y coordinates of the pixel corresponding to the Point
func (p Point) ToImage() (int, int) {
	return utility.IntRound(p.x - 0.5), utility.IntRound(p.y - 0.5)
}

// Draw paints the Point on the image with the given color
func (p Point) Draw(im *image.RGBA, c color.RGBA) {
	x, y := p.ToImage()
	im.Set(x, y, c)
}
