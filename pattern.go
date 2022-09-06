package main

import (
	"fmt"
	"image"
	"image/color"
)

func main() {
	W := color.Gray{Y: 255}
	B := color.Gray{}
	im := image.NewGray(image.Rectangle{Max: image.Point{X: 2, Y: 2}})
	im.Set(0, 0, W) // TL
	im.Set(0, 1, B) // BL
	im.Set(1, 0, B) // TR
	im.Set(1, 1, B) // BR
	fmt.Println(im.Pix)
}
