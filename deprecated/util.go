package deprecated

import (
	"image"
	"image/color"
	"math"
)

func widthHeight(im *image.RGBA) (float64, float64) {
	return float64(im.Rect.Max.X - im.Rect.Min.X), float64(im.Rect.Max.Y - im.Rect.Min.Y)
}

func magnitudes(mat [][]complex128) [][]float64 {
	width := len(mat[0])
	height := len(mat)

	matOut := make([][]float64, height)
	for y := range mat {
		matOut[y] = make([]float64, width)
		for x, v := range mat[y] {
			matOut[y][x] = math.Hypot(real(v), imag(v))
		}
	}
	return matOut
}

func Matrix(im image.Gray) [][]float64 {
	width := im.Rect.Max.X - im.Rect.Min.X
	height := im.Rect.Max.Y - im.Rect.Min.Y
	if width != height {
		panic("failed to satisfy im.width == im.height")
	}

	mat := make([][]float64, width)
	for y := range mat {
		mat[y] = make([]float64, width)
		for x := range mat[y] {
			v := im.GrayAt(x, y).Y
			if v == B {
				mat[y][x] = 0
			} else {
				mat[y][x] = 1
			}
		}
	}
	return mat
}

func UnMatrix(mat [][]float64) *image.Gray {
	width := len(mat[0])
	height := len(mat)
	if width != height {
		panic("failed to satisfy len(mat) == len(mat[0])")
	}

	imOut := image.NewGray(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: width, Y: width},
	})

	for y := range mat {
		for x := range mat[y] {
			imOut.Set(x, y, color.Gray{Y: uint8(mat[y][x])})
		}
	}
	return imOut
}
