package main

import (
	"screwSort/vision"
)

func main() {
	fn := "assets/data/single.png"
	im := vision.ToGray(vision.OpenPng(fn))
	im = vision.InverseThreshold(im, 150)
	out := vision.ToRgba(im)
	vision.SavePng(out, "assets/temp/out.png")
}
