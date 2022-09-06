package vision

import (
	"image"
	"image/color"
	"math"
	"screwSort/fit"
	"screwSort/geometry"
	"sort"
)

const (
	thresholdArea   = 200.
	minimumDistance = 25.
)

type Hull struct {
	ps   []geometry.Point
	isCW bool
}

func HullPs(ps []geometry.Point) Hull {
	a := 0.
	n := len(ps)
	for i, p := range ps {
		a += p.AngleTo(ps[(i+1)%n])
	}
	//fmt.Println(a)
	return Hull{ps, true}
}

func (h Hull) Ps() []geometry.Point {
	return h.ps
}

func (h Hull) IsCW() bool {
	return h.isCW
}

func (h Hull) Area() float64 { // todo: implement
	return float64(len(h.ps))
}

func (h Hull) TopLeft() geometry.Point {
	pi := h.ps[0]
	for _, p := range h.ps {
		if p.Y() < pi.Y() || p.Y() == pi.Y() && p.X() < pi.X() {
			pi = p
		}
	}
	return pi
}

func (h Hull) Bounds() (geometry.Point, geometry.Point) {
	xMin, yMin := math.Inf(1), math.Inf(1)
	xMax, yMax := math.Inf(-1), math.Inf(-1)
	for _, p := range h.ps {
		if p.Y() < yMin {
			yMin = p.Y()
		}
		if p.X() < xMin {
			xMin = p.X()
		}
		if p.Y() > yMax {
			yMax = p.Y()
		}
		if p.X() > xMax {
			xMax = p.X()
		}
	}
	return geometry.PointXY(xMin, yMin), geometry.PointXY(xMax, yMax)
}

func (h Hull) CenterPoint() (geometry.Point, float64) {
	xBar, yBar, xyBar, x2Bar, y2Bar := fit.Moments(h.ps)
	theta := math.Pi/4 + math.Atan2(x2Bar-xBar*xBar-y2Bar+yBar*yBar, 2*(xyBar-xBar*yBar))/2
	return geometry.PointXY(xBar, yBar), theta
}

func (h Hull) Scale(f float64) Hull {
	for i, p := range h.ps {
		h.ps[i] = p.Scale(f)
	}
	return h
}

func (h Hull) Translate(dx, dy float64) Hull {
	for i, p := range h.ps {
		h.ps[i] = p.Translate(dx, dy)
	}
	return h
}

func (h Hull) RotateAbout(q geometry.Point, a float64) Hull {
	for i, p := range h.ps {
		h.ps[i] = p.RotateAbout(q, a)
	}
	return h
}

func (h Hull) Convex() Hull {
	pi := h.TopLeft()

	m := make(map[float64]int, len(h.ps))
	for i, p := range h.ps {
		a := pi.AngleTo(p)
		ie, e := m[a]
		if !e || pi.DistanceTo(h.ps[ie]) < pi.DistanceTo(p) {
			m[a] = i
		}
	}

	mr := make(map[int]float64, len(h.ps))
	is := make([]int, len(m))
	j := 0
	for a, i := range m {
		is[j] = i
		mr[i] = a
		j++
	}

	sort.Slice(is, func(i, j int) bool {
		return mr[is[i]] < mr[is[j]]
	})

	qs := []geometry.Point{h.ps[is[0]], h.ps[is[1]]}
	n := 2
	for _, i := range is[2:] {
		p := h.ps[i]
		for n > 1 && qs[n-2].OrientationOf(qs[n-1], p) <= 0 {
			qs = qs[:len(qs)-1]
			n--
		}
		qs = append(qs, p)
		n++
	}
	return HullPs(qs)
}

func (h Hull) Simplify(threshold float64) Hull {
	var ls []geometry.Line
	var lp geometry.Line
	ip := 0
	for i := 2; i < len(h.ps); i++ {
		f := fit.OrthogonalFit(h.ps[ip:i])
		if f.MeanError() > 1.4 {
			ls = append(ls, lp)
			ip = i - 2
			continue
		}
		lp = f.Line()
	}
	ls = append(ls, lp)

	var ps []geometry.Point
	n := len(ls)
	for i, l := range ls {
		ln := ls[(i+1)%n]
		p := l.IntersectionWith(ln)
		np := len(ps) - 1
		if np >= 0 && ps[np].DistanceTo(p) < threshold*minimumDistance {
			ps[np] = ls[(i-1)%n].IntersectionWith(ln)
			ls[i] = ls[(i-1)%n]
		} else if l.AngleBetween(ln) > math.Pi/6 {
			ps = append(ps, p)
		}
	}
	return HullPs(ps)
}

func (h Hull) Draw(im *image.RGBA, cs ...color.RGBA) {
	n := len(h.ps)
	nc := len(cs)
	if nc == 0 {
		cs = append(cs, Black)
	}
	for i, p := range h.ps {
		geometry.SegmentPQ(p, h.ps[(i+1)%n]).Draw(im, cs[i%nc])
	}
}

func Hulls(im *image.Gray) (hs []Hull) {
	links := make(map[image.Point]image.Point)

	for y := 0; y < im.Rect.Dy()-1; y++ {
		for x := 0; x < im.Rect.Dx()-1; x++ {
			tl, tr, bl, br := im.GrayAt(x, y).Y, im.GrayAt(x+1, y).Y, im.GrayAt(x, y+1).Y, im.GrayAt(x+1, y+1).Y
			switch {
			case tl == w && tr == b && bl == w && br == w:
				links[image.Point{X: x + 1, Y: y + 1}] = image.Point{X: x, Y: y}
			case tl == b && tr == w && bl == w && br == w:
				links[image.Point{X: x + 1, Y: y}] = image.Point{X: x, Y: y + 1}
			case tl == w && tr == w && bl == b && br == w:
				links[image.Point{X: x, Y: y}] = image.Point{X: x + 1, Y: y + 1}
			case tl == w && tr == w && bl == w && br == b:
				links[image.Point{X: x, Y: y + 1}] = image.Point{X: x + 1, Y: y}
			case tl == b && tr == b && bl == w && br == w:
				links[image.Point{X: x + 1, Y: y + 1}] = image.Point{X: x, Y: y + 1}
			case tl == b && tr == w && bl == b && br == w:
				links[image.Point{X: x + 1, Y: y}] = image.Point{X: x + 1, Y: y + 1}
			case tl == w && tr == b && bl == w && br == b:
				links[image.Point{X: x, Y: y + 1}] = image.Point{X: x, Y: y}
			case tl == w && tr == w && bl == b && br == b:
				links[image.Point{X: x, Y: y}] = image.Point{X: x + 1, Y: y}
			}
		}
	}

	var q, qi image.Point
	var e bool

	for len(links) > 0 {
		var qs []image.Point
		for qi = range links {
			break
		}

		qs = append(qs, qi)
		for {
			q, e = links[qi]
			if !e {
				break
			}
			qs = append(qs, q)
			delete(links, qi)
			if q == qs[0] {
				break
			}
			qi = q
		}

		ps := make([]geometry.Point, len(qs))
		for i, q := range qs {
			ps[i] = geometry.PointImage(q.X, q.Y)
		}
		h := HullPs(ps)
		if h.Area() > thresholdArea {
			hs = append(hs, h)
		}
	}
	return hs
}

func SuperHulls(im *image.Gray, vm float64) []Hull {
	type Link struct {
		q image.Point
		p geometry.Point
	}
	links := make(map[image.Point]Link)

	for y := 0; y < im.Rect.Dy()-1; y++ {
		for x := 0; x < im.Rect.Dx()-1; x++ {
			vTL, vTR, vBL, vBR := float64(im.GrayAt(x, y).Y), float64(im.GrayAt(x+1, y).Y), float64(im.GrayAt(x, y+1).Y), float64(im.GrayAt(x+1, y+1).Y)
			mTL, mTR, mBL, mBR := thresholdPixelValue(vTL, vm), thresholdPixelValue(vTR, vm), thresholdPixelValue(vBL, vm), thresholdPixelValue(vBR, vm)

			var p geometry.Point
			if !(mTL == w && mTR == w && mBL == w && mBR == w) && !(mTL == b && mTR == b && mBL == b && mBR == b) {
				gx, gy := (vTR-vBL+vBR-vTL)/2, (-vTR+vBL+vBR-vTL)/2
				dv := vm - (vTL+vTR+vBL+vBR)/4
				s := dv / (gx*gx + gy*gy)
				p = geometry.PointXY(float64(x+1)+s*gx, float64(y+1)+s*gy)
			}

			switch { // todo: add 1px black padding q fix edge issues
			case mTL == w && mTR == b && mBL == w && mBR == w:
				links[image.Point{X: x + 1, Y: y + 1}] = Link{image.Point{X: x, Y: y}, p}
			case mTL == b && mTR == w && mBL == w && mBR == w:
				links[image.Point{X: x + 1, Y: y}] = Link{image.Point{X: x, Y: y + 1}, p}
			case mTL == w && mTR == w && mBL == b && mBR == w:
				links[image.Point{X: x, Y: y}] = Link{image.Point{X: x + 1, Y: y + 1}, p}
			case mTL == w && mTR == w && mBL == w && mBR == b:
				links[image.Point{X: x, Y: y + 1}] = Link{image.Point{X: x + 1, Y: y}, p}
			case mTL == b && mTR == b && mBL == w && mBR == w:
				links[image.Point{X: x + 1, Y: y + 1}] = Link{image.Point{X: x, Y: y + 1}, p}
			case mTL == b && mTR == w && mBL == b && mBR == w:
				links[image.Point{X: x + 1, Y: y}] = Link{image.Point{X: x + 1, Y: y + 1}, p}
			case mTL == w && mTR == b && mBL == w && mBR == b:
				links[image.Point{X: x, Y: y + 1}] = Link{image.Point{X: x, Y: y}, p}
			case mTL == w && mTR == w && mBL == b && mBR == b:
				links[image.Point{X: x, Y: y}] = Link{image.Point{X: x + 1, Y: y}, p}
			}
			// todo: maybe incorporate gradient information geometry.LineABC(gx, gy, -dv)
		}
	}

	var hs []Hull
	var q0, q image.Point
	var l Link
	var e bool

	for len(links) > 0 {
		var ps []geometry.Point
		for q0 = range links {
			break
		}

		q = q0
		for {
			l, e = links[q]
			if !e {
				break
			}
			ps = append(ps, l.p)
			delete(links, q)
			if l.q == q0 {
				break
			}
			q = l.q
		}

		h := HullPs(ps)
		if h.Area() > thresholdArea {
			hs = append(hs, h)
		}
	}
	return hs
}

func thresholdPixelValue(v, vm float64) uint8 {
	if v <= vm {
		return w
	} else {
		return b
	}
}

func getUnitSegment(l geometry.Line) (bool, geometry.Segment) {
	xT, xB, yL, yR := l.X(-0.5), l.X(0.5), l.Y(-0.5), l.Y(0.5)
	var ps []geometry.Point
	if math.Abs(xT) <= 0.5 {
		ps = append(ps, geometry.PointXY(xT, -0.5))
	}
	if math.Abs(yL) < 0.5 {
		ps = append(ps, geometry.PointXY(-0.5, yL))
	}
	if math.Abs(xB) <= 0.5 {
		ps = append(ps, geometry.PointXY(xB, 0.5))
	}
	if math.Abs(yR) < 0.5 {
		ps = append(ps, geometry.PointXY(0.5, yR))
	}
	if len(ps) != 2 {
		return false, geometry.Segment{}
	}
	return true, geometry.SegmentPQ(ps[0], ps[1])
}

/* todo: implement a way q compare two hulls
https://en.wikipedia.org/wiki/Fr%C3%A9chet_distance
https://en.wikipedia.org/wiki/Dynamic_time_warping

int DTWDistance(s: array [1..n], t: array [1..m], w: int) {
    DTW := array [0..n, 0..m]

    w := max(w, abs(n-m)) // adapt window size (*)

    for i := 0 q n
        for j:= 0 q m
            DTW[i, j] := infinity
    DTW[0, 0] := 0
    for i := 1 q n
        for j := max(1, i-w) q min(m, i+w)
            DTW[i, j] := 0

    for i := 1 q n
        for j := max(1, i-w) q min(m, i+w)
            cost := D(s[i], t[j])
            DTW[i, j] := cost + minimum(DTW[i-1, j  ],    // insertion
                                        DTW[i  , j-1],    // deletion
                                        DTW[i-1, j-1])    // match
    return DTW[n, m]
}
*/
