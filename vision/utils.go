package vision

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

var Red = color.RGBA{R: 255, A: 255}
var Green = color.RGBA{G: 255, A: 255}
var Blue = color.RGBA{B: 255, A: 255}
var Yellow = color.RGBA{R: 255, G: 255, A: 255}
var Magenta = color.RGBA{R: 255, B: 255, A: 255}
var Cyan = color.RGBA{G: 255, B: 255, A: 255}
var White = color.RGBA{R: 255, G: 255, B: 255, A: 255}
var Black = color.RGBA{A: 255}

func distance(a image.Point, b image.Point) float64 {
	return math.Hypot(float64(b.X-a.X), float64(b.Y-a.Y))
}

func minInt(a int, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

func maxInt(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func BlankRgba(w, h int, c color.RGBA) *image.RGBA {
	im := image.NewRGBA(
		image.Rectangle{
			Min: image.Point{},
			Max: image.Point{X: w, Y: h},
		})
	for y := im.Rect.Min.Y; y < im.Rect.Max.Y; y++ {
		for x := im.Rect.Min.X; x < im.Rect.Max.X; x++ {
			im.Set(x, y, c)
		}
	}
	return im
}

func Save(im image.Image, n string) {
	f, _ := os.Create(n)
	png.Encode(f, im)
	f.Close()
}

func Open(n string) image.Image {
	f, _ := os.Open(n)
	im, _ := png.Decode(f)
	f.Close()
	return im
}
