package vision

import (
	"image"
	"image/color"
)

const B uint8 = 0
const W uint8 = 255
const N uint8 = 8

func Threshold(im *image.Gray, t uint8) *image.Gray {
	imOut := image.NewGray(im.Rect)
	for y := im.Rect.Min.Y; y < im.Rect.Max.Y; y++ {
		for x := im.Rect.Min.X; x < im.Rect.Max.X; x++ {
			v := im.GrayAt(x, y).Y
			if v >= t {
				v = W
			} else {
				v = B
			}
			imOut.Set(x, y, color.Gray{Y: v})
		}
	}
	return imOut
}

func InverseThreshold(im *image.Gray, t uint8) *image.Gray {
	imOut := image.NewGray(im.Rect)
	for y := im.Rect.Min.Y; y < im.Rect.Max.Y; y++ {
		for x := im.Rect.Min.X; x < im.Rect.Max.X; x++ {
			v := im.GrayAt(x, y).Y
			if v <= t {
				v = W
			} else {
				v = B
			}
			imOut.Set(x, y, color.Gray{Y: v})
		}
	}
	return imOut
}

func Invert(im *image.Gray) *image.Gray {
	imOut := image.NewGray(im.Rect)
	for y := im.Rect.Min.Y; y < im.Rect.Max.Y; y++ {
		for x := im.Rect.Min.X; x < im.Rect.Max.X; x++ {
			v := im.GrayAt(x, y).Y
			imOut.Set(x, y, color.Gray{Y: W - v})
		}
	}
	return imOut
}

func ToGray(im image.Image) *image.Gray {
	imOut := image.NewGray(im.Bounds())
	for y := im.Bounds().Min.Y; y < im.Bounds().Max.Y; y++ {
		for x := im.Bounds().Min.X; x < im.Bounds().Max.X; x++ {
			imOut.Set(x, y, im.At(x, y))
		}
	}
	return imOut
}

func ToRgba(im image.Image) *image.RGBA {
	imOut := image.NewRGBA(im.Bounds())
	for y := im.Bounds().Min.Y; y < im.Bounds().Max.Y; y++ {
		for x := im.Bounds().Min.X; x < im.Bounds().Max.X; x++ {
			imOut.Set(x, y, im.At(x, y))
		}
	}
	return imOut
}

func UpscaleRgba(im *image.RGBA, f int) *image.RGBA {
	imOut := image.NewRGBA(image.Rectangle{Max: image.Point{X: f * im.Rect.Max.X, Y: f * im.Rect.Max.Y}})
	for y := 0; y < im.Rect.Dy(); y++ {
		for x := 0; x < im.Rect.Dx(); x++ {
			c := im.At(x, y)
			for iy := 0; iy < f; iy++ {
				for ix := 0; ix < f; ix++ {
					imOut.Set(f*x+ix, f*y+iy, c)
				}
			}
		}
	}
	return imOut
}

func CropSquare(im *image.Gray) *image.Gray {
	w := im.Rect.Max.X - im.Rect.Min.X
	h := im.Rect.Max.Y - im.Rect.Min.Y
	s := 2 * (minInt(w, h)/2 - 1)

	imOut := image.NewGray(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: s, Y: s},
	})

	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			imOut.Set(x, y, im.GrayAt(w/2-s/2+x, h/2-s/2+y))
		}
	}
	return imOut
}

func Erode(im *image.Gray, s uint8) *image.Gray {
	if s > 7 {
		panic("failed to satisfy s ∈ {0, ..., 7}")
	}
	imOut := image.NewGray(im.Rect)

	for y := im.Rect.Min.Y + 1; y < im.Rect.Max.Y-1; y++ {
		for x := im.Rect.Min.X + 1; x < im.Rect.Max.X-1; x++ {
			v := im.GrayAt(x, y).Y
			neighborCount := -(im.GrayAt(x-1, y-1).Y + im.GrayAt(x, y-1).Y + im.GrayAt(x+1, y-1).Y + im.GrayAt(x+1, y).Y + im.GrayAt(x+1, y+1).Y + im.GrayAt(x, y+1).Y + im.GrayAt(x-1, y+1).Y + im.GrayAt(x-1, y).Y)
			if neighborCount <= s {
				v = B
			}
			imOut.Set(x, y, color.Gray{Y: v})
		}
	}
	return imOut
}

func ErodeN(im *image.Gray, strength uint8, n int) *image.Gray {
	for i := 0; i < n; i++ {
		im = Erode(im, strength)
	}
	return im
}

func Dilate(im *image.Gray, strength uint8) *image.Gray {
	if strength > 7 {
		panic("failed to satisfy strength ∈ {0, ..., 7}")
	}
	imOut := image.NewGray(im.Rect)

	for y := im.Rect.Min.Y + 1; y < im.Rect.Max.Y-1; y++ {
		for x := im.Rect.Min.X + 1; x < im.Rect.Max.X-1; x++ {
			v := im.GrayAt(x, y).Y
			neighborCount := -(im.GrayAt(x-1, y-1).Y + im.GrayAt(x, y-1).Y + im.GrayAt(x+1, y-1).Y + im.GrayAt(x+1, y).Y + im.GrayAt(x+1, y+1).Y + im.GrayAt(x, y+1).Y + im.GrayAt(x-1, y+1).Y + im.GrayAt(x-1, y).Y)
			if neighborCount >= N-strength {
				v = W
			}
			imOut.Set(x, y, color.Gray{Y: v})
		}
	}
	return imOut
}

func DilateN(im *image.Gray, strength uint8, n int) *image.Gray {
	for i := 0; i < n; i++ {
		im = Dilate(im, strength)
	}
	return im
}

func FindEdge(im *image.Gray) *image.Gray {
	imOut := image.NewGray(im.Rect)

	for y := im.Rect.Min.Y + 1; y < im.Rect.Max.Y-1; y++ {
		for x := im.Rect.Min.X + 1; x < im.Rect.Max.X-1; x++ {
			neighborCount := -(im.GrayAt(x-1, y-1).Y + im.GrayAt(x, y-1).Y + im.GrayAt(x+1, y-1).Y + im.GrayAt(x+1, y).Y + im.GrayAt(x+1, y+1).Y + im.GrayAt(x, y+1).Y + im.GrayAt(x-1, y+1).Y + im.GrayAt(x-1, y).Y)
			if 5 <= neighborCount && neighborCount <= 7 {
				imOut.Set(x, y, color.Gray{Y: W})
			} else {
				imOut.Set(x, y, color.Gray{Y: B})
			}

		}
	}
	return imOut
}
