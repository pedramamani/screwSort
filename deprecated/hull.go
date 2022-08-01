package deprecated

import (
	"image"
	"math"
)

func MinimumWidth(ps []*Point) float64 {
	var d, dp float64
	var i int
	var p0, p1 *Point
	n := len(ps)
	j := 2
	dm := math.Inf(1)

	for i = 0; i < n; i++ {
		p0, p1 = ps[i], ps[(i+1)%n]
		dp = 0
		for {
			d = SegmentPQ(p0, p1).PerpDistanceTo(ps[j%n])
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

func SuperOutline(im *image.Gray, vm float64) (ps, qs []*Point, ss []*Line) {
	for y := 0; y < im.Rect.Dy()-1; y++ {
		for x := 0; x < im.Rect.Dx()-1; x++ {
			vmp, vpp, vmm, vpm := float64(im.GrayAt(x, y).Y), float64(im.GrayAt(x+1, y).Y), float64(im.GrayAt(x, y+1).Y), float64(im.GrayAt(x+1, y+1).Y)

			if minFloat64s(vmp, vpp, vmm, vpm) <= vm && vm <= maxFloat64s(vmp, vpp, vmm, vpm) {
				gx, gy := (vpp-vmm+vpm-vmp)/2, (vpp-vmm-vpm+vmp)/2
				dv := vm - (vmp+vpp+vmm+vpm)/4
				s := dv / (math.Pow(gx, 2) + math.Pow(gy, 2))
				q := PointXY(float64(x+1), float64(im.Rect.Dy()-y-1))
				ps = append(ps, PointXY(float64(x+1)+s*gx, float64(im.Rect.Dy()-y-1)+s*gy))
				qs = append(qs, q)
				if len(ss) < 2 {
					ss = append(ss, LineABC(gx, gy, -dv).Translate(q))
				}
			}
		}
	}
	return ps, qs, ss
}
