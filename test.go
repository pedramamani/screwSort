package main

import (
	"gometry/vision"
)

func main() {
	im := vision.ToGray(vision.Open("assets/leds.png"))
	imTest := vision.UpscaleRgba(vision.ToRgba(im), 8)

	ps, qs, ss := vision.SuperOutline(im, 168)
	for _, s := range ss {
		s.Scale(8).Draw(imTest, vision.Red)
	}
	for _, q := range qs {
		q.Scale(8).Draw(imTest, vision.Green)
	}
	for _, p := range ps {
		p.Scale(8).Draw(imTest, vision.Blue)
	}
	vision.Save(imTest, "assets/imTest.png")
}
