// Package delaunay computes 2D Delaunay triangulations.
//
// The core algorithm is Bowyer-Watson incremental insertion with a
// randomized insertion order and walking point location. The geometric
// predicates (Orient2D, InCircle) use a floating-point filter with
// conservative error bounds and fall back to exact rational arithmetic
// (math/big) whenever the filter is inconclusive, so collinear and
// cocircular degeneracies are resolved exactly.
//
// Insertion starts from a finite super-triangle scaled from the input
// bounding box, not from a point at infinity. If points on the convex
// hull are nearly collinear — their deviation from the shared line
// below roughly 1e-6 of the point-set span — a circumcircle through
// them can extend beyond the super-triangle, and triangles adjacent to
// the convex hull may then violate the empty-circumcircle property.
// Inputs without such near-collinear hull edges are unaffected.
//
// Usage:
//
//	tris, neighbors, err := delaunay.Triangulate(x, y)
//	// tris[i] = [3]int{a, b, c} — anticlockwise vertex indices
//	// neighbors[i][j] — adjacent triangle across edge j→(j+1)%3, or -1 on hull
package delaunay
