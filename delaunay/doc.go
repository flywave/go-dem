// Package delaunay computes 2D Delaunay triangulations using exact arithmetic.
//
// The core algorithm is Bowyer-Watson incremental insertion with exact
// rational-arithmetic predicates (Orient2D, InCircle), guaranteeing
// correct results for any input including degenerate/cocircular configurations.
//
// Usage:
//
//	tris, neighbors, err := delaunay.Triangulate(x, y)
//	// tris[i] = [3]int{a, b, c} — anticlockwise vertex indices
//	// neighbors[i][j] — adjacent triangle across edge j→(j+1)%3, or -1 on hull
package delaunay
