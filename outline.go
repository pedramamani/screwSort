package main

import (
	"fmt"
	"screwSort/vision"
	"strconv"
)

func main() {
	data := map[string][]float64{
		"single": {150, 0.38, 9.5},
		"washer": {80, 1, 8},
		"length": {120, 1, 8},
	}
	key := "length"

	fn := fmt.Sprintf("assets/data/%s.png", key)
	im := vision.ToGray(vision.OpenPng(fn))

	for i, h := range vision.SuperHulls(im, data[key][0]) {
		p, _ := h.CenterPoint()
		out := vision.ApplyAlpha(vision.ToRgba(vision.InverseThreshold(im, uint8(data[key][0]))), 0.3)
		s := h.Simplify(data[key][1], data[key][2])
		s.Draw(out, vision.Red, vision.Green)
		p.Draw(out, vision.Cyan)
		vision.SavePng(out, "assets/temp/out"+strconv.Itoa(i)+".png")
		fmt.Printf("CM: %s\nOutline points reduced from %d to %d\n\n", p, len(h.Ps()), len(s.Ps()))
	}
}
