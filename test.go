package main

import (
	"screwSort/vision"
)

func main() {
	fn := "assets/photos/head.png"
	im := vision.ToGray(vision.OpenPng(fn))
	im = vision.InverseThreshold(im, 60)
	out := vision.ToRgba(im)
	vision.SavePng(out, "assets/out.png")
}
