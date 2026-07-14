package delaunay

import "math/big"

func bigRat(v float64) *big.Rat {
	r := new(big.Rat)
	r.SetFloat64(v)
	return r
}

func Orient2D(ax, ay, bx, by, cx, cy float64) int {
	bax := new(big.Rat).Sub(bigRat(bx), bigRat(ax))
	bay := new(big.Rat).Sub(bigRat(by), bigRat(ay))
	cax := new(big.Rat).Sub(bigRat(cx), bigRat(ax))
	cay := new(big.Rat).Sub(bigRat(cy), bigRat(ay))
	left := new(big.Rat).Mul(bax, cay)
	right := new(big.Rat).Mul(bay, cax)
	return left.Sub(left, right).Sign()
}

func InCircle(ax, ay, bx, by, cx, cy, dx, dy float64) int {
	adx := new(big.Rat).Sub(bigRat(ax), bigRat(dx))
	ady := new(big.Rat).Sub(bigRat(ay), bigRat(dy))
	bdx := new(big.Rat).Sub(bigRat(bx), bigRat(dx))
	bdy := new(big.Rat).Sub(bigRat(by), bigRat(dy))
	cdx := new(big.Rat).Sub(bigRat(cx), bigRat(dx))
	cdy := new(big.Rat).Sub(bigRat(cy), bigRat(dy))

	ad2 := new(big.Rat).Add(new(big.Rat).Mul(adx, adx), new(big.Rat).Mul(ady, ady))
	bd2 := new(big.Rat).Add(new(big.Rat).Mul(bdx, bdx), new(big.Rat).Mul(bdy, bdy))
	cd2 := new(big.Rat).Add(new(big.Rat).Mul(cdx, cdx), new(big.Rat).Mul(cdy, cdy))

	t1 := new(big.Rat).Mul(adx, new(big.Rat).Sub(new(big.Rat).Mul(bdy, cd2), new(big.Rat).Mul(cdy, bd2)))
	t2 := new(big.Rat).Mul(ady, new(big.Rat).Sub(new(big.Rat).Mul(bdx, cd2), new(big.Rat).Mul(cdx, bd2)))
	t3 := new(big.Rat).Mul(ad2, new(big.Rat).Sub(new(big.Rat).Mul(bdx, cdy), new(big.Rat).Mul(cdx, bdy)))
	det := t1.Sub(t1, t2).Add(t1, t3)
	return det.Sign()
}
