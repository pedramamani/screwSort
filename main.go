package main

import (
	"math"
	"screwSort/geometry"
	"screwSort/linalg"
	"screwSort/vision"
)

func main() {
	im := vision.ToGray(vision.Open("assets/nail.png"))
	im = vision.InverseThreshold(im, 150)
	//im = vision.Erode(im, 3) // filter out noise
	//im = vision.DilateN(im, 4, 2)
	//im = vision.ErodeN(im, 4, 2)

	ps := vision.GetNonzero(im)
	pc, theta := vision.CenterLine(ps)
	ps = vision.Outline(im)
	for _, p := range ps {
		p.Rotate(math.Pi - theta)
	}
	pbl := vision.BottomLeftLimit(ps)
	ptr := vision.TopRightLimit(ps).Subtract(pbl)
	for _, p := range ps {
		p.Subtract(pbl)
	}
	ssf := linesToSegments(regressionLines(ps, 2))
	ssc := linesToSegments(regressionLines(ps, 12))

	imCenter := vision.ToRgba(im)
	geometry.LinePTheta(pc, theta).Draw(imCenter, vision.Red)
	geometry.LinePTheta(pc, theta+math.Pi/2).Draw(imCenter, vision.Red)
	imFinal := vision.BlankRgba(int(ptr.X()), int(ptr.Y()), vision.Black)
	//for _, p := range ps {
	//	p.Draw(imFinal, vision.White)
	//}
	for _, sf := range ssf {
		sf.Draw(imFinal, vision.Red)
	}
	for _, sc := range ssc {
		sc.Draw(imFinal, vision.Green)
	}
	//fmt.Println("Angle (rad), MetricLength (px)")
	//for i := 0; i < len(ss)-1; i++ {
	//	fmt.Printf("%f, %f\n", ss[i].AngleBetween(ss[i+1]), ss[i].D().R())
	//}
	//fmt.Printf("Outline points reduced from %d to %d\n", len(ps), len(ss))
	vision.Save(imCenter, "assets/imCenter.png")
	vision.Save(imFinal, "assets/imFinal.png")
}

func regressionLines(ps []*geometry.Point, t float64) []*geometry.Line {
	var ls []*geometry.Line
	var lp *geometry.Line
	ip := 0
	for i := 2; i < len(ps); i++ {
		l, r := linalg.OrthogonalRegression(ps[ip:i])
		if r > t {
			ls = append(ls, lp)
			ip = i - 2 // todo: or could be i - 1 if we don't want to chain with previous fit
			continue
		}
		lp = l
	}
	ls = append(ls, lp)
	return ls
}

func linesToSegments(ls []*geometry.Line) []*geometry.Segment {
	var ss []*geometry.Segment
	for i := 1; i < len(ls)-1; i++ {
		ss = append(ss, geometry.SegmentPQ(
			ls[i].IntersectionWith(ls[i-1]),
			ls[i].IntersectionWith(ls[i+1]),
		))
	}
	p := ls[0].IntersectionWith(ls[len(ls)-1])
	ss = append(ss,
		geometry.SegmentPQ(p, ss[0].P()),
		geometry.SegmentPQ(ss[len(ss)-1].Q(), p),
	)
	return ss
}

func linesToPoints(ls []*geometry.Line) []*geometry.Point {
	var ps []*geometry.Point
	n := len(ls)
	for i := 0; i < n; i++ {
		ps = append(ps, ls[i].IntersectionWith(ls[(i+1)%n]))
	}
	return ps
}
