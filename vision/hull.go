package vision

import (
	"image"
	"image/color"
	"screwSort/fit"
	"screwSort/geometry"
	"sort"
)

const (
	thresholdArea = 500.
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
		if f.MeanError() > threshold {
			ls = append(ls, lp)
			ip = i - 2
			continue
		}
		lp = f.Line()
	}
	ls = append(ls, lp)

	var ps []geometry.Point
	n := len(ls)
	for i := 0; i < n; i++ {
		ps = append(ps, ls[i].IntersectionWith(ls[(i+1)%n]))
	}
	return HullPs(ps)
}

func (h Hull) Draw(im *image.RGBA, c color.RGBA) {
	n := len(h.ps)
	for i, p := range h.ps {
		geometry.SegmentPQ(p, h.ps[(i+1)%n]).Draw(im, c)
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

	for len(links) > 0 {
		var qs []image.Point
		var q, qi image.Point
		var e bool

		maxY := 0
		for q = range links {
			if q.Y > maxY {
				maxY = q.Y
				qi = q
			}
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

/* todo: implement a way to compare two hulls
https://en.wikipedia.org/wiki/Fr%C3%A9chet_distance
https://en.wikipedia.org/wiki/Dynamic_time_warping

int DTWDistance(s: array [1..n], t: array [1..m], w: int) {
    DTW := array [0..n, 0..m]

    w := max(w, abs(n-m)) // adapt window size (*)

    for i := 0 to n
        for j:= 0 to m
            DTW[i, j] := infinity
    DTW[0, 0] := 0
    for i := 1 to n
        for j := max(1, i-w) to min(m, i+w)
            DTW[i, j] := 0

    for i := 1 to n
        for j := max(1, i-w) to min(m, i+w)
            cost := D(s[i], t[j])
            DTW[i, j] := cost + minimum(DTW[i-1, j  ],    // insertion
                                        DTW[i  , j-1],    // deletion
                                        DTW[i-1, j-1])    // match
    return DTW[n, m]
}
*/
