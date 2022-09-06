package main

import (
	"fmt"
	"screwSort/vision"
	"strconv"
)

func main() {
	fn := "assets/nut.png"
	im := vision.ToGray(vision.OpenPng(fn))

	for i, h := range vision.SuperHulls(im, 120) {
		p, _ := h.CenterPoint()
		//h = h.RotateAbout(p, a)
		//ps, pe := h.Bounds()
		//h = h.Translate(-ps.X(), -ps.Y())
		//dx, dy := pe.Subtract(ps).ToImage()
		//out := vision.BlackRgba(dx, dy)
		out := vision.ApplyAlpha(vision.ToRgba(vision.InverseThreshold(im, 120)), 0.3)
		s := h.Simplify(0.4)
		s.Draw(out, vision.Red, vision.Green)
		//h.Convex().Simplify(2).Draw(out, vision.Blue)
		//p.Subtract(ps).Draw(out, vision.Cyan)
		p.Draw(out, vision.Cyan)
		vision.SavePng(out, "assets/out"+strconv.Itoa(i)+".png")
		fmt.Println(p)
		fmt.Println(len(h.Ps()))
		fmt.Println(len(s.Ps()))
	}
}
