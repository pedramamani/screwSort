package vision

import (
	"gometry/geometry"
	"image"
	"math"
	"sort"
)

func GetNonzero(im *image.Gray) (ps []*geometry.Point) {
	for y := im.Rect.Min.Y; y < im.Rect.Max.Y; y++ {
		for x := im.Rect.Min.X; x < im.Rect.Max.X; x++ {
			if im.GrayAt(x, y).Y != B {
				ps = append(ps, geometry.PointImage(x, y, im))
			}
		}
	}
	return ps
}

func Outline(im *image.Gray) (ps []*geometry.Point) {
	links := make(map[image.Point]image.Point)

	for y := im.Rect.Min.Y; y < im.Rect.Max.Y-1; y++ {
		for x := im.Rect.Min.X; x < im.Rect.Max.X-1; x++ {
			tl, tr, bl, br := im.GrayAt(x, y).Y, im.GrayAt(x+1, y).Y, im.GrayAt(x, y+1).Y, im.GrayAt(x+1, y+1).Y
			switch {
			case tl == W && tr == B && bl == W && br == W:
				links[image.Point{X: x + 1, Y: y + 1}] = image.Point{X: x, Y: y}
			case tl == B && tr == W && bl == W && br == W:
				links[image.Point{X: x + 1, Y: y}] = image.Point{X: x, Y: y + 1}
			case tl == W && tr == W && bl == B && br == W:
				links[image.Point{X: x, Y: y}] = image.Point{X: x + 1, Y: y + 1}
			case tl == W && tr == W && bl == W && br == B:
				links[image.Point{X: x, Y: y + 1}] = image.Point{X: x + 1, Y: y}
			case tl == B && tr == B && bl == W && br == W:
				links[image.Point{X: x + 1, Y: y + 1}] = image.Point{X: x, Y: y + 1}
			case tl == B && tr == W && bl == B && br == W:
				links[image.Point{X: x + 1, Y: y}] = image.Point{X: x + 1, Y: y + 1}
			case tl == W && tr == B && bl == W && br == B:
				links[image.Point{X: x, Y: y + 1}] = image.Point{X: x, Y: y}
			case tl == W && tr == W && bl == B && br == B:
				links[image.Point{X: x, Y: y}] = image.Point{X: x + 1, Y: y}
			}
		}
	}

	var qs []image.Point
	var q, qi image.Point
	maxY := 0
	for q = range links {
		if q.Y > maxY {
			maxY = q.Y
			qi = q
		}
	}
	qs = append(qs, qi)
	for { // todo: update to include all outlines in image (check wrapping number to determine inner/outer)
		q = links[qi]
		qs = append(qs, q)
		if q == qs[0] {
			break
		}
		qi = q
	}

	for _, q = range qs {
		ps = append(ps, geometry.PointImage(q.X, q.Y, im))
	}
	return ps
}

func BottomLeftLimit(ps []*geometry.Point) *geometry.Point {
	minX := math.Inf(1)
	minY := math.Inf(1)
	for _, p := range ps {
		if p.Y() < minY {
			minY = p.Y()
		}
		if p.X() < minX {
			minX = p.X()
		}
	}
	return geometry.PointXY(minX, minY)
}

func TopRightLimit(ps []*geometry.Point) *geometry.Point {
	maxX := math.Inf(-1)
	maxY := math.Inf(-1)
	for _, p := range ps {
		if p.Y() > maxY {
			maxY = p.Y()
		}
		if p.X() > maxX {
			maxX = p.X()
		}
	}
	return geometry.PointXY(maxX, maxY)
}

func ConvexHull(ps []*geometry.Point) (qs []*geometry.Point) {
	pi := ps[0]
	for _, p := range ps {
		if p.Y() < pi.Y() || p.Y() == pi.Y() && p.X() < pi.X() {
			pi = p
		}
	}

	m := make(map[float64]int, len(ps))
	for i, p := range ps {
		a := pi.AngleTo(p)
		ie, e := m[a]
		if !e || pi.DistanceTo(ps[ie]) < pi.DistanceTo(p) {
			m[a] = i
		}
	}

	mr := make(map[int]float64, len(ps))
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

	qs = append(qs, ps[is[0]], ps[is[1]])
	n := 2
	for _, i := range is[2:] {
		p := ps[i]
		for n > 1 && qs[n-2].OrientationOf(qs[n-1], p) <= 0 {
			qs = qs[:len(qs)-1]
			n--
		}
		qs = append(qs, p)
		n++
	}
	return qs
}

func MinimumWidth(ps []*geometry.Point) float64 {
	var d, dp float64
	var i int
	var p0, p1 *geometry.Point
	n := len(ps)
	j := 2
	dm := math.Inf(1)

	for i = 0; i < n; i++ {
		p0, p1 = ps[i], ps[(i+1)%n]
		dp = 0
		for {
			d = geometry.SegmentPQ(p0, p1).PerpDistanceTo(ps[j%n])
			if d < dp {
				j--
				break
			}
			dp = d
			j++
		}
		if dp < dm {
			dm = dp
		}
	}
	return dm
}

func Moments(ps []*geometry.Point) (M00, M10, M01, M11, M20, M02 float64) {
	M00 = float64(len(ps))
	for _, p := range ps {
		M10 += p.X()
		M01 += p.Y()
		M11 += p.X() * p.Y()
		M20 += math.Pow(p.X(), 2)
		M02 += math.Pow(p.Y(), 2)
	}
	return
}

func CenterLine(ps []*geometry.Point) (*geometry.Point, float64) {
	M00, M10, M01, M11, M20, M02 := Moments(ps)
	xBar := M10 / M00
	yBar := M01 / M00
	m20 := M20/M00 - math.Pow(xBar, 2)
	m02 := M02/M00 - math.Pow(yBar, 2)
	m11 := M11/M00 - xBar*yBar
	theta := math.Atan(2*m11/(m20-m02))/2 + math.Pi/2
	pc := geometry.PointXY(M10/M00, M01/M00)
	return pc, theta
}

func SuperOutline(im *image.Gray, vm float64) (ps, qs []*geometry.Point, ss []*geometry.Line) {
	for y := 0; y < im.Rect.Dy()-1; y++ {
		for x := 0; x < im.Rect.Dx()-1; x++ {
			vmp, vpp, vmm, vpm := float64(im.GrayAt(x, y).Y), float64(im.GrayAt(x+1, y).Y), float64(im.GrayAt(x, y+1).Y), float64(im.GrayAt(x+1, y+1).Y)

			if MinFloat64s(vmp, vpp, vmm, vpm) <= vm && vm <= MaxFloat64s(vmp, vpp, vmm, vpm) {
				gx, gy := (vpp-vmm+vpm-vmp)/2, (vpp-vmm-vpm+vmp)/2
				dv := vm - (vmp+vpp+vmm+vpm)/4
				s := dv / (math.Pow(gx, 2) + math.Pow(gy, 2))
				q := geometry.PointXY(float64(x+1), float64(im.Rect.Dy()-y-1))
				ps = append(ps, geometry.PointXY(float64(x+1)+s*gx, float64(im.Rect.Dy()-y-1)+s*gy))
				qs = append(qs, q)
				if len(ss) < 2 {
					ss = append(ss, geometry.LineABC(gx, gy, -dv).Translate(q))
				}
			}
		}
	}
	return ps, qs, ss
}

func MinFloat64s(vs ...float64) float64 {
	vm := vs[0]
	for _, v := range vs {
		if v < vm {
			vm = v
		}
	}
	return vm
}

func MaxFloat64s(vs ...float64) float64 {
	vm := vs[0]
	for _, v := range vs {
		if v > vm {
			vm = v
		}
	}
	return vm
}

/* todo: implement a way to compare two outlines (used to match a blob to respective object)
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
            cost := d(s[i], t[j])
            DTW[i, j] := cost + minimum(DTW[i-1, j  ],    // insertion
                                        DTW[i  , j-1],    // deletion
                                        DTW[i-1, j-1])    // match
    return DTW[n, m]
}
*/
