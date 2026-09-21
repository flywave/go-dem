package delaunay

import (
	"math"
	"math/big"
)

const (
	epsilon          = 1.1102230246251565e-16
	orientErrbound   = (3.0 + 16.0*epsilon) * epsilon
	incircleErrbound = (10.0 + 96.0*epsilon) * epsilon
)

func bigRat(v float64) *big.Rat {
	r := new(big.Rat)
	r.SetFloat64(v)
	return r
}

func orientFilter(ax, ay, bx, by, cx, cy float64) (int, bool) {
	detleft := (ax - cx) * (by - cy)
	detright := (ay - cy) * (bx - cx)
	det := detleft - detright
	errbound := orientErrbound * (math.Abs(detleft) + math.Abs(detright))
	if det > errbound {
		return 1, true
	}
	if -det > errbound {
		return -1, true
	}
	return 0, false
}

func incircleFilter(ax, ay, bx, by, cx, cy, dx, dy float64) (int, bool) {
	adx := ax - dx
	ady := ay - dy
	bdx := bx - dx
	bdy := by - dy
	cdx := cx - dx
	cdy := cy - dy

	bdxcdy := bdx * cdy
	cdxbdy := cdx * bdy
	cdxady := cdx * ady
	adxcdy := adx * cdy
	adxbdy := adx * bdy
	bdxady := bdx * ady
	alift := adx*adx + ady*ady
	blift := bdx*bdx + bdy*bdy
	clift := cdx*cdx + cdy*cdy

	det := alift*(bdxcdy-cdxbdy) + blift*(cdxady-adxcdy) + clift*(adxbdy-bdxady)
	permanent := (math.Abs(bdxcdy)+math.Abs(cdxbdy))*alift +
		(math.Abs(cdxady)+math.Abs(adxcdy))*blift +
		(math.Abs(adxbdy)+math.Abs(bdxady))*clift
	errbound := incircleErrbound * permanent
	if det > errbound {
		return 1, true
	}
	if -det > errbound {
		return -1, true
	}
	return 0, false
}

func orientRat(ax, ay, bx, by, cx, cy *big.Rat) int {
	bax := new(big.Rat).Sub(bx, ax)
	bay := new(big.Rat).Sub(by, ay)
	cax := new(big.Rat).Sub(cx, ax)
	cay := new(big.Rat).Sub(cy, ay)
	left := new(big.Rat).Mul(bax, cay)
	right := new(big.Rat).Mul(bay, cax)
	return left.Sub(left, right).Sign()
}

func incircleRat(ax, ay, bx, by, cx, cy, dx, dy *big.Rat) int {
	adx := new(big.Rat).Sub(ax, dx)
	ady := new(big.Rat).Sub(ay, dy)
	bdx := new(big.Rat).Sub(bx, dx)
	bdy := new(big.Rat).Sub(by, dy)
	cdx := new(big.Rat).Sub(cx, dx)
	cdy := new(big.Rat).Sub(cy, dy)

	ad2 := new(big.Rat).Add(new(big.Rat).Mul(adx, adx), new(big.Rat).Mul(ady, ady))
	bd2 := new(big.Rat).Add(new(big.Rat).Mul(bdx, bdx), new(big.Rat).Mul(bdy, bdy))
	cd2 := new(big.Rat).Add(new(big.Rat).Mul(cdx, cdx), new(big.Rat).Mul(cdy, cdy))

	t1 := new(big.Rat).Mul(adx, new(big.Rat).Sub(new(big.Rat).Mul(bdy, cd2), new(big.Rat).Mul(cdy, bd2)))
	t2 := new(big.Rat).Mul(ady, new(big.Rat).Sub(new(big.Rat).Mul(bdx, cd2), new(big.Rat).Mul(cdx, bd2)))
	t3 := new(big.Rat).Mul(ad2, new(big.Rat).Sub(new(big.Rat).Mul(bdx, cdy), new(big.Rat).Mul(cdx, bdy)))
	det := t1.Sub(t1, t2).Add(t1, t3)
	return det.Sign()
}

func Orient2D(ax, ay, bx, by, cx, cy float64) int {
	if s, ok := orientFilter(ax, ay, bx, by, cx, cy); ok {
		return s
	}
	return orientRat(bigRat(ax), bigRat(ay), bigRat(bx), bigRat(by), bigRat(cx), bigRat(cy))
}

func InCircle(ax, ay, bx, by, cx, cy, dx, dy float64) int {
	if s, ok := incircleFilter(ax, ay, bx, by, cx, cy, dx, dy); ok {
		return s
	}
	return incircleRat(bigRat(ax), bigRat(ay), bigRat(bx), bigRat(by), bigRat(cx), bigRat(cy), bigRat(dx), bigRat(dy))
}
