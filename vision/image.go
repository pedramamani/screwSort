package vision

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	Black   = color.RGBA{A: 255}
	White   = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	Red     = color.RGBA{R: 255, A: 255}
	Green   = color.RGBA{G: 255, A: 255}
	Blue    = color.RGBA{B: 255, A: 255}
	Yellow  = color.RGBA{R: 255, G: 255, A: 255}
	Magenta = color.RGBA{R: 255, B: 255, A: 255}
	Cyan    = color.RGBA{G: 255, B: 255, A: 255}
)

func BlackRgba(dx, dy int) *image.RGBA {
	out := image.NewRGBA(
		image.Rectangle{
			Min: image.Point{},
			Max: image.Point{X: dx, Y: dy},
		})
	for y := out.Rect.Min.Y; y < out.Rect.Max.Y; y++ {
		for x := out.Rect.Min.X; x < out.Rect.Max.X; x++ {
			out.Set(x, y, Black)
		}
	}
	return out
}

func BlackGray(dx, dy int) *image.Gray {
	out := image.NewGray(
		image.Rectangle{
			Min: image.Point{},
			Max: image.Point{X: dx, Y: dy},
		})
	return out
}

func SavePng(im image.Image, n string) {
	f, _ := os.Create(n)
	_ = png.Encode(f, im)
	_ = f.Close()
}

func OpenPng(fn string) image.Image {
	f, _ := os.Open(fn)
	im, _ := png.Decode(f)
	_ = f.Close()
	return im
}

func ToGray(im image.Image) *image.Gray {
	out := image.NewGray(im.Bounds())
	for y := 0; y < im.Bounds().Dy(); y++ {
		for x := 0; x < im.Bounds().Dx(); x++ {
			out.Set(x, y, im.At(x, y))
		}
	}
	return out
}

func ToRgba(im image.Image) *image.RGBA {
	out := image.NewRGBA(im.Bounds())
	for y := 0; y < im.Bounds().Dy(); y++ {
		for x := 0; x < im.Bounds().Dx(); x++ {
			out.Set(x, y, im.At(x, y))
		}
	}
	return out
}
