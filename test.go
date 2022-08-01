package main

import (
	"fmt"
	"screwSort/fit"
	"screwSort/vision"
)

const (
	THRESHOLD = 100
)

func main() {
	fn := "assets/photos-20220704/length.png"
	im := vision.ToGray(vision.OpenPng(fn))
	im = vision.InverseThreshold(im, THRESHOLD)
	im = vision.Dilate(im, 3)
	out := vision.ToRgba(im)

	for _, h := range vision.Hulls(im) {
		h.Simplify(5).Draw(out, vision.Red)
		p, a := fit.CenterPoint(h.Ps())
		p.Draw(out, vision.Green)
		fmt.Println(a)
	}

	vision.SavePng(out, "assets/out.png")
}
