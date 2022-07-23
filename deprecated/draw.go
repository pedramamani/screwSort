package deprecated

import (
	"image"
	"image/color"
	"math"
)

func DrawPolygon(ps []image.Point, im *image.RGBA, c color.RGBA) {
	n := len(ps)
	if n < 3 {
		panic("failed to satisfy len(ps) >= 3")
	}
	for i := 0; i < n-1; i++ {
		DrawLine(ps[i], ps[i+1], im, c)
	}
	DrawLine(ps[n-1], ps[0], im, c)
}

func DrawPoints(ps []image.Point, im *image.RGBA, c color.RGBA) {
	for _, p := range ps {
		DrawPoint(p, im, c)
	}
}

func DrawLine(a image.Point, b image.Point, im *image.RGBA, c color.RGBA) {
	xi, yi := float64(a.X), float64(a.Y)
	xf, yf := float64(b.X), float64(b.Y)
	dx, dy := xf-xi, yf-yi
	m := dy / dx

	if math.Abs(m) < 1 {
		for x := xi; x*dx < xf*dx; x += math.Copysign(1, dx) {
			im.Set(int(math.Round(x)), int(math.Round(yi+(x-xi)*m)), c)
		}
	} else {
		for y := yi; y*dy < yf*dy; y += math.Copysign(1, dy) {
			im.Set(int(math.Round(xi+(y-yi)/m)), int(math.Round(y)), c)
		}
	}
}

func DrawPoint(a image.Point, im *image.RGBA, c color.RGBA) {
	im.Set(a.X, a.Y, c)
	im.Set(a.X-1, a.Y, c)
	im.Set(a.X, a.Y-1, c)
	im.Set(a.X+1, a.Y, c)
	im.Set(a.X, a.Y+1, c)
}
