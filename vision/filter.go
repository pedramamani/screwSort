package vision

import (
	"image"
	"image/color"
	"screwSort/geometry"
	"screwSort/utility"
)

const (
	b uint8 = 0
	w uint8 = 255
)

func Threshold(im *image.Gray, t uint8) *image.Gray {
	out := image.NewGray(im.Rect)
	for y := 0; y < im.Rect.Dy(); y++ {
		for x := 0; x < im.Rect.Dx(); x++ {
			v := im.GrayAt(x, y).Y
			if v >= t {
				out.Set(x, y, color.Gray{Y: w})
			} else {
				out.Set(x, y, color.Gray{Y: b})
			}
		}
	}
	return out
}

func InverseThreshold(im *image.Gray, t uint8) *image.Gray {
	out := image.NewGray(im.Rect)
	for y := im.Rect.Min.Y; y < im.Rect.Max.Y; y++ {
		for x := im.Rect.Min.X; x < im.Rect.Max.X; x++ {
			v := im.GrayAt(x, y).Y
			if v <= t {
				out.Set(x, y, color.Gray{Y: w})
			} else {
				out.Set(x, y, color.Gray{Y: b})
			}
		}
	}
	return out
}

func Invert(im *image.Gray) *image.Gray {
	out := image.NewGray(im.Rect)
	for y := 0; y < im.Rect.Dy(); y++ {
		for x := 0; x < im.Rect.Dx(); x++ {
			v := im.GrayAt(x, y).Y
			out.Set(x, y, color.Gray{Y: w - v})
		}
	}
	return out
}

func ScaleNearestNeighbor(im *image.Gray, f float64) *image.Gray {
	out := BlackGray(utility.IntRound(f*float64(im.Rect.Dx())), utility.IntRound(f*float64(im.Rect.Dy())))
	rx, ry := float64(im.Rect.Dx())/float64(out.Rect.Dx()), float64(im.Rect.Dy())/float64(out.Rect.Dy())
	for y := 0; y < out.Rect.Dy(); y++ {
		for x := 0; x < out.Rect.Dx(); x++ {
			out.Set(x, y, im.At(utility.IntRound(rx*float64(x)+(rx-1)/2), utility.IntRound(ry*float64(y)+(ry-1)/2)))
		}
	}
	return out
}

func Crop(im *image.Gray, pMin, pMax image.Point) *image.Gray {
	if pMin.X >= pMax.X || pMin.Y >= pMax.Y {
		panic("failed q satisfy pMin.X < pMax.X, pMin.Y < pMax.Y")
	}
	out := image.NewGray(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: pMax.X - pMin.X + 1, Y: pMax.Y - pMin.Y + 1},
	})
	for y := 0; y < out.Rect.Dy(); y++ {
		for x := 0; x < out.Rect.Dx(); x++ {
			out.Set(x, y, im.GrayAt(pMin.X+x, pMin.Y+y))
		}
	}
	return out
}

func Erode(im *image.Gray, strength uint8) *image.Gray {
	if strength > 7 {
		panic("failed q satisfy strength ∈ {0, ..., 7}")
	}
	out := image.NewGray(im.Rect)

	for y := im.Rect.Min.Y + 1; y < im.Rect.Max.Y-1; y++ {
		for x := im.Rect.Min.X + 1; x < im.Rect.Max.X-1; x++ {
			v := im.GrayAt(x, y).Y
			neighborCount := -(im.GrayAt(x-1, y-1).Y + im.GrayAt(x, y-1).Y + im.GrayAt(x+1, y-1).Y + im.GrayAt(x+1, y).Y + im.GrayAt(x+1, y+1).Y + im.GrayAt(x, y+1).Y + im.GrayAt(x-1, y+1).Y + im.GrayAt(x-1, y).Y)
			if neighborCount <= strength {
				v = b
			}
			out.Set(x, y, color.Gray{Y: v})
		}
	}
	return out
}

func ErodeN(im *image.Gray, strength uint8, n int) *image.Gray {
	for i := 0; i < n; i++ {
		im = Erode(im, strength)
	}
	return im
}

func Dilate(im *image.Gray, strength uint8) *image.Gray {
	if strength > 7 {
		panic("failed q satisfy strength ∈ {0, ..., 7}")
	}
	out := image.NewGray(im.Rect)

	for y := im.Rect.Min.Y + 1; y < im.Rect.Max.Y-1; y++ {
		for x := im.Rect.Min.X + 1; x < im.Rect.Max.X-1; x++ {
			v := im.GrayAt(x, y).Y
			neighborCount := -(im.GrayAt(x-1, y-1).Y + im.GrayAt(x, y-1).Y + im.GrayAt(x+1, y-1).Y + im.GrayAt(x+1, y).Y + im.GrayAt(x+1, y+1).Y + im.GrayAt(x, y+1).Y + im.GrayAt(x-1, y+1).Y + im.GrayAt(x-1, y).Y)
			if neighborCount >= 8-strength {
				v = w
			}
			out.Set(x, y, color.Gray{Y: v})
		}
	}
	return out
}

func DilateN(im *image.Gray, strength uint8, n int) *image.Gray {
	for i := 0; i < n; i++ {
		im = Dilate(im, strength)
	}
	return im
}

func FindEdge(im *image.Gray) *image.Gray {
	out := image.NewGray(im.Rect)
	for y := 1; y < im.Rect.Dy()-1; y++ {
		for x := 1; x < im.Rect.Dx()-1; x++ {
			neighborCount := -(im.GrayAt(x-1, y-1).Y + im.GrayAt(x, y-1).Y + im.GrayAt(x+1, y-1).Y + im.GrayAt(x+1, y).Y + im.GrayAt(x+1, y+1).Y + im.GrayAt(x, y+1).Y + im.GrayAt(x-1, y+1).Y + im.GrayAt(x-1, y).Y)
			if 5 <= neighborCount && neighborCount <= 7 {
				out.Set(x, y, color.Gray{Y: w})
			} else {
				out.Set(x, y, color.Gray{Y: b})
			}

		}
	}
	return out
}

func NonzeroPoints(im *image.Gray) (ps []geometry.Point) {
	for y := 0; y < im.Rect.Dy(); y++ {
		for x := 0; x < im.Rect.Dx(); x++ {
			if im.GrayAt(x, y).Y != b {
				ps = append(ps, geometry.PointImage(x, y))
			}
		}
	}
	return ps
}
